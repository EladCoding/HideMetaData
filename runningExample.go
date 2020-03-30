package main

import "os"
import "github.com/EladCoding/HideMetaData/mainpkg"
import "fmt"

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Please provide mode and name.")
		return
	}
	go mainpkg.StartUser("server", "001")
	go mainpkg.StartUser("server", "002")
	go mainpkg.StartUser("server", "003")
	go mainpkg.StartUser("mediator", "101")
	go mainpkg.StartUser("mediator", "102")
	go mainpkg.StartUser("mediator", "103")
	mode := os.Args[1]
	name := os.Args[2]
	mainpkg.StartUser(mode, name)
}
