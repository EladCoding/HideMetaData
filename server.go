package main

import (
	"fmt"
	"net"
)

type ClientManager struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
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
				mediatorManager.sendMessageToConnection(nextConnection, message)
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

func startServerMode(myAddress string, mediator bool, nextConnection *Client, mediatorManager ClientManager) {
	if !mediator {
		fmt.Println("Starting server...")
	}
	listener, err := net.Listen("tcp", myAddress)
	if err != nil {
		fmt.Println(err)
		// TODO handle
	}
	manager := ClientManager{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go manager.start()
	for {
		connection, err := listener.Accept()
		if err != nil { // TODO check i err (and nil) is good here (or ok)
			fmt.Println(err)
		}
		client := &Client{socket: connection, data: make(chan []byte)}
		manager.register <- client
		go manager.receive(client, mediator, nextConnection, mediatorManager)
		go manager.send(client)
	}
}
