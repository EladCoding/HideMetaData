package main

import (
	"fmt"
)

func splitMessageToNameAndMessage(message []byte) (string, []byte) {
	name := string(message[:UserNameLen])
	message = message[UserNameLen:]
	return name, message
}

func answerAsMediator(manager *ConnectionsManager, message []byte, sender *Connection) {
	nextConnectionName, nextMessage := splitMessageToNameAndMessage(message)
	if _, ok := manager.connectedServersConnections[nextConnectionName]; !ok {
		connectToNewConnection(manager.usersMap, nextConnectionName, manager.connectedServersConnections, manager.connectedServersPubkey)
	}
	manager.connectedServersConnections[nextConnectionName].socket.Write(nextMessage)
}

func startMediatorMode(myName string, usersMap userInfoMap) {
	fmt.Println("Starting mediator...")
	startServerMode(myName, usersMap, true)
}
