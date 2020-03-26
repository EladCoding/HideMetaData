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
	usersMap := getUsersMap()
	switch strings.ToLower(*flagMode) {
	case "server":
		startServerMode(name, usersMap)
	case "client":
		startClientMode(name, usersMap)
	case "mediator":
		startMediatorMode(name, usersMap)
	default:
		fmt.Println("You can only be client mediator or server")
	}
}
