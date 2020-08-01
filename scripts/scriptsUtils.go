package scripts

import (
	"github.com/EladCoding/HideMetaData/mixnet"
	"time"
)

// Run the mixnet architecture without clients.
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

// Run one node of the mixnet architecture.
func runOneNode(mode string, name string) {
	mixnet.StartUser(mode, name)
}

// Run a fixed client as a mixnet architecture node.
func runClient() {
	mode := "client"
	name := "201"
	runOneNode(mode, name)
}
