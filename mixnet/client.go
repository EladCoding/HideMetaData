package mixnet

import (
	"bufio"
	"crypto/ecdsa"
	"fmt"
	"github.com/EladCoding/HideMetaData/scripts"
	"log"
	"net/rpc"
	"os"
	"time"
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
		cipherMsg, symKeys := createOnionMessage(name, serverName, []byte(msgData), scripts.MediatorNames)
		var reply scripts.EncryptedMsg
		err = client.Call("CoordinatorListener.GetMessage", cipherMsg, &reply)
		scripts.CheckErrToLog(err)
		for index, _ := range symKeys {
			symKey := symKeys[len(symKeys)-index-1]
			reply, err = symmetricDecryption(reply, symKey)
			scripts.CheckErrAndPanic(err)
		}
		replyMsg := ConvertBytesToReplyMsg(reply)
		log.Printf("Reply: From: %v, Data: %v", replyMsg.From, string(replyMsg.Data))
		time.Sleep(100*time.Millisecond)
	}
}
