package mixnet

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/EladCoding/HideMetaData/scripts"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)


type MediatorListener struct {
	name      string
	msgMutex  *sync.Mutex
	roundMsgs []OnionMessage
	lastMediator bool
}


func (l *MediatorListener) GetMessage(msg OnionMessage, reply *Reply) error {
	encData := msg.Data
	from := msg.From
	*reply = Reply{l.name, msg.Data}
	msg = DecryptOnionLayer(msg, scripts.DecodePrivateKey(userPrivKeyMap[l.name]))
	to := msg.To
	fmt.Printf("Mediator %v Received OnionMessage:\nFrom: %v, To: %v, Data: %x\n", l.name, from, to, encData)
	return l.appendMsgToRound(msg)
}


func DecryptOnionLayer(onionMsg OnionMessage, privKey *ecdsa.PrivateKey) OnionMessage {
	pubKeys := onionMsg.PubKeyForSecret
	curPubKey := pubKeys[len(pubKeys) - 1]

	symKey := DecryptKeyForKeyExchange(curPubKey, privKey)
	decryptedData, err := symmetricDecryption(onionMsg.Data, symKey)
	scripts.CheckErrToLog(err)

	onionMsg.Data = decryptedData
	onionMsg.PubKeyForSecret = pubKeys[:len(pubKeys) - 1]
	return onionMsg
}


func (l *MediatorListener) appendMsgToRound(msg OnionMessage) error {
	l.msgMutex.Lock()
	l.roundMsgs = append(l.roundMsgs, msg)
	l.msgMutex.Unlock()
	return nil
}


func (l *MediatorListener) readRoundMsgs() []OnionMessage {
	l.msgMutex.Lock()
	curRoundMsgs := l.roundMsgs
	l.roundMsgs = make([]OnionMessage, 0)
	l.msgMutex.Unlock()
	return curRoundMsgs
}


func StartMediator(name string, nextHopName string, firstMediator bool, lastMediator bool) {
	fmt.Printf("Starting Mediator %v...\n", name)
	var client *rpc.Client
	var clients map[string]*rpc.Client
	var err error
	if !lastMediator { // dial to the next hop
		nextHopAddress := userAddressesMap[nextHopName]
		client, err = rpc.Dial("tcp", nextHopAddress)
		scripts.CheckErrToLog(err)
	} else { // dial to all servers
		clients = make(map[string]*rpc.Client, 0)
		for _, serverName := range scripts.ServerNames {
			serverAddress := userAddressesMap[serverName]
			client, err = rpc.Dial("tcp", serverAddress)
			scripts.CheckErrToLog(err)
			clients[serverName] = client
		}
	}

	address := userAddressesMap[name]
	fmt.Printf("name: %v. address: %v\n", name, address)
	addy, err := net.ResolveTCPAddr("tcp", address)
	scripts.CheckErrToLog(err)
	inbound, err := net.ListenTCP("tcp", addy)
	scripts.CheckErrToLog(err)
	listener := MediatorListener{
		name,
		&sync.Mutex{},
		make([]OnionMessage, 0),
		lastMediator,
	}

	rpc.Register(&listener)
	go rpc.Accept(inbound)

	nextRound := time.Now().Add(roundSlotTime) // TODO check what define round
	roundNumber := 1
	for { // for each round
		fmt.Printf("Round number: %v\n", roundNumber)
		roundNumber += 1
		if timeUntilNextRound := time.Until(nextRound); timeUntilNextRound > 0 {
			time.Sleep(timeUntilNextRound)
			//continue // TODO check if needed
		}
		nextRound = time.Now().Add(roundSlotTime)
		curRoundShuffledMsgs, curRoundPerm := shuffleMsgs(listener.readRoundMsgs())
		reverseShufflingMsgs(curRoundShuffledMsgs, curRoundPerm) // TODO use curRoundPerm for reply
		for _, msg := range curRoundShuffledMsgs {
			var reply Reply
			if lastMediator {
				destinitionServer := msg.To
				client = clients[destinitionServer]
				err = client.Call("ServerListener.GetMessage", msg, &reply)
				scripts.CheckErrToLog(err)
			} else {
				err = client.Call("MediatorListener.GetMessage", msg, &reply)
				scripts.CheckErrToLog(err)
			}
			log.Printf("Reply: %v, From: %v, Data: %v", reply, reply.From, reply.Data)
		}
	}
}
