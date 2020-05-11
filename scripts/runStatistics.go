package scripts

import (
	"fmt"
	"github.com/EladCoding/HideMetaData/mixnet"
	"time"
)


func sendAlittleMsg(i int, serverName string, clientDonePipe chan bool, serverNamePipe chan string, messagesPipe chan string) {
	msg := fmt.Sprintf("%v", i)
	clientDonePipe <- false
	serverNamePipe <- serverName
	messagesPipe <- msg
}


func spamMixNet(clientName string, serverName string, numberOfMsgs int,	durationPipe chan time.Duration) {

	serverNamePipe := make(chan string)
	messagesPipe := make(chan string)
	clientDonePipe := make(chan bool)
	go mixnet.StartClient(clientName, true, true, false, serverNamePipe, messagesPipe, clientDonePipe, nil)
	time.Sleep(time.Second)

	startTime := time.Now()
	nextRound := startTime
	for i := 0; i < numberOfMsgs; i += 1 {
		nextRound = nextRound.Add(1000 * time.Microsecond)
		time.Sleep(time.Until(nextRound))
		sendAlittleMsg(i, serverName, clientDonePipe, serverNamePipe, messagesPipe)
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
	time.Sleep(2*time.Second)

	spamClientName := "201"
	niceClientName := "202"
	serverName := "001"
	numberOfSpamMsgs := 100000
	numberOfNiceMsgs := 1
	spammingDurationPipe := make(chan time.Duration)
	latencyDurationPipe := make(chan time.Duration)
	go spamMixNet(spamClientName, serverName, numberOfSpamMsgs, spammingDurationPipe)
	time.Sleep(time.Second)

	go sendNiceMsgs(niceClientName, serverName, numberOfNiceMsgs, latencyDurationPipe)

	latencyDuration := <-latencyDurationPipe
	spammingDuration := <- spammingDurationPipe
	numberOfMsgsPerSecond := numberOfSpamMsgs / int(spammingDuration / time.Second)

	fmt.Printf("----------Statistics----------\nFinished sending:\n%v msgs\n%v fakeMsgs (mean) each round\n" +
		"after: %v\nlatencyDuration is: %v\n msgPerSecond (without fakes): %v\n",
		numberOfSpamMsgs, mixnet.FakeMsgsLaplaceMean, spammingDuration, latencyDuration, numberOfMsgsPerSecond)
}
