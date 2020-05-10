package scripts

import (
	"github.com/EladCoding/HideMetaData/mixnet"
	"time"
)


func runMixNetWithoutClients() {
	go mixnet.StartUser("server", "001")
	time.Sleep(200*time.Millisecond)
	go mixnet.StartUser("server", "002")
	time.Sleep(200*time.Millisecond)
	go mixnet.StartUser("server", "003")
	time.Sleep(200*time.Millisecond)
	go mixnet.StartUser("mediator", "103")
	time.Sleep(200*time.Millisecond)
	go mixnet.StartUser("mediator", "102")
	time.Sleep(200*time.Millisecond)
	mixnet.StartUser("mediator", "101")
}


func runOneNode(mode string, name string) {
	mixnet.StartUser(mode, name)
}


func runClient() {
	mode := "client"
	name := "201"
	runOneNode(mode, name)
}
