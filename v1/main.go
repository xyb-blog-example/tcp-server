package main

import (
	"os"
	"tcp-server/v1/client"
	"tcp-server/v1/server"
	"bufio"
)

func main() {
	params := ""
	if len(os.Args) < 2 {
		params = "client"
	} else {
		params = os.Args[1]
	}

	if params == "client" {
		conn := client.CreateClient()
		for {
			inputReader := bufio.NewReader(os.Stdin)
			input, _ := inputReader.ReadString('\n')
			client.SendMsgToServer(input, conn)
		}
	} else if params == "server" {
		server.CreateServer()
	}
}

