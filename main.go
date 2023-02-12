package main

import (
	"github.com/PylonSchema/server/server"
)

func main() {
	r := server.SetupRouter()

	r.Run(":8080")
	// log.Fatal(autotls.Run(r, "127.0.0.1"))
}
