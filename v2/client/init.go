package client

import (
	"github.com/xyb-blog-example/tcp-server/v2/protocol"
	"net"
	"fmt"
	"time"
	"strconv"
)

func CreateClient() {
	//1 创建待连接的远程节点
	tcpRAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("创建TCP远程节点失败,error:" + err.Error())
		return
	}

	//2 创建待连接的本地节点
	//tcpLAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8001")
	//if err != nil {
	//	fmt.Println("创建TCP远程节点失败,error:" + err.Error())
	//	return
	//}

	//3 连接远程节点
	tcpConn, err := net.DialTCP("tcp", nil, tcpRAddr)
	if err != nil {
		fmt.Println("连接服务端失败,error:" + err.Error())
		return
	}

	//4 启动goroutine，处理服务端的输入数据
	go handleServer(tcpConn)

	//5 客户端先向服务端发起发一点数据，然后服务端进行相应
	for i := 0; i < 10; i++ {
		err := protocol.SendTcpMsg(tcpConn, []byte("客户端发送了第"+ strconv.Itoa(i) + "条数据"))
		if err != nil {
			fmt.Println("发送消息失败,error:", err)
		}
		time.Sleep(1 * time.Second)
	}
}

func handleServer(conn *net.TCPConn) {
	for {
		//1 接受远程发送的数据包
		bodyBuffer, err := protocol.RecDataPack(conn)
		if err != nil {
			fmt.Printf("读取失败，断开了连接, error:%s\n", err.Error())
			return
		}

		//2 打印接收到的数据包
		fmt.Printf("收到的数据包为：%s\n", string(bodyBuffer))
	}
	return
}