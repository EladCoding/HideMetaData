package main

import (
	"crypto/rsa"
	"fmt"
	"net"
)

type ConnectionsManager struct {
	connections map[*Client]bool
	register    chan *Client
	unregister  chan *Client
	publicKey   *rsa.PublicKey
	privateKey  *rsa.PrivateKey
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

func (manager *ConnectionsManager) receive(client *Client, mediator bool, nextConnection *Client, mediatorManager ConnectionsManager) {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		fmt.Printf("msg: %x\n", message)
		if !mediator {
			fmt.Printf("len: %v\n", length)
			message = hybridDecryption(message[:length], manager.privateKey)
		}
		if err != nil {
			manager.unregister <- client
			client.socket.Close()
			break
		}
		if length > 0 {
			fmt.Println("RECEIVED:\n" + string(message))
			messageReceivedAns := append(MessageReceivedAnswer, message...)
			manager.sendMessageToConnection(client, messageReceivedAns)
			// TODO take only address part, and send to the next one
			if mediator {
				mediatorManager.sendMessageToConnection(nextConnection, message[:length])
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

func createConnctionsManager(privkey *rsa.PrivateKey, pubkey *rsa.PublicKey) ConnectionsManager {
	manager := ConnectionsManager{
		connections: make(map[*Client]bool),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		publicKey:   pubkey,
		privateKey:  privkey,
	}
	return manager
}

func startServerMode(myName string, serverMap userInfoMap, mediatorMap userInfoMap,  mediator bool, nextConnection *Client, mediatorManager ConnectionsManager) {
	// TODO check if mediator
	fmt.Println("Starting server...")
	myAddress := serverMap[myName][AddressSpot]
	myPublicKeyPath := serverMap[myName][PublicKeyPathSpot]
	privkey, pubkey := createKeys(myPublicKeyPath)
	clientsManager := createConnctionsManager(privkey, pubkey)

	listener, err := net.Listen("tcp", myAddress)
	checkErr(err)

	go clientsManager.start()
	for {
		connection, err := listener.Accept()
		checkErr(err)
		client := &Client{socket: connection, data: make(chan []byte)}
		clientsManager.register <- client
		go clientsManager.receive(client, mediator, nextConnection, mediatorManager)
		go clientsManager.send(client)
	}
}
