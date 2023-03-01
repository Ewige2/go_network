package util

import (
	"io"
	"net"
	"net/http"
	"time"
)

// 拷贝 http 头信息

func CloneHeader(src http.Header, dst *http.Header) {
	for k, v := range src {
		dst.Set(k, v[0])
	}
}
func RequestUrl(w http.ResponseWriter, r *http.Request, url string) {
	newreq, _ := http.NewRequest(r.Method, url, r.Body)
	CloneHeader(r.Header, &newreq.Header)
	// 添加转发地址
	newreq.Header.Add("x-forwarded-for", r.RemoteAddr)

	// 详细设置
	dt := http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		ResponseHeaderTimeout: 1 * time.Second,
	}

	newressponse, _ := dt.RoundTrip(newreq)
	getHeader := w.Header()
	//  将代理服务器的响应头，给客户端
	CloneHeader(newressponse.Header, &getHeader)
	// 写入   http  status
	w.WriteHeader(newressponse.StatusCode)
	defer newressponse.Body.Close()
	res_cont, _ := io.ReadAll(newressponse.Body)
	// 写入响应给客户端
	w.Write(res_cont)

}
