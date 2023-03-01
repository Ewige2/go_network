package main

import (
	"go_network/util"
	"io"
	"log"
	"net/http"
)

/*
请求转发： 通过  访问  这个 网站的服务，该服务器 会 进行 跳转  访问另一个网站的服务
*/

func RequestUrl(w http.ResponseWriter, r *http.Request, url string) {
	nearreq, _ := http.NewRequest(r.Method, url, r.Body)
	util.CloneHeader(r.Header, &nearreq.Header)

	nearresponse, _ := http.DefaultClient.Do(nearreq) // 请求连接获取响应
	getHeader := w.Header()
	util.CloneHeader(nearresponse.Header, &getHeader) //  拷贝头
	w.WriteHeader(nearresponse.StatusCode)            //  写入响应码

	defer nearresponse.Body.Close()
	// 读取响应信息
	all, _ := io.ReadAll(nearresponse.Body)
	w.Write(all)
	return
	// 访问不是在请求转发列表的目录，返回的信息
	w.Write([]byte("default  index "))
}

type ProxyHandler struct {
}

func (p ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 进行错误处理
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(500)
			log.Println(err)
		}
	}()

	if r.URL.Path == "/a" {
		// http.Get()，  是基于  NewRequest() 方法进行创建的
		// 获取访跳转目标的页面信息
		RequestUrl(w, r, "http://127.0.0.1:9001")
		return
	} else if r.URL.Path == "/b" {
		RequestUrl(w, r, "http://127.0.0.1:9002")
		return
	}

}

func main() {
	// 监听服务
	http.ListenAndServe(":8080", ProxyHandler{})
}
