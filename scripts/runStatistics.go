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
	donePipe := make(chan bool)

	go mixnet.StartClient("201", true, serverNamePipe, messagesPipe, donePipe)
	time.Sleep(time.Second)
	startTime := time.Now()
	for i := 0; i < 1000; i += 1 {
		donePipe <- false
		serverNamePipe <- "001"
		msg := fmt.Sprintf("%v", i)
		messagesPipe <- msg
	}
	donePipe <- true
	if <- donePipe {
		fmt.Printf("Great.\n")
	} else {
		panic("What.\n")
	}
	duration := time.Since(startTime)
	fmt.Printf("Finished after: %v\n", duration)

}
