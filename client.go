package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type Client struct {
	socket net.Conn
	data   chan []byte
}

func (client *Client) receive() {
	for {
		message := make([]byte, 0, 4096)
		length, err := client.socket.Read(message) // TODO fix double received issue
		if err != nil {
			client.socket.Close()
			break
		}
		if length > 0 {
			fmt.Println("RECEIVED:\n" + string(message))
		}
	}
}

func startClientMode(connectionAddress string) {
	fmt.Println("Starting client...")
	connection, err := net.Dial("tcp", connectionAddress)
	checkErr(err)
	client := &Client{socket: connection}
	go client.receive()

	//mediatorPubKey := ReadPublicKeyFromFile(MediatorPublicKeyPath)
	serverPubKey := ReadPublicKeyFromFile(ServerPublicKeyPath)
	//serverPrivateKey := ReadPrivateKeyFromFile(ServerPrivateKeyPath)

	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		cipherMessage := hybridEncryption([]byte(message), serverPubKey)
		fmt.Println("real msg:\n" + string(message))
		fmt.Printf("cipher msg: %x\n",cipherMessage)
		//plainMessage := hybridDecryption(cipherMessage, serverPrivateKey)
		//fmt.Println("%x\n", plainMessage)
		connection.Write(cipherMessage)
	}
}
