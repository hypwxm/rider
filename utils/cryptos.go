package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"crypto/rand"
	"encoding/base64"
)


//获取MD5值
func Md5String(s string) string {
	ms := md5.New()
	ms.Write([]byte(s))
	return hex.EncodeToString(ms.Sum(nil))
}

//创建指定长度的随机字符串
func RandString() string {
	b := make([]byte, 24)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return Md5String(base64.URLEncoding.EncodeToString(b))
}
