package main

import (
	"github.com/devhoodit/sse-chat/server"
)

func main() {
	r := server.SetupRouter()

	r.Run()
}
