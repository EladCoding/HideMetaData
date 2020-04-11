package mixnet

import (
	"bytes"
	"crypto/ecdsa"
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
var fakeMsgsToAppend = 3

// general vars
const PathLen = 3
const UserNameLen = 3
const AddressSpot = 0
const PublicKeyPathSpot = 1
const PrivateKeyPathSpot = 2
const roundSlotTime = 15*time.Second


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


func createOnionMessage(name string, serverName string, msgData []byte, mediatorsArr []string) OnionMessage {
	var curPubKey ecdsa.PublicKey
	var onionMsg OnionMessage

	curOnionData := msgData
	hopesArr := append(mediatorsArr, serverName)
	for index, _ := range hopesArr {
		if index > 0 {
			curOnionData = ConvertOnionMsgToBytes(onionMsg)
		}
		curHop := hopesArr[len(hopesArr)-index-1]
		curOnionData, curPubKey = hybridEncription(curOnionData, curHop)
		onionMsg = OnionMessage{
			name, // TODO check what about from
			serverName,
			curPubKey,
			curOnionData,
		}
	}
	return onionMsg
}


func appendFakeMsgs(curMsgs []OnionMessage, numOfMsgsToAppend int, name string, mediatorsLeft []string) []OnionMessage {
	for i := 0; i < numOfMsgsToAppend; i += 1 {
		fakeMsgData := make([]byte, 32) // TODO check size
		rand.Read(fakeMsgData)
		randServerName := scripts.ServerNames[rand.Intn(len(scripts.ServerNames))]
		cipherMsg := createOnionMessage(name, randServerName, fakeMsgData, mediatorsLeft)
		curMsgs = append(curMsgs, cipherMsg)
	}
	return curMsgs
}


var userAddressesMap = readUserAddressMap()
var userPubKeyMap = readUserPubKeyMap()
var userPrivKeyMap = readUserPrivKeyMap()


//var userAddressesMap scripts.UserAddressMap
//var userPubKeyMap scripts.UserPublicKeyMap
//var userPrivKeyMap scripts.UserPrivateKeyMap

