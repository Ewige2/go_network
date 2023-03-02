package util

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"sort"
	"time"
)

type HttpServers []*HttpServer

func (p HttpServers) Len() int           { return len(p) }
func (p HttpServers) Less(i, j int) bool { return p[i].CWeight > p[j].CWeight }
func (p HttpServers) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// 使用随机数  实现  最简单的负载均衡
type HttpServer struct { //  目标server 类
	Host      string
	Weight    int
	CWeight   int    //  当前权重
	Status    string //   服务状态,  默认UP
	FailCount int    //  记录  失败的次数
}

func NewHttpServer(host string, weight int, status string) *HttpServer {
	return &HttpServer{Host: host, Weight: weight, CWeight: weight, Status: status}
}

type LoadBalance struct { //  负载均衡类
	Servers  HttpServers
	CurIndex int
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
		if rad < v { //   逐一遍历 判断，
			fmt.Printf("rad: %d, v: %d", rad, v)
			return this.Servers[i]
		}

	}
	return this.Servers[0]
}

// 使用轮巡算法
func (this *LoadBalance) RoundRobin() *HttpServer {
	server := this.Servers[this.CurIndex]
	this.CurIndex = (this.CurIndex + 1) % len(this.Servers)
	if server.Status == "DOWN" {
		return this.RoundRobin()
	}
	return server
}

// 加权使用轮巡算法
func (this *LoadBalance) RoundRobinWeight() *HttpServer {
	server := this.Servers[ServerIndices[this.CurIndex]]
	this.CurIndex = (this.CurIndex + 1) % len(ServerIndices)
	return server
}

// 平滑加权使用轮巡算法
func (this *LoadBalance) RoundRobinWeight3() *HttpServer {
	for _, s := range this.Servers {
		s.CWeight = s.CWeight + s.Weight
	}
	sort.Sort(this.Servers)
	max := this.Servers[0] //  返回最大 作为命中服务
	max.CWeight = max.CWeight - SumWeight

	test := ""
	for _, s := range this.Servers {
		test += fmt.Sprintf("%d", s.CWeight)
	}
	fmt.Println(test)
	return max //  返回代理服务
}

var ServerIndices []int //  存储表示各个代理服务的权重信息

var LB *LoadBalance
var SumWeight int

func init() {

	LB = NewLoadBalance()
	LB.AddServer(NewHttpServer("http://localhost:9001", 5, "UP"))
	LB.AddServer(NewHttpServer("http://localhost:9002", 15, "UP"))
	LB.AddServer(NewHttpServer("http://localhost:9003", 15, "UP"))
	// 将地址信息的权重转化为， 对应比列存储
	for i, server := range LB.Servers {
		if server.Weight > 0 {
			for j := 0; j < server.Weight; j++ {
				ServerIndices = append(ServerIndices, i)
			}
		}
		// 计算权重总和
		SumWeight += server.Weight
	}
	// fmt.Println(ServerIndices)
	checkServers(LB.Servers)
}

func checkServers(servers HttpServers) {
	// 使用定时器
	t := time.NewTicker(time.Second * 3)
	check := NewHttpChecker(servers)
	for {
		select {
		case <-t.C:
			check.Check(time.Second * 2)
			for _, s := range servers {
				// 打印代理服务信息
				fmt.Println(s.Host, s.Status, s.FailCount)
			}
			fmt.Println("--------------------")

		}
	}

}
