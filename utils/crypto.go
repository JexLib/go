package utils

import (
	"crypto/md5"
	"encoding/hex"
)

//md5字符串计算
func Md5String(str string) string {
	md5h := md5.New()

	md5h.Write([]byte(str))
	cipherStr := md5h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
