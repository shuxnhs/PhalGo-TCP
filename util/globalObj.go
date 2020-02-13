package util

import (
	"PhalGo-TCP/core/api"
	"encoding/json"
	"io/ioutil"
)

/**
 * 存储全局的所有变量配置的对象
 * 可以通过config.json来配置
 */

type GlobalObject struct {

	// server配置
	TcpServer  api.IServer // 全局server对象
	TcpHost    string      // 当前服务器监听的IP
	TcpPort    int         // 当前服务器监听的端口
	ServerName string      // 当前服务器名字

	// 框架配置
	Version        string // 框架版本号
	MaxConnNum     int    // 当前服务器主机允许的最大连接数
	MaxPackageSize uint32 // 数据包的最大值

	WorkerPoolSize   uint32 // 业务工作worker池大小
	MaxWorkerTaskNum uint32 // 每个worker处理的任务的最大值
}

var GlobalObj *GlobalObject

// 初始化GlobalObj对象
func init() {
	GlobalObj = &GlobalObject{
		TcpHost:          "0.0.0.0",
		TcpPort:          9009,
		ServerName:       "PhalGo-TCP-Server",
		Version:          "V1.0",
		MaxConnNum:       1000,
		MaxPackageSize:   1024,
		WorkerPoolSize:   10,
		MaxWorkerTaskNum: 1024,
	}

	// 尝试从配置文件中加载配置
	GlobalObj.ReloadConfig()
}

// 加载配置文件
func (g *GlobalObject) ReloadConfig() {
	config, err := ioutil.ReadFile("conf/config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(config, &GlobalObj)
	if err != nil {
		panic(err)
	}
}
