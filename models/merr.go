package models

import (
	"github.com/jinzhu/gorm"
	"github.com/HYY-yu/LogLib"
)

func HandleErrByTx(err error, tx *gorm.DB) error {
	if tx != nil {
		tx.Rollback()
	}
	//处理err
	loglib.GetLogger().LogInfo("Has A err : %s", err.Error())

	return err
}
