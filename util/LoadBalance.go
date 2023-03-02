package util

import (
	"hash/crc32"
	"math/rand"
	"time"
)

// 使用随机数  实现  最简单的负载均衡

type HttpServer struct { //  目标server 类
	Host string
}

func NewHttpServer(host string) *HttpServer {
	return &HttpServer{Host: host}
}

type LoadBalance struct { //  负载均衡类
	Servers []*HttpServer
}

// 初始化负载均衡
func NewLoadBalance() *LoadBalance {
	return &LoadBalance{Servers: make([]*HttpServer, 0)}
}

// 添加服务
func (this *LoadBalance) AddServer(server *HttpServer) {
	this.Servers = append(this.Servers, server)
}

// 随机算法
func (this *LoadBalance) SelectForRand() *HttpServer {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(this.Servers))
	return this.Servers[index]
}

// 使用ip   hash  进行 负载均衡
func (this *LoadBalance) SelectByIpHash(ip string) *HttpServer {
	index := int(crc32.ChecksumIEEE([]byte(ip))) % len(this.Servers)
	return this.Servers[index]
}
