package main

import (
	"fmt"
	"strings"
)

func startUser(mode string, name string) {
	usersMap := getUsersMap()
	switch strings.ToLower(mode) {
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
