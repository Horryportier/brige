package client

import (
	"brige/app/event"
	"brige/app/msg"
	"bufio"
	"encoding/json"
	"log"
	"net"
	"os"
	"strconv"
)

var (
	msgch chan string = make(chan string, 10)
)

type ClientSetup struct {
	ServerIp   string              `json:"server_ip"`
	ServerPort int                 `json:"server_port"`
	Connected  []event.EventTriger `json:"connected"`
}

// BUG: adds one word at the time
func input() {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		if line != "" {
			msgch <- line
		}
	}
}

func LoadClientSetup(path string) (ClientSetup, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return ClientSetup{}, err
	}
	var setup ClientSetup
	err = json.Unmarshal(buf, &setup)
	if err != nil {
		return ClientSetup{}, err
	}

	return setup, nil
}

type Client struct {
	events []event.EventTriger
}

func New() Client {
	return Client{}
}

func (c Client) EventIds() []string {
	var names []string
	for _, event := range c.events {
		names = append(names, event.Name)
	}
	return names
}

func (c Client) Start(setup ClientSetup) error {
	c.events = setup.Connected
	log.Println("names: ", c.EventIds())
	conn, err := net.Dial("tcp", setup.ServerIp+":"+strconv.Itoa(setup.ServerPort))
	if err != nil {
		return err
	}
	defer conn.Close()
	// sending inital package
	inital_buf, err := json.Marshal(c.EventIds())
	if err != nil {
		return err
	}
	log.Println("sending inital data")
	m := append([]byte(strconv.Itoa(int(msg.Inital))), inital_buf...)
	_, err = conn.Write(append(m, byte('\n')))
	if err != nil {
		return err
	}

	log.Println("reading response")
	buf, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		log.Println("Error reading:", err.Error())
	}
	log.Println(string(buf))

	go input()
	for {
		select {
		case m := <-msgch:
			log.Println("msg: ", m)
			log.Println("writing msg")
			_, err = conn.Write([]byte(strconv.Itoa(int(msg.Event)) + m + "\n"))
			if err != nil {
				return err
			}

			buf := make([]byte, 1024)
			log.Println("reading response")
			_, err = conn.Read(buf)
			if err != nil {
				log.Println("Error reading:", err.Error())
			}
			log.Println(string(buf))
		}
	}
}
