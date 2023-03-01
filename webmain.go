package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

/*
基本web  服务创建
*/

type webhandler struct {
}

func (webhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 获取验证头
	auth := r.Header.Get("Authorization")
	if auth == "" {
		w.Header().Set("WWW-Authenticate", "Basic realm= 请输入用户名和密码")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// 判断用户名or  密码是否正确
	auth_list := strings.Split(auth, " ")
	if len(auth_list) == 2 && auth_list[0] == "Basic" {
		res, err2 := base64.StdEncoding.DecodeString(auth_list[1])
		if err2 == nil && string(res) == "hali:123" {
			w.Write([]byte("跳转访问成功"))
			return
		}
	}
	w.Write([]byte("用户名or密码错误"))
}

type web2handler struct {
}

func (web2handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("跳转网站，web2"))
}

func main() {

	// 监听 退出
	c := make(chan os.Signal)
	// 创建两个协程
	go (func() {
		http.ListenAndServe(":9001", webhandler{})
	})()

	go (func() {
		http.ListenAndServe(":9002", web2handler{})
	})()

	signal.Notify(c, os.Interrupt)
	s := <-c
	log.Println(s)
}
