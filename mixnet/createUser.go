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
			StartCoordinator(name, 1, "102")
		case "102":
			StartMediator(name, 2, "103")
		case "103":
			StartDistributor(name, 3)
		default:
			fmt.Println("Mediator name does not Exists!")
		}
		return
	case "client":
		for _, clientName := range scripts.ClientNames {
			if name == clientName {
				StartClient(name)
				return
			}
		}
		fmt.Println("Client name does not Exists!")
	default:
		fmt.Println("You can only be client mediator or server")
	}
}
