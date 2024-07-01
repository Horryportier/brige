package server

import (
	"brige/app/msg"
	"fmt"
	"net"
	"os"
	"strconv"
)

type ServerConfig struct {
	Host string
	Port int
}

func Entry(config ServerConfig) {
	// Listen for incoming connections.
	l, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + config.Host + ":" + strconv.Itoa(config.Port))
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.

	for {
		buf := make([]byte, 1024)
		// Read the incoming connection into the buffer.
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
		msg_type, err := strconv.Atoi(string(buf[0]))
		if err != nil {
			fmt.Println("Error first byte not an number:", err.Error())
		}

		conn.Write([]byte("Message received."))
		switch msg.MsgType(msg_type) {
		case msg.Exit:
			conn.Close()
			return
		case msg.Echo:
			conn.Write(buf)
			break
		default:
		}
		println(string(buf))
	}
}
