package main


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
