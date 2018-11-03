package client

import (
	"net"
	"fmt"
	"log"
)

func CreateClient() *net.TCPConn {
	//1 创建待连接的远程节点
	tcpRAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatalln("创建TCP远程节点失败,error:" + err.Error())
		return nil
	}

	//2 创建待连接的本地节点
	//tcpLAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8001")
	//if err != nil {
	//	log.Fatalln("创建TCP远程节点失败,error:" + err.Error())
	//	return nil
	//}

	//3 连接远程节点
	//tcpConn, err := net.DialTCP("tcp", tcpLAddr, tcpRAddr)
	tcpConn, err := net.DialTCP("tcp", nil, tcpRAddr)
	if err != nil {
		log.Fatalln("连接服务端失败,error:" + err.Error())
		return nil
	}

	//4 启动goroutine，处理服务端的输入数据
	go handleServer(tcpConn)
	return tcpConn
}

func handleServer(conn *net.TCPConn) {
	i := 0
	for {
		headBuffer		:= make([]byte, 100)
		_, err	:= conn.Read(headBuffer)
		if err != nil {
			log.Fatalln("从网络流中读取数据失败,error:", err.Error())
		}
		fmt.Printf("第%d次接收到服务端发送的数据：%s\n", i, string(headBuffer))
		i++
	}
	return
}

//给服务端发送信息
func SendMsgToServer(msg string, conn *net.TCPConn) {
	conn.Write([]byte(msg))
}