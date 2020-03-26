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

// TODO tmp for debug
var ServerAddress = "localhost:8000"
var MediatorAddress = "localhost:8001"
var ServerPublicKeyPath = "C:/repos/labProject/HideMetaData/keys/server/public_key"
var ServerPrivateKeyPath = "C:/repos/labProject/HideMetaData/keys/server/private_key"
var MediatorPublicKeyPath = "C:/repos/labProject/HideMetaData/keys/mediator/public_key"
var MediatorPrivateKeyPath = "C:/repos/labProject/HideMetaData/keys/mediator/private_key"

// TODO check how to create this map properly
// return a map containing {serverNum: {serverPublicKeyPath, serverAddress}}
func getServerMap() userInfoMap {
	serverMap := userInfoMap{
		"001": {fmt.Sprintf(ServerPublicKeyPathFormat, "001"), fmt.Sprintf(ServerAddressFormat, "001")},
		"002": {fmt.Sprintf(ServerPublicKeyPathFormat, "002"), fmt.Sprintf(ServerAddressFormat, "002")},
		"003": {fmt.Sprintf(ServerPublicKeyPathFormat, "003"), fmt.Sprintf(ServerAddressFormat, "003")},
		}
	return serverMap
}

// return a map containing {mediatorNum: {mediatorPublicKeyPath, mediatorAddress}}
func getMediatorMap() userInfoMap {
	mediatorMap := userInfoMap{
		"001": {fmt.Sprintf(MediatorPublicKeyPathFormat, "001"), fmt.Sprintf(MediatorAddressFormat, "001")},
		"002": {fmt.Sprintf(MediatorPublicKeyPathFormat, "002"), fmt.Sprintf(MediatorAddressFormat, "002")},
		"003": {fmt.Sprintf(MediatorPublicKeyPathFormat, "003"), fmt.Sprintf(MediatorAddressFormat, "003")},
	}
	return mediatorMap
}

// return a map containing {clientNum: {clientPublicKeyPath, clientAddress}}
func getClientMap() userInfoMap {
	clientMap := userInfoMap{
		"001": {fmt.Sprintf(ClientPublicKeyPathFormat, "001"), fmt.Sprintf(ClientAddressFormat, "001")},
		"002": {fmt.Sprintf(ClientPublicKeyPathFormat, "002"), fmt.Sprintf(ClientAddressFormat, "002")},
		"003": {fmt.Sprintf(ClientPublicKeyPathFormat, "003"), fmt.Sprintf(ClientAddressFormat, "003")},
	}
	return clientMap
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

func splitConnectionMap(connection connectionNameToClient) (string, *Client) {
	for k, v := range connection {
		return k, v
	}
	return "", nil
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
