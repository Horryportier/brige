package main

// NOTE:
// event based system

import (
	"brige/app/server"
	"errors"
	"log"
	"os"
	"strconv"
)

type Config struct {
	ServerConfig server.ServerConfig
}

func print_help() {
	println(os.Args[0] + " <HOST> <PORT>")
}

func parse_args() (Config, error) {
	if len(os.Args) != 3 {
		print_help()
		return Config{}, errors.New("incorect amount of arguments")
	}
	host := os.Args[1]
	port, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return Config{}, err
	}

	var config Config = Config{
		ServerConfig: server.ServerConfig{
			Host: host,
			Port: port,
		},
	}
	return config, nil
}

func main() {
	config, err := parse_args()
	if err != nil {
		log.Fatal(err)
	}
	server.Entry(config.ServerConfig)
}
