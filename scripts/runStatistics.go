package scripts

import (
	"fmt"
	"github.com/EladCoding/HideMetaData/mixnet"
	"time"
)

func RunStatistics() {
	CreateUsersMap()
	go runMixNetWithoutClients()
	time.Sleep(time.Second)

	serverNamePipe := make(chan string)
	messagesPipe := make(chan string)

	go mixnet.StartClient("201", true, serverNamePipe, messagesPipe)
	for i := 0; i < 100; i += 1 {
		serverNamePipe <- "001"
		msg := fmt.Sprintf("%v", i)
		messagesPipe <- msg
		fmt.Printf(msg)
		time.Sleep(time.Second / 10)
	}
	for {
		continue
	}
}
