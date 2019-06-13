package util

import (
	"log"
	"testing"
)

func TestMd5Encode(t *testing.T){
	md5 :="123456"
	res :="e10adc3949ba59abbe56e057f20f883e"
	str := Md5Encode(md5)
	if str!=res {
		log.Fatal("md5函数错误")
	}
}