package main

import (
	"fmt"
	//"github.com/EladCoding/HideMetaData/mainpkg"
	"github.com/EladCoding/HideMetaData/mixnet"
	"os"
)


//func main() {
//	if len(os.Args) < 3 {
//		fmt.Println("Please provide mode and name.")
//		return
//	}
//	go mainpkg.StartUser("server", "001")
//	go mainpkg.StartUser("server", "002")
//	go mainpkg.StartUser("server", "003")
//	go mainpkg.StartUser("mediator", "101")
//	go mainpkg.StartUser("mediator", "102")
//	go mainpkg.StartUser("mediator", "103")
//	mode := os.Args[1]
//	name := os.Args[2]
//	mainpkg.StartUser(mode, name)
//}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Please provide mode and name.")
		return
	}
	mode := os.Args[1]
	name := os.Args[2]
	if mode == "server" {
		mixnet.StartServer(name)
	}
	if mode == "mediator" {
		mixnet.StartMediator(name)
	}
	if mode == "client" {
		mixnet.StartClient(name)
	}
	//go mixnet.StartUser("server", "001")
	//go mixnet.StartUser("server", "002")
	//go mixnet.StartUser("server", "003")
	//go mixnet.StartUser("mediator", "101")
	//go mixnet.StartUser("mediator", "102")
	//go mixnet.StartUser("mediator", "103")
	//mode := os.Args[1]
	//name := os.Args[2]
	//mixnet.StartUser(mode, name)
}
