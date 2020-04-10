package mixnet

import (
	"bufio"
	"crypto/ecdsa"
	"fmt"
	"github.com/EladCoding/HideMetaData/scripts"
	"log"
	"net/rpc"
	"os"
)


type OnionMessage struct {
	From string
	To string
	PubKeyForSecret ecdsa.PublicKey
	Data []byte
}


// only client
func getServerNameFromUser() string {
	fmt.Print("what server do you want to send your message? (currently 001 002 or 003):")
	var serverName string
	fmt.Scanln(&serverName)
	return serverName
}

//only client
func getMessageFromUser(serverName string) string {
	fmt.Println("what is your message, for server " + serverName + "?")
	stdinReader := bufio.NewReader(os.Stdin)
	message, _ := stdinReader.ReadString('\n')
	return message
}


// only client
func createOnionMessage(name string, serverName string, msgData string) OnionMessage {
	var curPubKey ecdsa.PublicKey
	var onionMsg OnionMessage

	curOnionData := []byte(msgData)
	hopesArr := append(scripts.MediatorNames, serverName)
	for index, _ := range hopesArr {
		if index > 0 {
			curOnionData = ConvertOnionMsgToBytes(onionMsg)
		}
		curHop := hopesArr[len(hopesArr)-index-1]
		curOnionData, curPubKey = hybridEncription(curOnionData, curHop)
		onionMsg = OnionMessage{
			name, // TODO check what about from
			serverName,
			curPubKey,
			curOnionData,
		}
	}
	return onionMsg
}

func StartClient(name string) {
	fmt.Printf("Starting Client %v...\n", name)
	mediatorAddress := userAddressesMap["101"]
	fmt.Printf("name: %v. mediatorAddress: %v\n", name, mediatorAddress)
	client, err := rpc.Dial("tcp", mediatorAddress)
	scripts.CheckErrToLog(err)
	for {
		serverName := getServerNameFromUser()
		if !scripts.StringInSlice(serverName, scripts.ServerNames) {
			fmt.Printf("Server %s does not exists!\n", serverName)
			continue
		}
		if len(serverName) != UserNameLen {
			fmt.Printf("Server %s size is weird!\n", serverName)
			continue
		}
		msgData := getMessageFromUser(serverName)
		scripts.CheckErrAndPanic(err)
		cipherMsg := createOnionMessage(name, serverName, msgData)
		var reply Reply
		err = client.Call("MediatorListener.GetMessage", cipherMsg, &reply)
		scripts.CheckErrToLog(err)
		log.Printf("Reply: %v, From: %v, Data: %v\n", reply, reply.From, reply.Data)
	}
}
