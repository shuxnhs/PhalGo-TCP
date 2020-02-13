package main

import (
	"PhalGo-TCP/core"
	"PhalGo-TCP/core/api"
	"fmt"
)

// 自定义路由

type MyRouter struct {
	core.BaseRouter
}

//func (mr *MyRouter)PreHandle(request api.IRequest){
//	fmt.Println("call router preHandle...")
//	if _, err := request.GetConnection().GetTCPConnect().Write([]byte("before ping \n")); err !=nil{
//		fmt.Println("[error] preHandle error: ", err)
//	}
//}


// 处理conn业务的主方法
func (mr *MyRouter)Handle(request api.IRequest){
	fmt.Println("call router Handle...")
	//if _, err := request.GetConnection().GetTCPConnect().Write([]byte("ping \n")); err !=nil{
	//	fmt.Println("[error] Handle error: ", err)
	//}

	// 先读客户端数据
	fmt.Println("get msgId: ", request.GetMsgId(), ", msgData: ", string(request.GetMsgData()))

	// 发送数据
	if err := request.GetConnection().SendMsg(1, []byte("hello")); err != nil{
		fmt.Println("[error] ", err)
	}
}

type MyRouter2 struct {
	core.BaseRouter
}


// 处理conn业务的主方法
func (mr *MyRouter2)Handle(request api.IRequest){
	fmt.Println("call router Handle2...")
	// 先读客户端数据
	fmt.Println("get msgId: ", request.GetMsgId(), ", msgData: ", string(request.GetMsgData()))

	// 发送数据
	if err := request.GetConnection().SendMsg(2, []byte("world")); err != nil{
		fmt.Println("[error] ", err)
	}
}

// 处理conn业务之后的方法
//func (mr *MyRouter)AfterHandle(request api.IRequest){
//	fmt.Println("call router postHandle...")
//	if _, err := request.GetConnection().GetTCPConnect().Write([]byte("after ping \n")); err !=nil{
//		fmt.Println("[error] postHandle error: ", err)
//	}
//}

// conn连接前的钩子函数
func OnLineWarm(conn api.IConnection)  {
	fmt.Println("==========conn ", conn.GetConnId(), "is OnLine============")
	// 给客户端发送消息
	if err := conn.SendMsg(200, []byte("hello, sign")); err != nil{
		fmt.Println(err)
	}
	conn.SetProperty("country", "china")
}

func OffLineWarm(conn api.IConnection)  {
	fmt.Println("==========conn ", conn.GetConnId(), "is OffLine============")
	if value, err := conn.GetProperty("country"); err == nil{
		fmt.Println("client country is ", value)
	}else {
		fmt.Println("no exist country property")
	}
}

func main()  {

	// 创建一个server句柄，使用zinx的api
	s := core.NewServer()

	s.SetOnConnStart(OnLineWarm)
	s.SetOnConnStop(OffLineWarm)
	// 添加路由
	s.AddRouter(1, &MyRouter{})
	s.AddRouter(2, &MyRouter2{})

	// 启动
	s.Run()
	
}

