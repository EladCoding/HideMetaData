package main

import "os"

func main() {
	go startUser("server", "001")
	go startUser("server", "002")
	go startUser("server", "003")
	go startUser("mediator", "101")
	go startUser("mediator", "102")
	go startUser("mediator", "103")
	mode := os.Args[1]
	name := os.Args[2]
	startUser(mode, name)
}
