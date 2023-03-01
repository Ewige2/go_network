package util

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
)

// 读取配置文件,并保存到map  中

type EnvConfig *os.File

var ProxyConfigs map[string]string

/* init
初始化不能采用初始化表达式初始化的变量。
程序运行前的注册。
实现sync.Once功能。
其他s
*/

func init() { // 在 调用这个文件的时候 会自动进行初始化
	ProxyConfigs = make(map[string]string)
	EnvConfig, err := ini.Load("config.ini")

	if err != nil {
		fmt.Println(err)
		return
	}
	// 获取配置文件信息
	proxy, _ := EnvConfig.GetSection("proxy")
	// 判断扇区是否为空
	if proxy != nil {
		// 获取 子扇区
		sections := proxy.ChildSections()
		for _, sec := range sections {
			url, _ := sec.GetKey("url")
			path, _ := sec.GetKey("path")
			// 如果不为空 将数据保存到map数组中
			ProxyConfigs[path.Value()] = url.Value()
			
		}

	}
}
