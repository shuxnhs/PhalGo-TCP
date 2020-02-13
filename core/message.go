package core

type Message struct {

	MsgId 	uint32	// 消息id
	MsgData	[]byte	// 消息内容
	MsgLen	uint32	// 消息长度

}

// 创建一个Message方法
func NewMessage(msgId uint32, msgData []byte) *Message {
	return &Message{
		MsgId:   msgId,
		MsgData: msgData,
		MsgLen:  uint32(len(msgData)),
	}
}

// getter方法
func (m *Message)GetMsgId()	uint32{
	return m.MsgId
}
func (m *Message)GetMsgData() []byte{
	return m.MsgData
}
func (m *Message)GetMsgLen() uint32{
	return m.MsgLen
}

// setter方法
func (m *Message)SetMsgId(MsgId uint32){
	m.MsgId = MsgId
}
func (m *Message)SetMsgData(MsgData []byte){
	m.MsgData = MsgData
}
func (m *Message)SetMsgLen(MsgLen uint32){
	m.MsgLen = MsgLen
}
