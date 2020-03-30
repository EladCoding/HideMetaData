package mainpkg

import (
	"fmt"
	"net"
)


func (manager *ConnectionsManager) sendMessageToConnection(connection *Connection, message []byte) {
	select {
	case connection.data <- message:
	default:
		manager.terminateConnection(connection) // TODO check if needed
	}
}

func (manager *ConnectionsManager) terminateConnection(connection *Connection) {
	close(connection.data)
	delete(manager.connections, connection)
	fmt.Println("A connection has terminated!\n%v", connection)
}

func (manager *ConnectionsManager) start() {
	for {
		select {
		case connection := <-manager.register:
			manager.connections[connection] = true
			fmt.Println("Added new connection!\n%v", connection)
			connection.data <- ConnectionSuccessfulAnswer

		case connection := <-manager.unregister:
			if _, ok := manager.connections[connection]; ok {
				manager.terminateConnection(connection) // TODO add reason
			}
		}
	}
}

func answerAsServer(manager *ConnectionsManager, message []byte, sender *Connection) {
	messageReceivedAns := append(MessageReceivedAnswer, message...)
	manager.sendMessageToConnection(sender, messageReceivedAns)
}

func (manager *ConnectionsManager) receiveAsServer(client *Connection, mediator bool) {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		if err != nil {
			manager.unregister <- client
			client.socket.Close()
			break
		}
		message = hybridDecryption(message[:length], manager.privateKey)
		if length > 0 {
			fmt.Println(manager.myName + " RECEIVED:\n" + string(message))
			if mediator {
				answerAsMediator(manager, message, client)
			} else {
				answerAsServer(manager, message, client)
			}
		}
	}
}

func (manager *ConnectionsManager) send(client *Connection) {
	defer client.socket.Close()
	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				return
			}
			client.socket.Write(message)
		}
	}
}

func startServerMode(myName string, usersMap userInfoMap, mediator bool) {
	fmt.Println("Starting Server...")
	manager := createGeneralManager(usersMap, myName)
	myAddress := usersMap[myName][AddressSpot]
	listener, err := net.Listen("tcp", myAddress)
	checkErr(err)
	go manager.start()

	for {
		connection, err := listener.Accept()
		checkErr(err)
		client := &Connection{socket: connection, data: make(chan []byte)}
		manager.register <- client
		go manager.receiveAsServer(client, mediator)
		go manager.send(client)
	}
}
