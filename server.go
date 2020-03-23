package main

import (
	"crypto/rsa"
	"fmt"
	"net"
)

type ClientManager struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	publicKey *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func (manager *ClientManager) sendMessageToConnection(connection *Client, message []byte) {
	select {
	case connection.data <- message:
	default:
		manager.terminateConnection(connection) // TODO check if needed
	}
}

func (manager *ClientManager) terminateConnection(connection *Client) {
	close(connection.data)
	delete(manager.clients, connection)
	fmt.Println("A connection has terminated!\n%v", connection)
}

func (manager *ClientManager) start() {
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = true
			fmt.Println("Added new connection!\n%v", connection)
			connection.data <- ConnectionSuccessfulAnswer

		case connection := <-manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				manager.terminateConnection(connection) // TODO add reason
			}
		}
	}
}

func (manager *ClientManager) receive(client *Client, mediator bool, nextConnection *Client, mediatorManager ClientManager) {
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

func (manager *ClientManager) send(client *Client) {
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

func createServerKeys(mediator bool) (*rsa.PrivateKey, *rsa.PublicKey) {
	var myPublicKeyPath string
	var myPrivateKeyPath string
	if mediator {
		myPublicKeyPath = MediatorPublicKeyPath
		myPrivateKeyPath = MediatorPrivateKeyPath
	} else {
		fmt.Println("Starting server...")
		myPublicKeyPath = ServerPublicKeyPath
		myPrivateKeyPath = ServerPrivateKeyPath
	}
	privkey, pubkey := GenerateRsaKeyPair(RsaKeyBits)
	WritePublicKeyToFile(myPublicKeyPath, pubkey)
	WritePrivateKeyToFile(myPrivateKeyPath, privkey)
	return privkey, pubkey
}

func startServerMode(myAddress string, mediator bool, nextConnection *Client, mediatorManager ClientManager) {
	privkey, pubkey := createServerKeys(mediator)

	listener, err := net.Listen("tcp", myAddress)
	checkErr(err)

	manager := ClientManager{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		publicKey: pubkey,
		privateKey: privkey,
	}
	go manager.start()
	for {
		connection, err := listener.Accept()
		checkErr(err)
		client := &Client{socket: connection, data: make(chan []byte)}
		manager.register <- client
		go manager.receive(client, mediator, nextConnection, mediatorManager)
		go manager.send(client)
	}
}
