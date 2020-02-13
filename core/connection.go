package core

import (
	"PhalGo-TCP/core/api"
	"PhalGo-TCP/util"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

// IConnection实现
type Connection struct {

	// 当前connection属于哪个server
	TcpServer api.IServer

	// 当前连接的socket：套接字
	Conn *net.TCPConn

	// 连接的ID
	ConnId uint32

	// 当前的连接状态
	IsClose bool

	// 当前连接所绑定的业务方法API
	HandApi api.HandleFunc

	// 告知当前连接已退出的channel
	ExitChan chan bool

	// 无缓冲通道，用于读写routinue通信
	MsgChan chan []byte

	// 路由处理的方法
	//Router api.IRouter

	// 由消息管理来处理路由
	MsgHandler api.IMsgHandler

	// 链接属性集合
	property map[string]interface{}

	// 链接属性锁
	propertyLock sync.RWMutex
}

func NewConnection(server api.IServer, conn *net.TCPConn, connId uint32, handler api.IMsgHandler) *Connection {
	connection := &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnId:    connId,
		IsClose:   false,
		//Router:   router,
		MsgHandler: handler,
		ExitChan:   make(chan bool, 1),
		MsgChan:    make(chan []byte),
		property:   make(map[string]interface{}),
	}
	// 将当前的connection加进连接管理器中
	server.GetConnManger().AddConn(connection)
	return connection
}

// 启动连接
func (c *Connection) StartConnect() {
	fmt.Sprintln("[connection] con n:", c.ConnId, " Start")
	// 启动从当前连接读数据的业务
	go c.startRead()
	// 启动从当前连接写数据的业务
	go c.startWrite()

	// 执行用户创建的钩子函数
	c.TcpServer.CallOnConnStart(c)
}

// 连接的读业务
func (c *Connection) startRead() {
	fmt.Println("Reader GoRoutinue is running")
	defer fmt.Println("conId: ", c.ConnId, " reader is close! RemoteAddr is", c.GetRemoteAddr().String())
	defer c.StopConnect()
	for {
		//buffer := make([]byte, util.GlobalObj.MaxPackageSize)
		//_, err := c.Conn.Read(buffer)
		//if err != nil{
		//	fmt.Println("[error] recv error: ", err)
		//	continue
		//}
		//// 调用当前连接所绑定的api
		//if err := c.HandApi(c.Conn, buffer, num); err != nil{
		//	fmt.Println("[error] handle is error: ", err)
		//	break
		//}

		// 创建一个拆包解包对象
		pack := NewTlvPackage()

		// 1.从conn把包的header读出来
		headBuffer := make([]byte, pack.GetPackHeadLen())
		if _, err := io.ReadFull(c.Conn, headBuffer); err != nil {
			fmt.Println("read header error: ", err)
			break
		}
		// 开始解包,读出MsgLen，MsgId
		msg, err := pack.TlvUnSerialize(headBuffer)
		if err != nil {
			fmt.Println("TlvUnSerialize error: ", err)
		}
		var msgData []byte
		if msg.GetMsgLen() > 0 {
			// 有数据，从conn第二次读MsgData
			msgData = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.Conn, msgData); err != nil {
				fmt.Println("TlvUnSerialize error: ", err)
			}
		}
		msg.SetMsgData(msgData)

		// 路由
		req := Request{
			conn: c,
			msg:  msg,
		}
		// 执行注册的路由方法
		//go func(request api.IRequest) {
		//	c.Router.PreHandle(request)
		//	c.Router.Handle(request)
		//	c.Router.AfterHandle(request)
		//}(&req)
		if util.GlobalObj.WorkerPoolSize > 0 {
			// 开启工作池，加入任务队列中
			c.MsgHandler.SendReqToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgRouter(&req)
		}

	}
}

// 连接的写业务
func (c *Connection) startWrite() {
	fmt.Println("Writer GoRoutines is running")
	defer fmt.Println(c.GetRemoteAddr().String(), " : writer goRoutines is finished")

	// 阻塞的等待channel的消息，写给客户端
	for {
		select {
		case data := <-c.MsgChan:
			// 有数据要返回给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("[error] conn write error: ", err)
				return
			}
		case <-c.ExitChan:
			// reader告知要推出，write也退出
			return
		}

	}
}

// 将数据发送给远程客户端的封包办法
func (c *Connection) SendMsg(msgId uint32, msgData []byte) error {
	if c.IsClose == true {
		return errors.New("Connection Close ")
	}
	tp := TlvPackage{}
	binaryData, err := tp.TlvSerialize(NewMessage(msgId, msgData))
	if err != nil {
		fmt.Println("[error] TlvSerialize error: ", err)
		return errors.New("TlvSerialize Error")
	}
	//if _, err := c.Conn.Write(binaryData); err != nil{
	//	fmt.Println("[error] server write: ", err)
	//	return errors.New("Conn Write Error ")
	//}
	// 发送给msgChan
	c.MsgChan <- binaryData
	return nil
}

// 结束连接
func (c *Connection) StopConnect() {
	fmt.Sprintln("[close] conn:", c.ConnId, " Stop")
	if c.IsClose == true {
		return
	}
	c.IsClose = true

	// 执行用户测试的钩子函数
	c.TcpServer.CallOnConnStop(c)

	// 关闭socket连接
	err := c.Conn.Close()
	if err != nil {
		fmt.Println("[error] close conn error : ", err)
	}
	// 告知write
	c.ExitChan <- true
	// 关闭chan
	close(c.ExitChan)
	close(c.MsgChan)

	// 将当前连接从连接管理器中删除
	c.TcpServer.GetConnManger().DelConn(c)
	fmt.Println("connection ", c.GetConnId(), "stop , now connection exist ", c.TcpServer.GetConnManger().GetConnNum())
}

// 获取当前连接绑定的socket conn
func (c *Connection) GetTCPConnect() *net.TCPConn {
	return c.Conn
}

// 获取当前连接模块的id
func (c *Connection) GetConnId() uint32 {
	return c.ConnId
}

// 获取远程客户端的TCP状态和IP，port
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 设置链接属性
func (c *Connection)SetProperty(key string, value interface{}){
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

// 获取链接属性
func (c *Connection)GetProperty(key string)(interface{}, error){
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok{
		return value, nil
	}
	return nil, errors.New("No Exist Property ")
}

// 删除链接属性
func (c *Connection)DelProperty(key string){
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}