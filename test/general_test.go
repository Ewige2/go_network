package test

import (
	"fmt"
	"gopkg.in/ini.v1"
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
