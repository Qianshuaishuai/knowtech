package helper

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

//类型转化 string  to int
func StrToInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

//类型转化 string  to uint64
func StrToUint64(str string) uint64 {
	i, _ := strconv.ParseUint(str, 0, 64)
	return i
}

//类型转化 string  to float64
func StrToFloat64(str string) float64 {
	f, _ := strconv.ParseFloat(str, 64)
	return f
}

//类型转化 int to string
func IntToString(i int) string {
	return fmt.Sprintf("%d", i)
}

//类型转化 int64 to string
func Int64ToString(i int64) string {
	return fmt.Sprintf("%d", i)
}

//类型转化 uint64 to string
func Uint64ToString(i uint64) string {
	return fmt.Sprintf("%d", i)
}

//类型转化 uint32 to string
func Uint32ToString(i uint32) string {
	return fmt.Sprintf("%d", i)
}

//类型转换inerface to string
func InterfaceToString(data interface{}) string {
	return fmt.Sprintf("%s", data)
}

//合并int数组
func JoinInt(list []int, flag string) string {
	result := ""
	if len(list) > 0 {
		for _, v := range list {
			result += IntToString(v) + flag
		}
		result = strings.Trim(result, flag)
	}
	return result
}

func JoinInt64(list []int64, flag string) string {
	result := ""
	if len(list) > 0 {
		for _, v := range list {
			result += Int64ToString(v) + flag
		}
		result = strings.Trim(result, flag)
	}
	return result
}

//获取interface的类型
func GetInterfaceType(i interface{}) string {
	typeObj := reflect.TypeOf(i)
	return typeObj.Kind().String()
}

//深度复制一个对象
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

//去重一个int数组
func RmDuplicateInt(list *[]int) []int {
	var x []int = []int{}
	for _, i := range *list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if i == v {
					break
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}
	return x
}

//uint数组中是否有元素elem
func UIntContainer(values []uint, elem uint) bool {
	found := false
	for _, searchValue := range values {
		if elem == searchValue {
			found = true
			break
		}
	}
	if !found {
		return false
	}
	return true
}

//将insertion插入到slice中，index指定插入位置
func Insert(slice []byte, insertion byte, index int) []byte {
	rear := make([]byte, len(slice)+1)

	if len(slice) == 1 && index == 0 {
		rear[0] = insertion
		rear[1] = slice[0]
		return rear
	}

	for i := range slice[:index] {
		rear[i] = slice[i]
	}

	rear[index] = insertion

	for i := index; i < len(slice); i++ {
		rear[i+1] = slice[i]
	}
	return rear
}
