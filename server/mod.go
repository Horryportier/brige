package server

import (
	"brige/app/event"
	"brige/app/msg"
	"encoding/json"
	"log"
	"net"
	"os"
	"strconv"
)

var (
	server = Server{make(chan []event.Event), make(chan []Connection)}
)

type Connection map[net.Conn][]string

type Server struct {
	event_queue chan []event.Event
	connections chan []Connection
}

type ServerSetup struct {
	Host string
	Port int
}

func LoadServerSetup(path string) (ServerSetup, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return ServerSetup{}, err
	}
	var setup ServerSetup
	err = json.Unmarshal(buf, &setup)
	if err != nil {
		return ServerSetup{}, err
	}

	return setup, nil
}

func Start(config ServerSetup) error {

	l, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		return err
	}
	defer l.Close()

	log.Println("Listening on " + config.Host + ":" + strconv.Itoa(config.Port))

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	log.Println("conn:  " + conn.RemoteAddr().String())
	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			log.Println("Error reading:", err.Error())
		}

		msg_type, err := strconv.Atoi(string(buf[0]))
		if err != nil {
			log.Println("Error first byte not an number:", err.Error())
		}
		buf = buf[1:]

		conn.Write([]byte("Message received."))
		switch msg.MsgType(msg_type) {
		case msg.MsgErr:
			log.Println(string(buf))
			return
		case msg.Exit:
			conn.Close()
			return
		case msg.Echo:
			conn.Write(buf)
			break
		case msg.Event:
			// TODO:
			break
		default:
		}

		log.Print(string(buf))
	}
}
