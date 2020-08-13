package helper

import (
	"encoding/json"
	"strings"
	"unicode/utf8"
	"crypto/md5"
	"fmt"
	"time"
	"math/rand"
)

// 字母 ABC 转换成 123
func MapABCto123(s string) string {
	l, _ := utf8.DecodeRuneInString(s)
	return IntToString(int(l) - 65)
}

func Map123toABC(i int) string {
	var ru rune
	ru = rune(65 + i)
	return "答案" + string(ru)
}

//检查一个字符串是否在字符串数组里面
func StringInArray(value string, list []string) bool {
	result := false
	for _, item := range list {
		if value == item {
			result = true
			break
		}
	}
	return result
}

//合并字符串数组
func JoinString(list []string, flag string) string {
	result := ""
	if len(list) > 0 {
		for _, v := range list {
			result += v + flag
		}
		result = strings.Trim(result, flag)
	}
	return result
}

//将[id,id,id]字符串转换成id数组
func TransformStringToInt64Arr(idsString string) ([]int64, error) {
	resourceIdList := make([]int64, 0)
	dec := json.NewDecoder(strings.NewReader(idsString))
	dec.UseNumber()
	errJ := dec.Decode(&resourceIdList)
	return resourceIdList, errJ
}

//将id数组转[id,id,id]字符串
func TransformInt64ArrToString(idsString []int64) string {
	return "[" + JoinInt64(idsString, ",") + "]"
}


func Md5(buf []byte) string {
	hash := md5.New()
	hash.Write(buf)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func Password(len int, pwdO string) (pwd string, salt string) {
	salt = GetRandomString(len)
	defaultPwd := "george518"
	if pwdO != "" {
		defaultPwd = pwdO
	}
	pwd = Md5([]byte(defaultPwd + salt))
	return pwd, salt
}

//生成随机字符串
func GetRandomString(lens int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < lens; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
