package main

// NOTE:
// brige is server to server task dispatcher
// IDEA:
// two servers talking to each other and sending information and commands.

import (
	"brige/app/server"
)

func main() {
	server.Entry(server.ServerConfig{Host: "localhost", Port: 8000})
	println("Hello World!")
}
