package models

import (
	"errors"
	"fmt"
)

type User struct {
	Id        int    `gorm:"primary_key;column:F_user_id;type:INT(11)"`
	LoginName string `gorm:"column:F_user_name;type:VARCHAR(100)"`
	Contact   string `gorm:"column:F_contact;"`
	Role      int    `gorm:"column:F_role;type:TINYINT(4)"`
	Enable    int    `gorm:"column:F_enable;type:TINYINT(4)"`
	Password  string `gorm:"column:F_password;type:VARCHAR(255)"`
	Salt      string `gorm:"column:F_salt;type:VARCHAR(10)"`
}

func AdminGetByName(loginName string) (*User, error) {
	//查表
	var newUser User
	err := GetDb().Table("t_admins").Where("F_user_name = ?", loginName).First(&newUser).Error

	if err != nil {
		return nil, errors.New("数据库异常" + err.Error())
	}

	return &newUser, nil
}

func AdminGetById(id int) (*User, error) {
	//查表
	var newUser User
	err := GetDb().Table("t_admins").Where("F_user_id = ?", id).First(&newUser).Error

	if err != nil {
		return nil, errors.New("数据库异常" + err.Error())
	}

	return &newUser, nil
}

func HasUserName(userName string) bool {
	//查表
	var n User
	notFound := GetDb().Table("t_admins").Where(" binary F_user_name = ?", userName).Select("F_user_id").Find(&n).RecordNotFound()
	return !notFound
}

func SaveUser(id int, userName string, Contact string, passMd5 string, passSalt string) error {
	if id != 0 {
		tx := GetDb().Table("t_admins").Begin()
		updated := make(map[string]interface{})

		if len(userName) > 0 {
			updated["F_user_name"] = userName
		}

		if len(Contact) > 0 {
			updated["F_contact"] = Contact
		}
		fmt.Println("passMd5:", passMd5)
		fmt.Println("passSalt:", passSalt)
		if len(passMd5) > 0 {
			updated["F_password"] = passMd5

			if len(passSalt) > 0 {
				updated["F_salt"] = passSalt
			} else {
				delete(updated, "F_password")
			}
		}

		err := tx.Model(&User{}).Where("F_user_id = ?", id).Updates(updated).Error
		if err != nil {
			return HandleErrByTx(errors.New("更新用户信息失败:"+err.Error()), tx)
		}

		tx.Commit()
		return nil
	}
	return errors.New("ID 不能为0")
}

func GetAdminList(limit int, page int) (result []User, count int64) {
	result = make([]User, 0)

	//处理分页参数
	var offset int
	if limit > 0 && page > 0 {
		offset = (page - 1) * limit
	}

	GetDb().Table("t_admins").Count(&count)
	GetDb().Table("t_admins").Limit(limit).Offset(offset).Find(&result)

	return
}

func ChangeUserStatus(newStatus, id int) error {
	if newStatus != 0 && newStatus != 1 {
		return errors.New("只能为启用/禁用")
	}

	if id != 0 {
		var roles []int
		GetDb().Table("t_admins").Model(&User{}).Where("F_user_id = ?", id).Pluck("F_role", &roles)

		if len(roles) > 0 {
			if roles[0] == ADMIN_SUPER {
				return errors.New("不能修改Super的状态")
			}
			err := GetDb().Table("t_admins").Model(&User{}).Where("F_user_id = ?", id).
				UpdateColumn("F_enable", newStatus).Error
			return err
		}
		return errors.New("没有该用户")
	}
	return errors.New("ID不能为0")
}

func AddUser(loginName string, contact string, passMd5 string, passSalt string, role int) error {
	user := User{
		LoginName: loginName,
		Contact:   contact,
		Password:  passMd5,
		Salt:      passSalt,
		Role:      role,
		Enable:    1,
	}
	fmt.Println(user)
	err := GetDb().Table("t_admins").Create(&user).Error
	fmt.Println(err)
	if err != nil {
		return err
	}
	return nil
}
