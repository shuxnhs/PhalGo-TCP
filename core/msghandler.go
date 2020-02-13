package core

import (
	"PhalGo-TCP/core/api"
	"PhalGo-TCP/util"
	"errors"
	"fmt"
)

type MsgHandler struct {

	// 路由集合
	routers map[uint32] api.IRouter

	// 业务工作worker池的worker数量
	workPoolSize	uint32

	// worker取任务的消息队列(每个worker对应一个任务队列)
	taskQueue	[]chan api.IRequest
}

// 创建方法
func NewMsgHandler() * MsgHandler {
	return &MsgHandler{
		routers:      make(map[uint32]api.IRouter),
		workPoolSize: util.GlobalObj.WorkerPoolSize,
		taskQueue:    make([]chan api.IRequest, util.GlobalObj.WorkerPoolSize),
	}
}

func (m *MsgHandler)DoMsgRouter(request api.IRequest){
	if router,ok := m.routers[request.GetMsgId()]; ok{
		router.PreHandle(request)
		router.Handle(request)
		router.AfterHandle(request)
	}
}

// 添加路由
func (m *MsgHandler)AddRouter(msgId uint32, router api.IRouter) error{
	// 是否已经存在
	if _,ok := m.routers[msgId]; ok{
		return errors.New("Router Exist ")
	}
	m.routers[msgId] = router
	return nil
}


// 启动一个Worker工作池
func (m *MsgHandler)StartWorkerPool(){
	// 创建最多有WorkerPoolSize个worker
	for i := 0; i < int(util.GlobalObj.WorkerPoolSize); i++  {
		// 任务队列初始化
		fmt.Println("[msgHandler] worker ", i, " init")
		m.taskQueue[i] = make(chan api.IRequest, util.GlobalObj.MaxWorkerTaskNum)
		go m.startOneWorker(i, m.taskQueue[i])
	}
}

// 启动一个Worker工作流程(不暴露对外)
func (m *MsgHandler)startOneWorker(workerId int, taskQueue chan api.IRequest)  {
	fmt.Println("worker ", workerId, " begin to work")
	// 不断的阻塞等待任务队列的消息
	for  {
		select {
			// 有消息过来，就去执行当前的request所绑定的业务
			case request := <- taskQueue :
				m.DoMsgRouter(request)
		}
	}
}

// 将消息添加到任务队列
func (m *MsgHandler)SendReqToTaskQueue(request api.IRequest)  {
	// 使用ConnID平均轮训分配（todo：分布式环境可以改为根据IP分配）
	workId := request.GetConnection().GetConnId() % m.workPoolSize
	fmt.Println("request message ", request.GetMsgId(), " add to taskQueue ", workId)
	// 将消息分配到对应的worker
	m.taskQueue[workId]<- request
}

