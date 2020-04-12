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
				return
			}
		}
		fmt.Println("Server name does not Exists!")
	case "mediator":
		switch name {
		case "101":
			StartCoordinator(name, 1, "103")
		//case "102":
		//	StartMediator(name, 2, "103", false, false)
		case "103":
			StartDistributor(name, 3)
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
