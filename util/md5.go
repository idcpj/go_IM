package util

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

//小写
func Md5Encode(data string)string{
	hash := md5.New()
	hash.Write([]byte(data))
	sum := hash.Sum(nil)
	return hex.EncodeToString(sum)
}
//大写
func MD5Encode(data string)string{
	return strings.ToUpper(Md5Encode(data))
}

func ValidatePasswd(plainPwd,salt,passpwd string) bool{
	return Md5Encode(plainPwd+salt)==passpwd
}
func MakePasswd(plainPwd,slat string)string{
	return Md5Encode(plainPwd+slat)
}