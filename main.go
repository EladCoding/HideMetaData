package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	mode := os.Args[1]
	name := os.Args[2]
	flagMode := flag.String("mode", mode, "start in client or server mode")
	flag.Parse()
	mediatorMap := getMediatorMap()
	serversMap := getServerMap()
	switch strings.ToLower(*flagMode) {
	case "server":
		nilManager := ConnectionsManager{ // TODO check how to fix that (we don't use this arg)
			connections: make(map[*Client]bool),
			register:    make(chan *Client),
			unregister:  make(chan *Client),
		}
		startServerMode(name, serversMap, mediatorMap, false, nil, nilManager)
	case "client":
		startClientMode(name, serversMap, mediatorMap)
	//case "mediator":
	//	startMediatorMode(name, serversMap, mediatorMap)
	default:
		fmt.Println("You can only be client mediator or server")
	}
}
