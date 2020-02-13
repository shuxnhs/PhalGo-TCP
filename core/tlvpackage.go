package core

import (
	"PhalGo-TCP/core/api"
	"PhalGo-TCP/util"
	"bytes"
	"encoding/binary"
	"errors"
)

type TlvPackage struct {

}

// 实例化方法
func NewTlvPackage() *TlvPackage {
	return &TlvPackage{}
}

//获取包头长度的方法
func (t *TlvPackage) GetPackHeadLen() uint32{
	// MsgLen + MsgId = uint32(4个字节) + uint32(4个字节)
	return 8
}

// 封包
func (t *TlvPackage)TlvSerialize(msg api.IMessage) ([]byte, error){
	dataBuffer := bytes.NewBuffer([]byte{})
	// 将MsgLen写入buffer
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgLen()); err != nil{
		return nil, err
	}
	// 将MsgId写入buffer
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgId()); err != nil{
		return nil, err
	}
	// 将MsgData写入buffer
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgData() ); err != nil{
		return nil, err
	}
	return dataBuffer.Bytes(), nil
}

// 拆包
func (t *TlvPackage)TlvUnSerialize(data []byte) (api.IMessage, error){
	// 先读取Header的MsgLen，再根据MsgLen去读数据
	dataBuffer := bytes.NewReader(data)

	msg := &Message{}
	// 读MsgLen
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.MsgLen); err != nil{
		return nil, err
	}
	// 读MsgId
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.MsgId); err != nil{
		return nil, err
	}
	// 判断长度是否超过系统限制最大
	if util.GlobalObj.MaxPackageSize > 0 && util.GlobalObj.MaxPackageSize < msg.MsgLen{
		return nil, errors.New("Too Large MsgData ")
	}

	return msg, nil
}