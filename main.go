package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	mode := os.Args[1]
	flagMode := flag.String("mode", mode, "start in client or server mode")
	flag.Parse()
	switch strings.ToLower(*flagMode) {
	case "server":
		nilManager := ClientManager{ // TODO check how to fix that (we don't use this arg)
			clients:    make(map[*Client]bool),
			register:   make(chan *Client),
			unregister: make(chan *Client),
		}
		startServerMode(ServerAddress, false, nil, nilManager)
	case "client":
		startClientMode(MediatorAddress)
	case "mediator":
		startMediatorMode(MediatorAddress, ServerAddress)
	default:
		fmt.Println("You can only be client mediator or server")
	}
}
