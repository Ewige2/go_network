package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
)

/*
基本web  服务创建
*/

type webhandler struct {
}

func (webhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("跳转网站，web1"))
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
