package mixnet

import (
	"bufio"
	"crypto/ecdsa"
	"fmt"
	"github.com/EladCoding/HideMetaData/scripts"
	"log"
	"net/rpc"
	"os"
	"strings"
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
	for {
		fmt.Println("what is your message, for server " + serverName + "?")
		stdinReader := bufio.NewReader(os.Stdin)
		message, _ := stdinReader.ReadString('\n')
		message = strings.TrimRight(message, "\n")
		if msgLen := len([]byte(message)); msgLen > MsgBytes {
			fmt.Printf("Message len is too long (%d). max len is %d\n", msgLen, MsgBytes)
		} else {
			return message
		}
	}
}


func StartClient(name string) {
	fmt.Printf("Starting Client %v...\n", name)
	mediatorAddress := userAddressesMap["101"]
	client, err := rpc.Dial("tcp", mediatorAddress)
	scripts.CheckErrToLog(err)
	for {
		serverName := getServerNameFromUser()
		if serverName == "exit" {
			fmt.Printf("Exiting!\n")
			return
		} else if len(serverName) != UserNameLen {
			fmt.Printf("Server %s size is weird!\n", serverName)
			continue
		} else if !scripts.StringInSlice(serverName, scripts.ServerNames) {
			fmt.Printf("Server %s does not exists!\n", serverName)
			continue
		}
		msgData := getMessageFromUser(serverName)
		padMessage, err := pkcs7padding([]byte(msgData), MsgBytes)
		scripts.CheckErrToLog(err)
		scripts.CheckErrAndPanic(err)
		cipherMsg, symKeys := createOnionMessage(name, serverName, padMessage, scripts.MediatorNames)
		var reply scripts.EncryptedMsg
		err = client.Call("CoordinatorListener.GetMessage", cipherMsg, &reply)
		scripts.CheckErrToLog(err)
		for index, _ := range symKeys {
			symKey := symKeys[len(symKeys)-index-1]
			reply, err = symmetricDecryption(reply, symKey)
			scripts.CheckErrAndPanic(err)
		}
		replyMsg := ConvertBytesToReplyMsg(reply)
		log.Printf("Client %v got reply message:\nFrom: %v, Data: %v\n", name, replyMsg.From, string(replyMsg.Data))
		time.Sleep(100*time.Millisecond)
	}
}
