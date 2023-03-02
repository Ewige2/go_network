package util

import (
	"net/http"
	"time"
)

type HttpChecker struct {
	Servers HttpServers
}

func NewHttpChecker(servers HttpServers) *HttpChecker {
	return &HttpChecker{Servers: servers}
}

func (this *HttpChecker) Check(timeout time.Duration) {
	// 使用一个http  客户端 每隔一段时间请求  一次 代理服务
	client := http.Client{}
	for _, server := range this.Servers {
		resp, err := client.Head(server.Host) //  使用呢Head  请求 就减少数据传输量
		if resp != nil {
			defer resp.Body.Close()
		}

		if err != nil { // 宕机了
			server.Status = "DOWN"
			continue
		}
		// 其他情况判断
		if resp.StatusCode >= 200 && resp.StatusCode < 400 {
			server.Status = "UP"
		} else {
			server.Status = "DOWN"
		}
	}

}
