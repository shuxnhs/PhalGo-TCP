package main

import (
	"PhalGo-TCP/core"
	"fmt"
	"io"
	"net"
	"time"
)

func main()  {
	fmt.Println("client start")
	conn, err := net.Dial("tcp", "127.0.0.1:9010")
	if err != nil{
		fmt.Println("connect error,", err, "exit!")
		return
	}

	for{

		tp := core.TlvPackage{}
		binaryData, err :=  tp.TlvSerialize(core.NewMessage(2, []byte("client")))
		if err != nil{
			fmt.Println("[error] TlvSerialize error: ", err)
			break
		}
		if _, err := conn.Write(binaryData); err != nil{
			fmt.Println("[error] server write: ", err)
			break
		}


		// 读取服务端返回的数据
		// 1.从conn把包的header读出来
		headBuffer := make([]byte, tp.GetPackHeadLen())
		if _,err := io.ReadFull(conn, headBuffer); err != nil{
			fmt.Println("read header error: ", err)
			break
		}
		// 开始解包,读出MsgLen，MsgId
		msg, err := tp.TlvUnSerialize(headBuffer)
		if err != nil{
			fmt.Println("TlvUnSerialize error: ", err)
		}

		var msgData []byte
		if msg.GetMsgLen() > 0{
			// 有数据，从conn第二次读MsgData
			msgData = make([]byte, msg.GetMsgLen())
			if _,err := io.ReadFull(conn, msgData); err != nil{
				fmt.Println("TlvUnSerialize error: ", err)
			}
			msg.SetMsgData(msgData)
			fmt.Println("Recv MsgId: ", msg.GetMsgId(), " MsgLen: ", msg.GetMsgLen(), " data: ", string(msg.GetMsgData()))
		}

		time.Sleep(1 * time.Second)
	}

}
