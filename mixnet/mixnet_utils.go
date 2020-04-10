package mixnet

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"github.com/EladCoding/HideMetaData/scripts"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

// cipher vars
var RsaKeyBits = 2048
var CipherRsaLen = RsaKeyBits / 8
var AesKeyBytes = 32

// general vars
var PathLen = 3
var UserNameLen = 3
var AddressSpot = 0
var PublicKeyPathSpot = 1
var PrivateKeyPathSpot = 2
var roundSlotTime = 15*time.Second


func readUserAddressMap() scripts.UserAddressMap { // TODO change
	var usersMap scripts.UserAddressMap
	jsonFile, err := os.Open(scripts.UserAddressesMapPath)
	scripts.CheckErrAndPanic(err)
	defer jsonFile.Close()
	jsonByteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(jsonByteValue, &usersMap)
	return usersMap
}


func readUserPubKeyMap() scripts.UserPublicKeyMap { // TODO change
	var usersMap scripts.UserPublicKeyMap
	jsonFile, err := os.Open(scripts.UserPublicKeysMapPath)
	scripts.CheckErrAndPanic(err)
	defer jsonFile.Close()
	jsonByteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(jsonByteValue, &usersMap)
	return usersMap
}


func readUserPrivKeyMap() scripts.UserPrivateKeyMap { // TODO change
	var usersMap scripts.UserPrivateKeyMap
	jsonFile, err := os.Open(scripts.UserPrivateKeysMapPath)
	scripts.CheckErrAndPanic(err)
	defer jsonFile.Close()
	jsonByteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(jsonByteValue, &usersMap)
	return usersMap
}


func shuffleMsgs(msgs []OnionMessage) ([]OnionMessage, []int) { // TODO check how to shuffle properly (cryptographlly)
	shuffledMsgs := make([]OnionMessage, len(msgs))
	perm := rand.Perm(len(msgs))
	for i, v := range perm {
		shuffledMsgs[v] = msgs[i]
	}
	return shuffledMsgs, perm
}


func reverseShufflingMsgs(msgs []OnionMessage, perm []int) []OnionMessage { // TODO check how to shuffle properly (cryptographlly)
	reversedMsgs := make([]OnionMessage, len(msgs))
	for i, v := range perm {
		reversedMsgs[i] = msgs[v]
	}
	return reversedMsgs
}


func ConvertOnionMsgToBytes(onionMsg OnionMessage) []byte {
	var inpBuf bytes.Buffer
	enc := gob.NewEncoder(&inpBuf)
	err := enc.Encode(onionMsg)
	scripts.CheckErrToLog(err)
	return inpBuf.Bytes()
}


func ConvertBytesToOnionMsg(onionBytes []byte) OnionMessage {
	var outpBuf bytes.Buffer
	var onionMsg OnionMessage
	outpBuf.Write(onionBytes)
	dec := gob.NewDecoder(&outpBuf) // Will read from network.
	err := dec.Decode(&onionMsg)
	scripts.CheckErrToLog(err)
	return onionMsg
}


var userAddressesMap = readUserAddressMap()
var userPubKeyMap = readUserPubKeyMap()
var userPrivKeyMap = readUserPrivKeyMap()


//var userAddressesMap scripts.UserAddressMap
//var userPubKeyMap scripts.UserPublicKeyMap
//var userPrivKeyMap scripts.UserPrivateKeyMap

