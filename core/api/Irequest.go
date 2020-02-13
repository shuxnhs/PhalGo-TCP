package api

// 将客户端的请求的链接信息和数据包装到一个request
type IRequest interface {

	GetConnection() IConnection

	GetMsgData() []byte

	GetMsgId() uint32

}
