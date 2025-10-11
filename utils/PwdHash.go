package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

func GetHashStr(str string) string {

	s := sha256.New()
	io.WriteString(s, str) //将str写入到s中
	bw := s.Sum(nil)       //w.Sum(nil)将w的hash转成[]byte格式

	return  hex.EncodeToString(bw) //将 bw 转成字符串
}

func CheckHashStr(str string,hashstr string) bool{
	return  GetHashStr(str) == hashstr
}

