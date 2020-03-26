package main

import (
	"fmt"
	"net"
)

func splitMessageToNameAndMessage(message []byte) (string, []byte) {
	name := string(message[:UserNameLen])
	message = message[UserNameLen:]
	return name, message
}

func answerAsMediator(manager *ConnectionsManager, message []byte, sender *Client) {
	nextConnectionName, nextMessage := splitMessageToNameAndMessage(message)
	if _, ok := manager.connectedServersConnections[nextConnectionName]; !ok {
		createNewConnection(manager.usersMap, nextConnectionName, manager.connectedServersConnections, manager.connectedServersPubkey)
	}
	manager.connectedServersConnections[nextConnectionName].socket.Write(nextMessage)
}

func startMediatorMode(myName string, usersMap userInfoMap) {
	fmt.Println("Starting mediator...")
	manager := createGeneralManager(usersMap, myName)
	myAddress := usersMap[myName][AddressSpot]
	listener, err := net.Listen("tcp", myAddress)
	checkErr(err)
	go manager.start()
	for {
		connection, err := listener.Accept()
		checkErr(err)
		client := &Client{socket: connection, data: make(chan []byte)}
		manager.register <- client
		go manager.receive(client, true)
		go manager.send(client)
	}
}
