package api

// 数据包TLV的序列化抽象模块——处理粘包问题
type ITlvPackage interface {

	//获取包头长度的方法
	GetPackHeadLen() uint32

	// 封包
	TlvSerialize(IMessage) ([]byte, error)

	// 拆包
	TlvUnSerialize([]byte) (IMessage, error)
}
