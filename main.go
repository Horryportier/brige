package main

// NOTE:
// event based system

import (
	"brige/app/client"
	"brige/app/server"
	"log"
	"os"
)

func main() {
	switch os.Args[1] {
	case "-s", "--server":
		var server_setup server.ServerSetup
		var err error
		if len(os.Args) == 3 {
			server_setup, err = server.LoadServerSetup(os.Args[2])
		} else {
			server_setup, err = server.LoadServerSetup("server_config.json")
		}
		if err != nil {
			log.Fatalln(err)
		}
		err = server.Start(server_setup)
		if err != nil {
			log.Fatalln(err)
		}
		return
	case "-c", "--client":
		var client_setup client.ClientSetup
		var err error
		if len(os.Args) == 3 {
			client_setup, err = client.LoadClientSetup(os.Args[2])
		} else {
			client_setup, err = client.LoadClientSetup("client_config.json")
		}
		if err != nil {
			log.Fatalln(err)
		}
		_ = client_setup
		//err = client.Start(client_setup)
		if err != nil {
			log.Fatalln(err)
		}
		return

	}
}
