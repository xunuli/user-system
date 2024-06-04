package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// 比对，限制一些字段的范围，例如：性别只有男和女
func Contains(source []string, tg string) bool {
	for _, s := range source {
		if s == tg {
			return true
		}
	}
	return false
}

// 加密算法
func Md5string(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	str := hex.EncodeToString(h.Sum(nil))
	return str
}

// 返回经过M5加密后的字符串session
func GenerateSession(userName string) string {
	return Md5string(fmt.Sprintf("%s:%s", userName, "session"))
}
