package util

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"time"
)

// 使用随机数  实现  最简单的负载均衡

type HttpServer struct { //  目标server 类
	Host   string
	Weight int
}

func NewHttpServer(host string, weight int) *HttpServer {
	return &HttpServer{Host: host, Weight: weight}
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

// 加权随机算法（如果用户输入的  权重很大， 将会极大的消耗内存）
func (this *LoadBalance) SelectByWeigthRand() *HttpServer {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(ServerIndices))
	return this.Servers[ServerIndices[index]]
}

// 加权随机算法（改良版本）
func (this *LoadBalance) SelectByWeigthRand2() *HttpServer {
	rand.Seed(time.Now().UnixNano())
	sumList := make([]int, len(this.Servers))
	sum := 0
	// 算出出权重信息
	for i := 0; i < len(this.Servers); i++ {
		sum += this.Servers[i].Weight
		sumList[i] = sum
	}
	rad := rand.Intn(sum) //  [)
	for i, v := range sumList {
		if rad < v {
			fmt.Printf("rad: %d, v: %d", rad, v)
			return this.Servers[i]
		}

	}
	return this.Servers[0]
}

var ServerIndices []int //  存储表示各个代理服务的权重信息

var LB *LoadBalance

func init() {
	LB = NewLoadBalance()
	LB.AddServer(NewHttpServer("http://localhost:9001", 5))
	LB.AddServer(NewHttpServer("http://localhost:9002", 15))
	// 将地址信息的权重转化为， 对应比列存储
	for i, server := range LB.Servers {
		if server.Weight > 0 {
			for j := 0; j < server.Weight; j++ {
				ServerIndices = append(ServerIndices, i)
			}
		}

	}
	fmt.Println(ServerIndices)
}
