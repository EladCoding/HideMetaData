package mixnet

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/gob"
	"encoding/json"
	laplaceRand "golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)


type ReplyMessage struct {
	From string
	To string
	Data []byte
}


// cipher vars
var RsaKeyBits = 2048
var CipherRsaLen = RsaKeyBits / 8
var AesKeyBytes = 32
var MsgBytes = 256  // TODO check size
var maxUserMsgSize = MsgBytes - 1

// general vars
const PathLen = 3
const UserNameLen = 3
const AddressSpot = 0
const PublicKeyPathSpot = 1
const PrivateKeyPathSpot = 2
const roundSlotTime = 3*time.Second
const fakeMsgsLaplaceMean = 100.0 // TODO change and check
const fakeMsgsLaplaceScale = (fakeMsgsLaplaceMean / 16) // TODO change and check
const minRoundSlotTime = roundSlotTime / 10


func ReadUserAddressMap() UserAddressMapType { // TODO change
	var usersMap UserAddressMapType
	jsonFile, err := os.Open(UserAddressesMapPath)
	CheckErrAndPanic(err)
	defer jsonFile.Close()
	jsonByteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(jsonByteValue, &usersMap)
	return usersMap
}


func ReadUserPubKeyMap() UserPublicKeyMapType { // TODO change
	var usersMap UserPublicKeyMapType
	jsonFile, err := os.Open(UserPublicKeysMapPath)
	CheckErrAndPanic(err)
	defer jsonFile.Close()
	jsonByteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(jsonByteValue, &usersMap)
	return usersMap
}


func ReadUserPrivKeyMap() UserPrivateKeyMapType { // TODO change
	var usersMap UserPrivateKeyMapType
	jsonFile, err := os.Open(UserPrivateKeysMapPath)
	CheckErrAndPanic(err)
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


func reverseShufflingReplyMsgs(msgs []EncryptedMsg, perm []int) []EncryptedMsg { // TODO check how to shuffle properly (cryptographlly)
	reversedMsgs := make([]EncryptedMsg, len(msgs))
	for i, v := range perm {
		reversedMsgs[i] = msgs[v]
	}
	return reversedMsgs
}


func ConvertMsgToBytes(onionMsg interface{}) EncryptedMsg {
	var inpBuf bytes.Buffer
	enc := gob.NewEncoder(&inpBuf)
	err := enc.Encode(onionMsg)
	CheckErrToLog(err)
	return inpBuf.Bytes()
}


func ConvertBytesToOnionMsg(onionBytes EncryptedMsg) OnionMessage {
	var outpBuf bytes.Buffer
	var onionMsg OnionMessage
	outpBuf.Write(onionBytes)
	dec := gob.NewDecoder(&outpBuf) // Will read from network.
	err := dec.Decode(&onionMsg)
	CheckErrToLog(err)
	return onionMsg
}


func ConvertBytesToReplyMsg(replyBytes EncryptedMsg) ReplyMessage {
	var outpBuf bytes.Buffer
	var replyMsg ReplyMessage
	outpBuf.Write(replyBytes)
	dec := gob.NewDecoder(&outpBuf) // Will read from network.
	err := dec.Decode(&replyMsg)
	CheckErrToLog(err)
	return replyMsg
}


func createOnionMessage(name string, serverName string, msgData []byte, mediatorsArr []string) (OnionMessage, []SecretKey) {
	var curPubKey ecdsa.PublicKey
	var curSymKey SecretKey
	var symKeys []SecretKey
	var onionMsg OnionMessage

	curOnionData := msgData
	hopesArr := append(mediatorsArr, serverName)
	for index, _ := range hopesArr {
		if index > 0 {
			curOnionData = ConvertMsgToBytes(onionMsg)
		}
		curHop := hopesArr[len(hopesArr)-index-1]
		curOnionData, curPubKey, curSymKey = hybridEncription(curOnionData, curHop)
		onionMsg = OnionMessage{
			name, // TODO check what about from
			serverName,
			curPubKey,
			curOnionData,
		}
		symKeys = append(symKeys, curSymKey)
	}
	return onionMsg, symKeys
}


func appendFakeMsgs(curMsgs []OnionMessage, numOfMsgsToAppend int, name string, mediatorsLeft []string) []OnionMessage {
	for i := 0; i < numOfMsgsToAppend; i += 1 {
		fakeMsgData, err := pkcs7padding([]byte("Fake"),MsgBytes)
		CheckErrToLog(err)
		//randServerName := scripts.ServerNames[rand.Intn(len(scripts.ServerNames))] // TODO back this when release (commented for debugging)
		randServerName := "001"
		cipherMsg, _ := createOnionMessage(name, randServerName, fakeMsgData, mediatorsLeft)
		curMsgs = append(curMsgs, cipherMsg)
	}
	return curMsgs
}


func DecryptOnionLayer(onionMsg OnionMessage, privKey *ecdsa.PrivateKey) (OnionMessage, []byte) {
	pubKey := onionMsg.PubKeyForSecret

	symKey := DecryptKeyForKeyExchange(pubKey, privKey)
	decryptedData, err := symmetricDecryption(onionMsg.Data, symKey)
	CheckErrToLog(err)

	onionMsg = ConvertBytesToOnionMsg(decryptedData)
	return onionMsg, symKey
}


func sampleFromLaplace() int {
	numOfFakes := int(myLaplace.Rand())
	return intMax(numOfFakes, 0)
}


var UserAddressesMap UserAddressMapType
var UserPubKeyMap UserPublicKeyMapType
var UserPrivKeyMap UserPrivateKeyMapType
var	myLaplace = distuv.Laplace{
	fakeMsgsLaplaceMean,
	fakeMsgsLaplaceScale,
	laplaceRand.NewSource(uint64(time.Now().UTC().UnixNano())),
}