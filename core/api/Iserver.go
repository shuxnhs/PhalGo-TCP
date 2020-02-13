package api

// 定义服务器模块的抽象层
type IServer interface {
	// 启动服务器
	Start()
	// 停止服务器
	Stop()
	// 运行服务器
	Run()

	// 路由功能: 给当前服务注册一个路由方法，供客户端链接处理使用
	AddRouter(msgId uint32, router IRouter)

	// 获取当前的连接管理器
	GetConnManger() IConnManger

	// 创建连接后自动运行的钩子函数创建
	SetOnConnStart(func(IConnection))

	// 销毁连接后自动运行的钩子函数创建
	SetOnConnStop(func(IConnection))

	// 创建连接后自动运行的钩子函数调用
	CallOnConnStart(IConnection)

	// 销毁连接后自动运行的钩子函数调用
	CallOnConnStop(IConnection)
}
