package scripts

import (
	"fmt"
	"github.com/EladCoding/HideMetaData/mixnet"
	"time"
)


func spamMixNet(clientName string, serverName string, numberOfMsgs int,	durationPipe chan time.Duration) {

	serverNamePipe := make(chan string)
	messagesPipe := make(chan string)
	clientDonePipe := make(chan bool)
	go mixnet.StartClient(clientName, true, true, false, serverNamePipe, messagesPipe, clientDonePipe, nil)
	time.Sleep(time.Second)

	startTime := time.Now()
	for i := 0; i < numberOfMsgs; i += 1 {
		clientDonePipe <- false
		serverNamePipe <- serverName
		msg := fmt.Sprintf("%v", i)
		messagesPipe <- msg
		time.Sleep(2*time.Millisecond)
	}
	clientDonePipe <- true
	if <- clientDonePipe {
		fmt.Printf("Statistics: FinishSpamming.\n")
	} else {
		panic("What.\n")
	}
	durationPipe <- time.Since(startTime)
}


func sendNiceMsgs(clientName string, serverName string, numberOfMsgs int, durationPipe chan time.Duration) {

	serverNamePipe := make(chan string)
	messagesPipe := make(chan string)
	clientDonePipe := make(chan bool)
	clientDurationPipe := make(chan time.Duration)
	go mixnet.StartClient(clientName, true, false, true, serverNamePipe, messagesPipe, clientDonePipe, clientDurationPipe)
	time.Sleep(time.Second)

	totalDuration := time.Duration(0)
	var curMsgLatency time.Duration
	for i := 0; i < numberOfMsgs; i += 1 {
		clientDonePipe <- false
		serverNamePipe <- serverName
		msg := fmt.Sprintf("%v", i)
		messagesPipe <- msg
		curMsgLatency = <- clientDurationPipe
		totalDuration += curMsgLatency
		fmt.Printf("Statistics: curMsg latency: %v\n", curMsgLatency)
	}
	clientDonePipe <- true
	fmt.Printf("Statistics: Finish latency..\n")
	durationPipe <- totalDuration / time.Duration(numberOfMsgs)
}


func RunStatistics() {
	CreateUsersMap()
	go runMixNetWithoutClients()
	time.Sleep(time.Second)

	spamClientName := "201"
	niceClientName := "202"
	serverName := "001"
	numberOfSpamMsgs := 1000000
	numberOfNiceMsgs := 100
	spammingDurationPipe := make(chan time.Duration)
	latencyDurationPipe := make(chan time.Duration)
	go spamMixNet(spamClientName, serverName, numberOfSpamMsgs, spammingDurationPipe)
	time.Sleep(time.Second)

	go sendNiceMsgs(niceClientName, serverName, numberOfNiceMsgs, latencyDurationPipe)

	latencyDuration := <-latencyDurationPipe
	spammingDuration := <- spammingDurationPipe

	fmt.Printf("----------Statistics----------\nFinished sending:\n%v msgs\n%v fakeMsgs (mean) each round\nafter: %v\nlatencyDuration is: %v\n",
		numberOfSpamMsgs, mixnet.FakeMsgsLaplaceMean, spammingDuration, latencyDuration)
}
