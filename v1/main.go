package main

import (
	"os"
	"tcp-server/v1/client"
	"tcp-server/v1/server"
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

