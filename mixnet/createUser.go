package mixnet

import (
	"fmt"
	"github.com/EladCoding/HideMetaData/scripts"
	"strings"
)

func StartUser(mode string, name string) {
	switch strings.ToLower(mode) {
	case "server":
		for _, serverName := range scripts.ServerNames {
			if name == serverName {
				StartServer(name)
			}
			break
		}
		fmt.Println("Server name does not Exists!")
	case "mediator":
		switch name {
		case "101":
			StartMediator(name, 1, "102", true, false)
		case "102":
			StartMediator(name, 2, "103", false, false)
		case "103":
			StartMediator(name, 3, "001", false, true)
		default:
			fmt.Println("Mediator name does not Exists!")
		}
	case "client":
		for _, serverName := range scripts.ClientNames {
			if name == serverName {
				StartClient(name)
			}
			break
		}
		fmt.Println("Client name does not Exists!")
	default:
		fmt.Println("You can only be client mediator or server")
	}
}
