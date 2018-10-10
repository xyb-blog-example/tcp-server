package client

import (
	"net"
	"fmt"
	"time"
)

func CreateClient() {
	//1 创建待连接的远程节点
	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("创建TCP远程节点失败,error:" + err.Error())
		return
	}

	//2 连接远程节点
	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("连接服务端失败,error:" + err.Error())
		return
	}

	//3 启动goroutine，处理服务端的输入数据
	go handleServer(tcpConn)

	//4 客户端先向服务端发起发一点数据，然后服务端进行相应
	for i := 0; i < 10; i++ {
		text := fmt.Sprintf("%d + %d = ?", i, i)
		fmt.Println("向服务端发送了以下数据:", text)
		tcpConn.Write([]byte(text))
		time.Sleep(5 * time.Second)
	}
}

func handleServer(conn *net.TCPConn) {
	i := 0
	for {
		headBuffer		:= make([]byte, 20)
		_, err	:= conn.Read(headBuffer)
		if err != nil {
			fmt.Println("从网络流中读取数据失败,error:", err.Error())
		}
		fmt.Printf("第%d次接收到服务端发送的数据：%s\n", i, string(headBuffer))
		i++
	}
	return
}