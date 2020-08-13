package helper

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	mathrand "math/rand"
	"time"
)

//sha1
func Sha1(str string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(str)))
}

//生成Guid字串(32位) 3a0f6874b37e0f12f8e9ea113985ad89
func GetGuid() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return Md5([]byte(base64.URLEncoding.EncodeToString(b)))
}

func GetSaveToken() string {
	//token = 随机8位数 + 时间戳
	ntime := time.Now().Unix()
	rnd := mathrand.New(mathrand.NewSource(ntime))
	vcode := fmt.Sprintf("%08v", rnd.Int31n(10000000))

	return vcode + Int64ToString(ntime)
}
