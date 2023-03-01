package main

import "encoding/base64"

func main() {
	// 测试 base64  加解密
	s := "hali:123"
	println(base64.StdEncoding.EncodeToString([]byte(s)))

	// 解密
	decodeString, _ := base64.StdEncoding.DecodeString("aGFsaToxMjM=")
	println(decodeString)

}
