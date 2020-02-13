package core

import (
	"PhalGo-TCP/core/api"
	"PhalGo-TCP/util"
	"fmt"
	"net"
)

// Iserver的接口实现
type Server struct {
	// 服务器的名称
	Name string
	// 服务器绑定的IP版本
	IpVersion string
	// 服务器监听的IP
	Ip string
	// 服务器监听的端口
	Port int
	//
	//// 暂时只能注册一个路由
	//Router api.IRouter
	// 由消息处理去选择路由
	MsgHandler api.IMsgHandler

	ConnManger api.IConnManger

	// 创建连接后自动运行的钩子函数
	OnConnStart func(conn api.IConnection)

	// 销毁连接后自动运行的钩子函数
	OnConnStop func(conn api.IConnection)
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server %s listening at IP: %s, Port : %d is starting \n", util.GlobalObj.ServerName, s.Ip, s.Port)

	go func() {
		// 先开辟一个worker工作池
		s.MsgHandler.StartWorkerPool()

		// 1. 获取一个TCP的Addr（创建套接字）
		addr, err := net.ResolveTCPAddr(s.IpVersion, fmt.Sprintf("%s:%d", s.Ip, s.Port))
		if err != nil {
			fmt.Println("resolveTcpAddr error : ", err)
			return
		}
		// 2. 监听服务器地址
		listener, err := net.ListenTCP(s.IpVersion, addr)
		if err != nil {
			fmt.Println("listenTcp error : ", err)
			return
		}
		fmt.Println("[server] server ", s.Name, "start success， listening")

		// 3. 阻塞的等待客服端连接，处理客户端的链接业务
		var cid uint32 = 0
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("[server] accept err :", err)
				continue
			}
			fmt.Println(cid)
			fmt.Println("client connect success")

			// 判断是否超出最大连接数
			if s.ConnManger.GetConnNum() > util.GlobalObj.MaxConnNum {
				fmt.Println("reach max conn")
				// todo: 给客户端发送错误包
				_ = conn.Close()
				continue
			}
			// 链接模块
			dealConnection := NewConnection(s, conn, cid, s.MsgHandler)
			go dealConnection.StartConnect()
			cid++

		}
	}()

}

func (s *Server) Stop() {
	// 将一些服务器的资源，状态进行回收或者停止
	s.ConnManger.DelAllConn()
}

func (s *Server) Run() {

	// 启动服务器
	s.Start()

	// todo：做一些启动后的功能

	// 阻塞状态
	select {}
}

// 初始化模块
func NewServer() api.IServer {
	server := &Server{
		Name:      util.GlobalObj.ServerName,
		IpVersion: "tcp4", // todo:暂时只实现tcp4
		Ip:        util.GlobalObj.TcpHost,
		Port:      util.GlobalObj.TcpPort,
		//Router:    nil,
		MsgHandler: NewMsgHandler(),
		ConnManger: NewConnManger(),
	}
	return server
}

//
//// 定义客户端连接后的handleApi
//// todo：暂时写死
//func CallBackFuncToClient(conn *net.TCPConn, buffer []byte, num int) error{
//	fmt.Println("[conn handle], callback to client")
//	if _,err := conn.Write(buffer[:num]);  err != nil{
//		fmt.Println("[error] write handle callback error: ", err)
//		return errors.New("CallBackFuncToClient Error")
//	}
//	return nil
//}

// 添加路由
func (s *Server) AddRouter(msgId uint32, router api.IRouter) {
	//s.Router = router
	if err := s.MsgHandler.AddRouter(msgId, router); err != nil {
		fmt.Println("[error]Add router error: ", err)
	}
	fmt.Println("[success]Add router success!")
}

// 获取当前的连接管理器
func (s *Server) GetConnManger() api.IConnManger {
	return s.ConnManger
}

// 创建连接后自动运行的钩子函数创建
func (s *Server)SetOnConnStart(hookFunc func(api.IConnection)){
	fmt.Println("Conn Start Hook Func Set")
	s.OnConnStart = hookFunc
}

// 销毁连接后自动运行的钩子函数创建
func (s *Server)SetOnConnStop(hookFunc func(api.IConnection)){
	fmt.Println("Conn Stop Hook Func Set")
	s.OnConnStop = hookFunc
}

// 创建连接后自动运行的钩子函数调用
func (s *Server)CallOnConnStart(conn api.IConnection){
	if s.OnConnStart != nil {
		fmt.Println("Conn Start Hook Func Call")
		s.OnConnStart(conn)
	}
}

// 销毁连接后自动运行的钩子函数调用
func (s *Server)CallOnConnStop(conn api.IConnection){
	if s.OnConnStop != nil {
		fmt.Println("Conn Stop Hook Func Call")
		s.OnConnStop(conn)
	}
}
