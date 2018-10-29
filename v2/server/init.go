package server

import (
    "fmt"
    "net"
    "io"
    "os"
)

func CreateServer() {
    //1 创建listener，内部创建了socket，开始监听端口
    listener, err := net.Listen("tcp", "localhost:8000")
    if err != nil {
        fmt.Println("监听端口失败,error:", err.Error())
        os.Exit(0) //终止程序
    }

    //2 持续等待来自客户端的连接并对连接进行处理
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
        go handleClient(conn)

    }
}

func handleClient(conn net.Conn) {
    i := 0
    for {
        //1 创建长度为20字节的缓冲区
        headBuffer := make([]byte, 20)
        //headBuffer := make([]byte, 3)

        //2 开始从网络流中读取数据
        readSize, err := conn.Read(headBuffer)
        if err != nil {
            if err == io.EOF {
                return
            }
            fmt.Println("从网络流中读取数据失败,error:", err.Error())
        }
        fmt.Printf("第%d次接收到客户端发送的长度为%d的数据：%s\n", i, readSize, string(headBuffer))
        i++

        conn.Write([]byte("0123456789"))
    }
}