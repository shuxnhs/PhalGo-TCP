package api

// 多路由及消息处理模块抽象
type IMsgHandler interface {

	// 调度执行对应的router
	DoMsgRouter(request IRequest)

	// 添加路由
	AddRouter(msgId uint32, router IRouter) error

	// 启动一个Worker工作池
	StartWorkerPool()

	// 将消息添加到任务队列
	SendReqToTaskQueue(request IRequest)

}
