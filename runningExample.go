package main

import (
	"crypto/elliptic"
	"encoding/gob"
	"github.com/EladCoding/HideMetaData/mixnet"
	"os"
	"time"
)


func runMixNetWithoutClient() {
	go mixnet.StartUser("server", "001")
	time.Sleep(100*time.Millisecond)
	go mixnet.StartUser("server", "002")
	time.Sleep(100*time.Millisecond)
	go mixnet.StartUser("server", "003")
	time.Sleep(100*time.Millisecond)
	go mixnet.StartUser("mediator", "103")
	time.Sleep(100*time.Millisecond)
	go mixnet.StartUser("mediator", "102")
	time.Sleep(100*time.Millisecond)
	go mixnet.StartUser("mediator", "101")
	for {
		continue
	}
	//time.Sleep(100*time.Millisecond)
	//mixnet.StartUser("client", "201")
}


func runOneNode(mode string, name string) {
	mixnet.StartUser(mode, name)
}


func main() {
	gob.Register(elliptic.P256())
	//convertStructToBytes()
	mode := os.Args[1]
	name := os.Args[2]
	runOneNode(mode, name)
	//runWholeMixNet()
	//runMixNetWithoutClient()
	//scripts.Main()
}
