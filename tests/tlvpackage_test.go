package tests

import (
	net2 "PhalGo-TCP/core"
	"fmt"
	"io"
	"net"
	"testing"
)

func TestTlvUnSerialize(t *testing.T) {
	/**
	 * 模拟服务器
	 */
	listener, err := net.Listen("tcp", "127.0.0.1:9010")
	if err != nil {
		fmt.Println("server listen error: ", err)
		return
	}

	// 负责处理解包
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error: ", err)
			}
			go func(conn net.Conn) {
				// 定义一个包对象
				pack := net2.NewTlvPackage()
				for {
					// 1.从conn把包的header读出来
					headBuffer := make([]byte, pack.GetPackHeadLen())
					if _, err := io.ReadFull(conn, headBuffer); err != nil {
						fmt.Println("read header error: ", err)
						break
					}
					// 开始解包,读出MsgLen，MsgId
					msgHeader, err := pack.TlvUnSerialize(headBuffer)
					if err != nil {
						fmt.Println("TlvUnSerialize error: ", err)
						return
					}
					if msgHeader.GetMsgLen() > 0 {
						// 有数据，从conn第二次读MsgData
						msg := msgHeader.(*net2.Message)
						msg.MsgData = make([]byte, msg.GetMsgLen())
						if _, err := io.ReadFull(conn, msg.MsgData); err != nil {
							fmt.Println("TlvUnSerialize error: ", err)
							return
						}
						fmt.Println("Recv MsgId: ", msg.MsgId, " MsgLen: ", msg.MsgLen, " data: ", string(msg.MsgData))
					}
				}
			}(conn)
		}
	}()

	/**
	 * 模拟客户端
	 */
	conn, err := net.Dial("tcp", "127.0.0.1:9010")
	if err != nil {
		fmt.Println("client connect error: ", err)
	}
	// 创建一个粘包
	tp := net2.NewTlvPackage()
	// 创建一个包
	msg1 := &net2.Message{
		MsgId:   1,
		MsgData: []byte{'h', 'e', 'l', 'l', 'o'},
		MsgLen:  5,
	}
	msgdata1, err := tp.TlvSerialize(msg1)
	if err != nil {
		fmt.Println("TlvSerialize msg1 error: ", err)
		return
	}
	// 创建第二个包
	msg2 := &net2.Message{
		MsgId:   2,
		MsgData: []byte{'w', 'o', 'r', 'l', 'd', '!'},
		MsgLen:  6,
	}
	msgdata2, err := tp.TlvSerialize(msg2)
	if err != nil {
		fmt.Println("TlvSerialize msg2 error: ", err)
		return
	}
	msgdata1 = append(msgdata1, msgdata2...)
	_, _ = conn.Write(msgdata1)
	select {} // 阻塞

}
