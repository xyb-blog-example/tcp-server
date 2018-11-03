package client

import (
	"testing"
	"log"
)

func TestCreateClient(t *testing.T) {
	conn := CreateClient()
	if conn == nil {
		log.Fatalln("与服务端建立连接失败")
	}
	log.Println("成功建立连接")
}

func TestSendMsgToServer(t *testing.T) {
	conn := CreateClient()
	if conn == nil {
		log.Fatalln("与服务端建立连接失败")
	}
	SendMsgToServer("第一次", conn)
	SendMsgToServer("第二次", conn)
	SendMsgToServer("第三次", conn)
	SendMsgToServer("第四次", conn)
	conn.Close()
}
