package main


import "io/ioutil"


var ConnectionSuccessfulAnswer = []byte("Success!\nYou are now connected")
var MessageReceivedAnswer = []byte("Success!\nYour message has received\n")
var ServerHost = "localhost"
var ServerPort = "8000"
var ServerAddress = ServerHost + ":" + ServerPort
var MediatorHost = "localhost"
var MediatorPort = "8001"
var MediatorAddress = MediatorHost + ":" + MediatorPort

var RsaKeyBits = 2048
var CipherRsaLen = RsaKeyBits / 8
var AesKeyBytes = 32
var ServerPublicKeyPath = "C:/repos/labProject/HideMetaData/keys/server/public_key"
var ServerPrivateKeyPath = "C:/repos/labProject/HideMetaData/keys/server/private_key"
var MediatorPublicKeyPath = "C:/repos/labProject/HideMetaData/keys/mediator/public_key"
var MediatorPrivateKeyPath = "C:/repos/labProject/HideMetaData/keys/mediator/private_key"

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func writeToFile(filePath string, data []byte) {
	// write the whole body at once
	err := ioutil.WriteFile(filePath, data, 0644) // 0644 is the permission to create if file does not exists
	checkErr(err)
}

func readFromFile(filePath string) []byte {
	// read the whole file at once
	data, err := ioutil.ReadFile(filePath)
	checkErr(err)
	return data
}
