package main

import (
	"flag"

	"github.com/sargo/kodicast/server"
)

func main() {
	flag.Parse()

	server.Serve()
}
