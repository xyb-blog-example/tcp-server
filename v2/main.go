package main

import (
	"github.com/xyb-blog-example/tcp-server/v2/client"
	"github.com/xyb-blog-example/tcp-server/v2/server"
	"os"
)

func main() {
	params := ""
	if len(os.Args) < 2 {
		params = "client"
	} else {
		params = os.Args[1]
	}

	if params == "client" {
		client.CreateClient()
	} else if params == "server" {
		server.CreateServer()
	}
}

