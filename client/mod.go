package client

import (
	"brige/app/event"
	"encoding/json"
	"os"
)

type ClientSetup struct {
	ServerIp   string
	ServerPort int
}

type Client struct {
	events []event.Event
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
