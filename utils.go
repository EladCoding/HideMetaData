package main

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
)

// built in messages
var ConnectionSuccessfulAnswer = []byte("Success!\nYou are now connected")
var MessageReceivedAnswer = []byte("Success!\nYour message has received\n")

// cipher vars
var RsaKeyBits = 2048
var CipherRsaLen = RsaKeyBits / 8
var AesKeyBytes = 32
var PathLen = 3

var working_dir, _ = os.Getwd()
// servers map
var ServerPublicKeyPathFormat = working_dir + "/keys/server%s/public_key.txt"
var ServerPrivateKeyPathFormat = working_dir + "/keys/server%s/private_key.txt"
var ServerHost = "localhost"
var ServerPortFormat = "9%s"
var ServerAddressFormat = ServerHost + ":" + ServerPortFormat
// mediators map
var MediatorPublicKeyPathFormat = working_dir + "/keys/mediator%s/public_key.txt"
var MediatorPrivateKeyPathFormat = working_dir + "/keys/mediator%s/private_key.txt"
var MediatorHost = "localhost"
var MediatorHostFormat = "8%s"
var MediatorAddressFormat = MediatorHost + ":" + MediatorHostFormat
// clients map
var ClientPublicKeyPathFormat = working_dir + "/keys/client%s/public_key.txt"
var ClientPrivateKeyPathFormat = working_dir + "/keys/client%s/private_key.txt"
var ClientHost = "localhost"
var ClientHostFormat = "7%s"
var ClientAddressFormat = ClientHost + ":" + ClientHostFormat

var MediatorNames = []string{"101", "102", "103"}

var PublicKeyPathSpot = 0
var AddressSpot = 1
var UserNameLen = 3

type userInfoMap map[string][2]string
type connectionNameToClient map[string]*Connection
type connectionNameToPubkey map[string]*rsa.PublicKey
type Connection struct {
	socket net.Conn
	data   chan []byte
}
type ConnectionsManager struct {
	connectedServersConnections connectionNameToClient
	connectedServersPubkey connectionNameToPubkey
	connections map[*Connection]bool
	register    chan *Connection
	unregister  chan *Connection
	publicKey   *rsa.PublicKey
	privateKey  *rsa.PrivateKey
	usersMap userInfoMap
	myName string
}

// TODO check how to create this map properly
// return a map containing {serverNum: {serverPublicKeyPath, serverAddress}}
func getUsersMap() userInfoMap {
	usersMap := userInfoMap{
		"001": {fmt.Sprintf(ServerPublicKeyPathFormat, "001"), fmt.Sprintf(ServerAddressFormat, "001")},
		"002": {fmt.Sprintf(ServerPublicKeyPathFormat, "002"), fmt.Sprintf(ServerAddressFormat, "002")},
		"003": {fmt.Sprintf(ServerPublicKeyPathFormat, "003"), fmt.Sprintf(ServerAddressFormat, "003")},
		"101": {fmt.Sprintf(MediatorPublicKeyPathFormat, "101"), fmt.Sprintf(MediatorAddressFormat, "101")},
		"102": {fmt.Sprintf(MediatorPublicKeyPathFormat, "102"), fmt.Sprintf(MediatorAddressFormat, "102")},
		"103": {fmt.Sprintf(MediatorPublicKeyPathFormat, "103"), fmt.Sprintf(MediatorAddressFormat, "103")},
		"201": {fmt.Sprintf(ClientPublicKeyPathFormat, "201"), fmt.Sprintf(ClientAddressFormat, "201")},
		"202": {fmt.Sprintf(ClientPublicKeyPathFormat, "202"), fmt.Sprintf(ClientAddressFormat, "202")},
		"203": {fmt.Sprintf(ClientPublicKeyPathFormat, "203"), fmt.Sprintf(ClientAddressFormat, "203")},
		}
	return usersMap
}

func createGeneralManager(usersMap userInfoMap, myName string) ConnectionsManager {
	myPublicKeyPath := usersMap[myName][PublicKeyPathSpot]
	privkey, pubkey := createKeys(myPublicKeyPath)
	manager := ConnectionsManager{
		connections: make(map[*Connection]bool),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
		publicKey:   pubkey,
		privateKey:  privkey,
		connectedServersConnections : make(connectionNameToClient),
		connectedServersPubkey : make(connectionNameToPubkey),
		usersMap : usersMap,
		myName: myName,
	}
	return manager
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func writeToFile(filePath string, data []byte) {
	// write the whole body at once
	os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	err := ioutil.WriteFile(filePath, data, 0644) // 0644 is the permission to create if file does not exists
	checkErr(err)
}

func readFromFile(filePath string) []byte {
	// read the whole file at once
	data, err := ioutil.ReadFile(filePath)
	checkErr(err)
	return data
}

func findValue(slic []string, val string) int {
		for i, n := range slic {
			if val == n {
				return i
			}
		}
		return len(slic)
}

func deleteValue(slic []string, val string) []string {
	if i := findValue(slic, val); i == len(slic) {
		return slic
	} else {
		return append(slic[:i], slic[i+1:]...)
	}
}
