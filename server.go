package main

import (
	"crypto/rsa"
	"fmt"
	"net"
)

type ConnectionsManager struct {
	connectedServersConnections connectionNameToClient
	connectedServersPubkey connectionNameToPubkey
	connections map[*Client]bool
	register    chan *Client
	unregister  chan *Client
	publicKey   *rsa.PublicKey
	privateKey  *rsa.PrivateKey
	usersMap userInfoMap
	myName string
}

func (manager *ConnectionsManager) sendMessageToConnection(connection *Client, message []byte) {
	select {
	case connection.data <- message:
	default:
		manager.terminateConnection(connection) // TODO check if needed
	}
}

func (manager *ConnectionsManager) terminateConnection(connection *Client) {
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

func answerAsServer(manager *ConnectionsManager, message []byte, sender *Client) {
	messageReceivedAns := append(MessageReceivedAnswer, message...)
	manager.sendMessageToConnection(sender, messageReceivedAns)
}

func (manager *ConnectionsManager) receive(client *Client, mediator bool) {
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

func (manager *ConnectionsManager) send(client *Client) {
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

func createKeys(myPublicKeyPath string) (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, pubkey := GenerateRsaKeyPair(RsaKeyBits)
	WritePublicKeyToFile(myPublicKeyPath, pubkey)
	return privkey, pubkey
}

func startServerMode(myName string, usersMap userInfoMap) {
	fmt.Println("Starting server...")
	clientsManager := createGeneralManager(usersMap, myName)
	myAddress := usersMap[myName][AddressSpot]
	listener, err := net.Listen("tcp", myAddress)
	checkErr(err)

	go clientsManager.start()
	for {
		connection, err := listener.Accept()
		checkErr(err)
		client := &Client{socket: connection, data: make(chan []byte)}
		clientsManager.register <- client
		go clientsManager.receive(client, false)
		go clientsManager.send(client)
	}
}
