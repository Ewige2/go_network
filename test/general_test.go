package test

import (
	"fmt"
	"gopkg.in/ini.v1"
	"hash/crc32"
	"log"
	"testing"
)

func TestInit(t *testing.T) {
	// 测试init 读取配置文件
	cfg, err := ini.Load("../config.ini")
	if err != nil {
		log.Println(err)
		return
	}
	sec, err := cfg.GetSection("proxy")
	fmt.Println(sec.GetKey("url"))
	fmt.Println(sec.GetKey("path"))

}

func TestIpHash(t *testing.T) {
	// 测试使用ip   hash    负载均衡
	ip := "127.0.0.1"
	fmt.Println(crc32.ChecksumIEEE([]byte(ip))) //    hash  之后对  size  进行取余
}
