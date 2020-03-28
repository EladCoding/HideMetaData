package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)


// only client
func (client *Connection) receiveAsClient() {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message) // TODO fix double received issue
		if err != nil {
			client.socket.Close()
			break
		}
		if length > 0 {
			fmt.Println("client RECEIVED:\n" + string(message))
		}
	}
}

// main client
func connectToNewConnection(usersMap userInfoMap, connectionName string,
	connectedServersConnections connectionNameToClient, connectedServersPubkey connectionNameToPubkey) {
	addConnectionSocket(usersMap, connectedServersConnections, connectionName)
	addConnectionPubKey(usersMap, connectedServersPubkey, connectionName)
}

// main client
func addConnectionSocket(usersMap userInfoMap, connectedServersConnections connectionNameToClient, connectionName string) {
	if _, ok := connectedServersConnections[connectionName]; !ok {
		serverAddress := usersMap[connectionName][AddressSpot]
		connectionSocket, err := net.Dial("tcp", serverAddress)
		checkErr(err)
		newConnection := &Connection{socket: connectionSocket}
		connectedServersConnections[connectionName] = newConnection
		go newConnection.receiveAsClient()
	}
}

// main client
func addConnectionPubKey(usersMap userInfoMap, connectedServersPubkey connectionNameToPubkey, connectionName string) {
	if _, ok := connectedServersPubkey[connectionName]; !ok {
		serverPubKeyPath := usersMap[connectionName][PublicKeyPathSpot]
		connectionPubKey := ReadPublicKeyFromFile(serverPubKeyPath)
		connectedServersPubkey[connectionName] = connectionPubKey
	}
}

// only client
func chooseNewChannel(curPath []string) string {
	lastChannel := curPath[len(curPath)-1]
	rand.Seed(time.Now().Unix())
	channelOptions := deleteValue(MediatorNames, lastChannel)
	n := rand.Int() % len(channelOptions)
	return channelOptions[n]
}

// only client
func createCipherPathMessage(message string, destination string, manager ConnectionsManager) ([]byte, string) {
	cipherMessage := []byte(message)
	prevChannel := destination
	curChannel := destination
	curPath := make([]string, 0)

	for i := 0; i < PathLen; i += 1 {
		curPath = append(curPath, curChannel)
		if i > 0 {
			cipherMessage = append([]byte(prevChannel), cipherMessage...)
		}
		addConnectionPubKey(manager.usersMap, manager.connectedServersPubkey, curChannel)
		curChannelPubKey := manager.connectedServersPubkey[curChannel]
		cipherMessage = hybridEncryption(cipherMessage, curChannelPubKey)
		prevChannel = curChannel
		curChannel = chooseNewChannel(curPath)
	}
	return cipherMessage, prevChannel
}

// only client
func getServerNameFromUser() string {
	fmt.Println("what server you want to send your message? (currently 001 002 or 003)")
	stdinReader := bufio.NewReader(os.Stdin)
	serverName, _ := stdinReader.ReadString('\n')
	serverName = strings.TrimRight(serverName, "\n")
	return serverName
}

//only client
func getMessageFromUser(serverName string) string {
	fmt.Println("what is your message, for server " + serverName + "?")
	stdinReader := bufio.NewReader(os.Stdin)
	message, _ := stdinReader.ReadString('\n')
	return message
}

func startClientMode(myName string, usersMap userInfoMap) {
	fmt.Println("Starting client...")
	manager := createGeneralManager(usersMap, myName)
	for {
		time.Sleep(1 * time.Second)
		serverName := getServerNameFromUser()
		if _, ok := usersMap[serverName]; !ok {
			fmt.Println("The server does not exists!\n")
			continue
		}
		message := getMessageFromUser(serverName)
		cipherMessage, nextChannel := createCipherPathMessage(message, serverName, manager)
		if _, ok := manager.connectedServersConnections[nextChannel]; !ok {
			connectToNewConnection(usersMap, nextChannel, manager.connectedServersConnections, manager.connectedServersPubkey)
		}
		nextChannelConnection := manager.connectedServersConnections[nextChannel]
		nextChannelConnection.socket.Write(cipherMessage)
	}
}
