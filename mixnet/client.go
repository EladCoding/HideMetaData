package mixnet

import (
	"bufio"
	"crypto/ecdsa"
	"fmt"
	"net/rpc"
	"os"
	"strings"
	"sync"
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
func getMessageFromUser(serverName string) [][]byte {
	fmt.Println("what is your message, for server " + serverName + "?")
	stdinReader := bufio.NewReader(os.Stdin)
	message, _ := stdinReader.ReadString('\n')
	message = strings.TrimRight(message, "\n")
	return convertStringToMessages(message)
}


func convertStringToMessages(message string) [][]byte {
	splittedMsg := make([][]byte, 0)
	var endMsgSpot int
	for startMsgSpot := 0; startMsgSpot < len(message); startMsgSpot += maxUserMsgSize {
		if startMsgSpot + maxUserMsgSize < len(message) {
			endMsgSpot = startMsgSpot + maxUserMsgSize
		} else {
			endMsgSpot = len(message)
		}
		unPaddedMsg := message[startMsgSpot:endMsgSpot]
		paddedMsg, err := pkcs7padding([]byte(unPaddedMsg),MsgBytes)
		CheckErrToLog(err)
		splittedMsg = append(splittedMsg, paddedMsg)
	}
	return splittedMsg
}


func sendSpecificMessage(name string, serverName string, msg []byte, client *rpc.Client, wg *sync.WaitGroup) {
	cipherMsg, symKeys := createOnionMessage(name, serverName, msg, MediatorNames)
	var reply EncryptedMsg
	err := client.Call("CoordinatorListener.GetMessageFromClient", cipherMsg, &reply)
	if len(reply) == 0 {
		return
	}
	CheckErrToLog(err)
	for index, _ := range symKeys {
		symKey := symKeys[len(symKeys)-index-1]
		reply, err = symmetricDecryption(reply, symKey)
		CheckErrAndPanic(err)
	}
	replyMsg := ConvertBytesToReplyMsg(reply)
	orgMsg, _ := pkcs7strip(msg, MsgBytes)
	if string(replyMsg.Data) == string(orgMsg) {
		wg.Done()
	} else {
		panic(fmt.Sprintf("Reply don't feet.\nmsg: %v\nreply: %v\n", string(orgMsg), string(replyMsg.Data)))
	}
	//log.Printf("Client %v got reply message:\nFrom: %v, Data: %v\n", name, replyMsg.From, string(replyMsg.Data))
}


func StartClient(name string, statistics bool, serverNamePipe chan string, massagePipe chan string, donePipe chan bool) {
	var serverName string
	var userSplittedMsg [][]byte
	fmt.Printf("Starting Client %v...\n", name)
	mediatorAddress := UserAddressesMap["101"]
	client, err := rpc.Dial("tcp", mediatorAddress)
	CheckErrToLog(err)
	wg := &sync.WaitGroup{}
	for {
		if statistics {
			if <-donePipe {
				wg.Wait()
				break
			}
			serverName = <-serverNamePipe
		} else {
			serverName = getServerNameFromUser()
		}
		if serverName == "exit" {
			fmt.Printf("Exiting!\n")
			return
		} else if !StringInSlice(serverName, ServerNames) {
			fmt.Printf("Server %s does not exists!\n", serverName)
			continue
		}
		if len(serverName) != UserNameLen {
			fmt.Printf("Server %s size is weird!\n", serverName)
			continue
		}
		if statistics {
			userSplittedMsg = convertStringToMessages(<-massagePipe)
		} else {
			userSplittedMsg = getMessageFromUser(serverName)
		}
		for _, msg := range userSplittedMsg {
			if statistics {
				wg.Add(1)
				go sendSpecificMessage(name, serverName, msg, client, wg)
			} else {
				sendSpecificMessage(name, serverName, msg, client, nil)
			}
		}
		if !statistics {
			time.Sleep(100 * time.Millisecond)
		}
	}
	donePipe <- true
}
