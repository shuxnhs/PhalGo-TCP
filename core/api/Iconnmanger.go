package api

// 链接管理的抽象层
type IConnManger interface {

	// 添加链接
	AddConn(connection IConnection)

	// 删除链接
	DelConn(connection IConnection)

	// 根据connId获取链接
	GetConn(connId uint32) (IConnection, error)

	// 获取所有链接的总数
	GetConnNum() int

	// 删除所有链接
	DelAllConn()

}
