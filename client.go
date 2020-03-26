package main

import (
	"bufio"
	"crypto/rsa"
	"fmt"
	"net"
	"os"
	"strings"
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

func startClientMode(myName string, serverMap userInfoMap, mediatorMap userInfoMap) {
	fmt.Println("Starting client...")
	connectedServersConnections := make(connectionNameToClient)
	connectedServersPubkey := make(connectionNameToPubkey)
	var serverConnection *Client
	var serverPubKey *rsa.PublicKey
	var stdinReader *bufio.Reader

	for {
		fmt.Println("what server you want to send your message? (currently 001 002 or 003)")
		stdinReader = bufio.NewReader(os.Stdin)
		serverName, _ := stdinReader.ReadString('\n')
		serverName = strings.TrimRight(serverName, "\n")
		if _, ok := serverMap[serverName]; !ok {
			fmt.Println("The server does not exists!\n")
			continue
		}
		if _, ok := connectedServersConnections[serverName]; !ok {
			serverAddress := serverMap[serverName][AddressSpot]
			connectionSocket, err := net.Dial("tcp", serverAddress)
			checkErr(err)
			serverConnection = &Client{socket: connectionSocket}

			serverPubKeyPath := serverMap[serverName][PublicKeyPathSpot]
			serverPubKey = ReadPublicKeyFromFile(serverPubKeyPath)

			connectedServersConnections[serverName] = serverConnection
			connectedServersPubkey[serverName] = serverPubKey

			go serverConnection.receive()
		}

		serverConnection = connectedServersConnections[serverName]
		serverPubKey := connectedServersPubkey[serverName]

		fmt.Println("what is your message, for server " + serverName + "?")
		stdinReader = bufio.NewReader(os.Stdin)
		message, _ := stdinReader.ReadString('\n')
		cipherMessage := hybridEncryption([]byte(message), serverPubKey)
		fmt.Println("real msg:\n" + string(message))
		fmt.Printf("cipher msg: %x\n", cipherMessage)

		serverConnection.socket.Write(cipherMessage)
	}
}
