package api

// 定义消息模块的抽象层
type IMessage interface {

	// getter方法
	GetMsgId()		uint32
	GetMsgData()	[]byte
	GetMsgLen()		uint32

	// setter方法
	SetMsgId(uint32)
	SetMsgData([]byte)
	SetMsgLen(uint32)

}
