package scripts

import (
	"fmt"
	"github.com/EladCoding/HideMetaData/mixnet"
	"log"
	"time"
)

// Send a specific dummy message to a specific server.
func sendSpammingMsg(i int, serverName string, clientDonePipe chan bool, serverNamePipe chan string,
	messagesPipe chan string) {
	msg := fmt.Sprintf("%v", i)
	clientDonePipe <- false
	serverNamePipe <- serverName
	messagesPipe <- msg
}

// Send a lot of dummy messages to a specific server, and check the throughput.
func spamMixNet(clientName string, serverName string, numberOfMsgs int, durationPipe chan time.Duration,
	slotDuration time.Duration) {
	serverNamePipe := make(chan string)
	messagesPipe := make(chan string)
	clientDonePipe := make(chan bool)
	go mixnet.StartClient(clientName, false, true, true, false, serverNamePipe, messagesPipe, clientDonePipe, nil, nil)
	time.Sleep(time.Second)

	startTime := time.Now()
	nextRound := startTime

	for i := 0; i < numberOfMsgs; i += 1 {
		time.Sleep(time.Until(nextRound))
		memUsage := mixnet.GetMemUsage()
		if memUsage > 0.73 {
			fmt.Printf("Memory Usage: %v\n", memUsage)
			time.Sleep(10*slotDuration)
		}
		sendSpammingMsg(i, serverName, clientDonePipe, serverNamePipe, messagesPipe)
		nextRound = nextRound.Add(slotDuration)
	}
	clientDonePipe <- true
	if <-clientDonePipe {
		log.Printf("Statistics: FinishSpamming.\n")
	} else {
		panic("What.\n")
	}
	durationPipe <- time.Since(startTime)
}

// Send a specific dummy message to a specific server, and check the goodput.
func sendNiceMsgs(clientName string, serverName string, numberOfMsgs int, durationPipe chan time.Duration) {
	serverNamePipe := make(chan string)
	messagesPipe := make(chan string)
	clientDonePipe := make(chan bool)
	clientDurationPipe := make(chan time.Duration)
	go mixnet.StartClient(clientName, false, true, false, true, serverNamePipe, messagesPipe, clientDonePipe, clientDurationPipe, nil)
	time.Sleep(time.Second)

	totalDuration := time.Duration(0)
	var curMsgLatency time.Duration
	for i := 0; i < numberOfMsgs; i += 1 {
		clientDonePipe <- false
		serverNamePipe <- serverName
		msg := fmt.Sprintf("%v", i)
		messagesPipe <- msg
		curMsgLatency = <-clientDurationPipe
		totalDuration += curMsgLatency
		log.Printf("Statistics: curMsg latency: %v\n", curMsgLatency)
	}
	clientDonePipe <- true
	log.Printf("Statistics: Finish latency..\n")
	durationPipe <- totalDuration / time.Duration(numberOfMsgs)
}

// Run mixnet statistics.
func RunStatistics() {
	CreateNodesMap()
	go runMixNetWithoutClients()
	time.Sleep(2 * time.Second)

	spamClientName := "201"
	niceClientName := "202"
	serverName := "001"
	numberOfSpamMsgs := 200000
	slotDuration := 280 * time.Microsecond
	roundDuration := mixnet.RoundSlotTime
	maxMsgsPerRound := int(roundDuration / slotDuration)
	minimumRounds := numberOfSpamMsgs / maxMsgsPerRound
	numberOfNiceMsgs := minimumRounds / 4
	if numberOfNiceMsgs < 1 {
		numberOfNiceMsgs = 1
	} else if numberOfNiceMsgs > 4 {
		numberOfNiceMsgs = 4
	}
	spammingDurationPipe := make(chan time.Duration)
	latencyDurationPipe := make(chan time.Duration)
	go spamMixNet(spamClientName, serverName, numberOfSpamMsgs, spammingDurationPipe, slotDuration)
	time.Sleep(time.Second)

	go sendNiceMsgs(niceClientName, serverName, numberOfNiceMsgs, latencyDurationPipe)

	latencyDuration := <-latencyDurationPipe
	spammingDuration := <-spammingDurationPipe
	numberOfMsgsPerSecond := (float64(numberOfSpamMsgs) / float64(spammingDuration) * float64(time.Second))
	msgsLimitPerSecond := int(time.Second / slotDuration)

	fmt.Printf("----------Statistics----------\n"+
		"Finished sending:\n"+
		"%v msgs\n"+
		"%v fakeMsgs (mean) each round\n"+
		"after: %v\n"+
		"latencyDuration is: %v\n"+
		"msgPerSecond (without fakes): %v\n"+
		"msgLimitPerSecond (without fakes): %v\n",
		numberOfSpamMsgs, mixnet.FakeMsgsLaplaceMean, spammingDuration, latencyDuration, numberOfMsgsPerSecond, msgsLimitPerSecond)
}
