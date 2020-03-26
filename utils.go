package main

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
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
var PathLen = 2

// servers map
var ServerPublicKeyPathFormat = "C:/repos/labProject/HideMetaData/keys/server%s/public_key.txt"
var ServerPrivateKeyPathFormat = "C:/repos/labProject/HideMetaData/keys/server%s/private_key.txt"
var ServerHost = "localhost"
var ServerPortFormat = "9%s"
var ServerAddressFormat = ServerHost + ":" + ServerPortFormat
// mediators map
var MediatorPublicKeyPathFormat = "C:/repos/labProject/HideMetaData/keys/mediator%s/public_key.txt"
var MediatorPrivateKeyPathFormat = "C:/repos/labProject/HideMetaData/keys/mediator%s/private_key.txt"
var MediatorHost = "localhost"
var MediatorHostFormat = "8%s"
var MediatorAddressFormat = MediatorHost + ":" + MediatorHostFormat
// clients map
var ClientPublicKeyPathFormat = "C:/repos/labProject/HideMetaData/keys/client%s/public_key.txt"
var ClientPrivateKeyPathFormat = "C:/repos/labProject/HideMetaData/keys/client%s/private_key.txt"
var ClientHost = "localhost"
var ClientHostFormat = "7%s"
var ClientAddressFormat = ClientHost + ":" + ClientHostFormat


type userInfoMap map[string][2]string
type connectionNameToClient map[string]*Client
type connectionNameToPubkey map[string]*rsa.PublicKey

var PublicKeyPathSpot = 0
var AddressSpot = 1
var UserNameLen = 3

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
		connections: make(map[*Client]bool),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		publicKey:   pubkey,
		privateKey:  privkey,
		connectedServersConnections : make(connectionNameToClient),
		connectedServersPubkey : make(connectionNameToPubkey),
		usersMap : usersMap,
	}
	return manager
}

func getAddress(serverMap userInfoMap, mediatorMap userInfoMap, mediator bool, addressNum string) string {
	if mediator {
		return mediatorMap[addressNum][AddressSpot]
	} else {
		return serverMap[addressNum][AddressSpot]
	}
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
