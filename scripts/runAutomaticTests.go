package scripts

import (
	"fmt"
	"github.com/EladCoding/HideMetaData/mixnet"
	"time"
)

// send a specific test message to a specific server, and validate that the message return successfully and anonymously.
func sendTestingMsg(msg string, serverName string, clientDonePipe chan bool, serverNamePipe chan string,
	messagesPipe chan string, receivedMsgPipe chan string) bool {
	clientDonePipe <- false
	serverNamePipe <- serverName
	messagesPipe <- msg
	receivedMsg := <- receivedMsgPipe
	if receivedMsg != msg {
		fmt.Printf(
			"Test Failed.\n" +
			"Original massage: %v\n" +
			"Recived massage: %v",
			msg, receivedMsg)
		return false
	} else {
		return true
	}
}

// send a few messages to a some servers, and validate that the messages return successfully and anonymously.
func testMixNet(clientName string, serverName string, numberOfMsgs int, testSucceededPipe chan bool) {
	serverNamePipe := make(chan string)
	messagesPipe := make(chan string)
	receivedMessagesPipe := make(chan string)
	clientDonePipe := make(chan bool)
	testSucceeded := true

	go mixnet.StartClient(clientName, true, false, false, false, serverNamePipe, messagesPipe, clientDonePipe, nil, receivedMessagesPipe)
	time.Sleep(time.Second)

	for i := 0; i < numberOfMsgs; i += 1 {
		testSucceeded = sendTestingMsg(clientName + string(i), serverName, clientDonePipe, serverNamePipe, messagesPipe, receivedMessagesPipe)
		if testSucceeded == false {
			break
		}
	}

	clientDonePipe <- true
	if ! (<- clientDonePipe) {
		panic("What.\n")
	}

	testSucceededPipe <- testSucceeded
}

// Run the automatic tests.
func RunAutomaticTests() {
	mixnet.RoundSlotTime = time.Second
	CreateNodesMap()
	go runMixNetWithoutClients()
	time.Sleep(2*time.Second)

	firstClientName := "201"
	secondClientName := "202"
	serverName := "001"
	numberOfTestMsgs := 10
	firstTestSuccededPipe := make(chan bool)
	secondTestSuccededPipe := make(chan bool)

	go testMixNet(firstClientName, serverName, numberOfTestMsgs, firstTestSuccededPipe)
	time.Sleep(time.Second)

	go testMixNet(secondClientName, serverName, numberOfTestMsgs, secondTestSuccededPipe)

	firstTestSucceded := <- firstTestSuccededPipe
	secondTestSucceded := <- secondTestSuccededPipe

	if firstTestSucceded && secondTestSucceded{
		fmt.Printf("--------------------Automatic Tests--------------------\n" +
			"Test completed succesfully\n")
	}
}
