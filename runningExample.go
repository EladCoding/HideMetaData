package main

import (
	"crypto/elliptic"
	"encoding/gob"
	"github.com/EladCoding/HideMetaData/mixnet"
	"time"
)


func runMixNetWithoutClients() {
	go mixnet.StartUser("server", "001")
	time.Sleep(100*time.Millisecond)
	go mixnet.StartUser("server", "002")
	time.Sleep(100*time.Millisecond)
	go mixnet.StartUser("server", "003")
	time.Sleep(100*time.Millisecond)
	go mixnet.StartUser("mediator", "103")
	time.Sleep(200*time.Millisecond)
	go mixnet.StartUser("mediator", "102")
	time.Sleep(200*time.Millisecond)
	go mixnet.StartUser("mediator", "101")
	time.Sleep(200*time.Millisecond)
}


func runOneNode(mode string, name string) {
	mixnet.StartUser(mode, name)
}


func runClient() {
	mode := "client"
	name := "201"
	runOneNode(mode, name)
}


func main() {
	gob.Register(elliptic.P256())
	runMixNetWithoutClients()
	runClient()
}
