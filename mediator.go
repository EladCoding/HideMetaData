package main

import (
	"fmt"
	"net"
)

func startMediatorMode(myAddress string, nextConnectionAddress string) {
	fmt.Println("Starting mediator...")
	connection, err := net.Dial("tcp", nextConnectionAddress)
	if err != nil {
		fmt.Println(err)
	}
	client := &Client{socket: connection, data: make(chan []byte)}
	go client.receive()


	mediatorManager := ClientManager{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	mediatorManager.clients[client] = true
	go mediatorManager.send(client)
	startServerMode(myAddress, true, client, mediatorManager)
}
