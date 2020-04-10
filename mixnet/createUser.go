package mixnet

import (
	"fmt"
	"strings"
)

func StartUser(mode string, name string) {
	switch strings.ToLower(mode) {
	case "server":
		StartServer(name)
	case "mediator":
		switch name {
		case "101":
			StartMediator(name, "102", true, false)
		case "102":
			StartMediator(name, "103", false, false)
		case "103":
			StartMediator(name, "001", false, true)
		default:
			fmt.Println("You can only chose mediator 101 102 or 103")
		}
	case "client":
		StartClient(name)
	default:
		fmt.Println("You can only be client mediator or server")
	}
}
