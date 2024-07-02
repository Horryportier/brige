package server

import (
	"brige/app/event"
	"brige/app/msg"
	"bufio"
	"io"

	//"brige/app/utils"
	"encoding/json"
	"log"
	"net"
	"os"
	"strconv"
)

type Connection struct {
	conn          net.Conn
	observing     []string
	current_event chan event.Event
}

type Server struct {
	event_queue chan event.Event
	connections []chan Connection
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

func New() Server {
	return Server{make(chan event.Event, 100), make([]chan Connection, 0)}
}

func (s Server) Start(setup ServerSetup) error {

	l, err := net.Listen("tcp", setup.Host+":"+strconv.Itoa(setup.Port))
	if err != nil {
		return err
	}
	defer l.Close()

	log.Println("Listening on " + setup.Host + ":" + strconv.Itoa(setup.Port))

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		log.Println("acceping connection")
		connection_info := make(chan Connection, 10)
		connection_info <- Connection{conn: conn, current_event: make(chan event.Event, 10)}
		s.connections = append(s.connections, connection_info)
		go handleRequest(connection_info, s)
		select {
		// TODO: ya that is blocking server bad idea
		case e := <-s.event_queue:
			for _, c := range s.connections {
				connection := <-c
				for _, name := range connection.observing {
					if e.Name == name {
						connection.current_event <- e
					}
				}
				c <- connection
			}
		}

	}
}

func handleRequest(connectionch chan Connection, server Server) {
	log.Println("handling request")
	connection := <-connectionch
	log.Println("conn:  " + connection.conn.RemoteAddr().String())
	for {
		buf, err := bufio.NewReader(connection.conn).ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				connectionch <- Connection{}
				close(connectionch)
				return
			}
			log.Fatal("Failed to read from conn: ", err)
		}

		msg_type, err := strconv.Atoi(string(buf[0]))
		if err != nil {
			log.Println("Error first byte not an number:", err.Error())
		}
		buf = buf[1:]

		switch msg.MsgType(msg_type) {
		case msg.MsgErr:
			log.Println(string(buf))
			return
		case msg.Inital:
			log.Println("reciving inital package")
			var names []string
			json.Unmarshal(buf, &names)
			connection.observing = names
			connectionch <- connection
			connection.conn.Write(append([]byte(strconv.Itoa(int(msg.Ok))), '\n'))
			break
		case msg.Exit:
			connection.conn.Close()
			return
		case msg.Echo:
			connection.conn.Write(buf)
			break
		case msg.Event:
			var event event.Event
			err := json.Unmarshal(buf, &event)
			if err != nil {
				log.Println("Can't unmarshal event:", err)
				continue
			}
			server.event_queue <- event
			break
		default:
		}
		select {
		case e := <-connection.current_event:
			buf, err := json.Marshal(e)
			if err != nil {
				log.Println("failed to send event:", err)
			}
			connection.conn.Write(append(buf, '\n'))

		}

		log.Print("msg: ", string(buf))
	}
}
