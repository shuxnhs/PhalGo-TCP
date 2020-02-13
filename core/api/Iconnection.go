package api

import "net"

// 定义连接模块的抽象层
type IConnection interface {

	// 启动连接
	StartConnect()

	// 结束连接
	StopConnect()

	// 获取当前连接绑定的socket conn
	GetTCPConnect() *net.TCPConn

	// 获取当前连接模块的id
	GetConnId() uint32

	// 获取远程客户端的TCP状态和IP，port
	GetRemoteAddr() net.Addr

	// 发送数据，将数据发送给远程客户端
	SendMsg(msgId uint32, data []byte) error

	// 设置链接属性
	SetProperty(key string, value interface{})

	// 获取链接属性
	GetProperty(key string)(interface{}, error)

	// 删除链接属性
	DelProperty(key string)
}

// 定义一个处理连接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error