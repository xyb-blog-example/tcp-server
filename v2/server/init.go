package server

import (
    "fmt"
    "net"
    "os"
	"tcp-server/v2/protocol"
	"time"
	"strconv"
)

func CreateServer() {
    //1 创建listener，内部创建了socket，开始监听端口
    listener, err := net.Listen("tcp", "localhost:8000")
    if err != nil {
        fmt.Println("监听端口失败,error:", err.Error())
        os.Exit(0) //终止程序
    }

    //2 持续等待来自客户端的连接并对连接进行处理
	i := 0
    for {
        //2.1 listener.Accept是一个阻塞方法，只有客户端连接时才会有返回
        fmt.Println("正在等待客户端的连接...")
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("接受连接失败,error:", err.Error())
            continue
        }
        fmt.Println("接收到了一个客户端的连接")

        //2.2 开goroutine对连接上来的客户端进行处理
        go handleClient(i, conn)
        i++
    }
}

func handleClient(i int, conn net.Conn) {
    for {
        //1 接受远程发送的数据包
        bodyBuffer, err := protocol.RecDataPack(conn)
        if err != nil {
			fmt.Printf("读取失败，断开了第%d个连接,error:%s\n", i, err.Error())
        	return
		}

        //2 打印接收到的数据包
        fmt.Printf("收到的数据包为：%s\n", string(bodyBuffer))

        //3 睡一秒以后回一个包
        time.Sleep(1 * time.Second)
        protocol.SendTcpMsg(conn, []byte("服务器" + strconv.Itoa(i) + "号收到了发送的数据"))
    }
}
