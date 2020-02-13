package core

import "PhalGo-TCP/core/api"

type Request struct {
	// 已经和客户端建立好的链接
	conn api.IConnection

	// 客户请求的数据
	msg api.IMessage

}

func (r *Request)GetConnection() api.IConnection{
	return r.conn
}

func (r *Request)GetMsgData() []byte{
	return r.msg.GetMsgData()
}

func (r *Request)GetMsgId() uint32 {
	return r.msg.GetMsgId()
}

