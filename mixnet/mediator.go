package mixnet

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)


type GeneralListener struct {
	name             string
	num              int
	msgMutex         *sync.Mutex
	roundMsgs        []OnionMessage
	roundSymKeys     []SecretKey
	roundFinished    *sync.Cond
	roundRepliedMsgs []EncryptedMsg
	nextHop          *rpc.Client
	isCoordinator    bool
	isDistributor    bool
	clients          ClientsMap
}


type MediatorListener struct {
	GeneralListener GeneralListener
}


func (l *GeneralListener) readMessage(msg OnionMessage) (OnionMessage, int) {
	from := msg.From
	encMsg := msg
	msg, symKey := DecryptOnionLayer(msg, DecodePrivateKey(userPrivKeyMap[l.name]))
	to := msg.To
	log.Printf("Mediator %v Received OnionMessage:\nFrom: %v, To: %v, len: %v\n", l.name, from, to, len(encMsg.Data))
	msgIndex := l.appendMsgToRound(msg, symKey)
	return msg, msgIndex
}


func (l *GeneralListener) sendRoundMessagesToNextHop(nextHop *rpc.Client, curRoundMsgs []OnionMessage, curRoundSymKeys []SecretKey) []EncryptedMsg {
	var curRoundRepliedMsgs []EncryptedMsg
	var err error
	roundMsgsLen := len(curRoundMsgs)
	if len(curRoundMsgs) == 0 { // TODO remove from here
		return nil
	}

	// TODO random number of fakes
	curRoundMsgs = appendFakeMsgs(curRoundMsgs, fakeMsgsToAppend, l.name, MediatorNames[l.num:])

	curRoundShuffledMsgs, curRoundPerm := shuffleMsgs(curRoundMsgs)

	if l.isDistributor {
		curRoundRepliedMsgs = make([]EncryptedMsg, 0)
		for _, msg := range curRoundShuffledMsgs {
			var reply *EncryptedMsg
			destinitionServer := msg.To
			client := l.clients[destinitionServer]
			err := client.Call("ServerListener.GetMessage", msg, &reply)
			CheckErrToLog(err)
			curRoundRepliedMsgs = append(curRoundRepliedMsgs, *reply)
		}
	} else if l.isCoordinator {
		err := nextHop.Call("MediatorListener.GetRoundMsgs", curRoundShuffledMsgs, &curRoundRepliedMsgs)
		CheckErrToLog(err)
	} else {
		err := nextHop.Call("DistributorListener.GetRoundMsgs", curRoundShuffledMsgs, &curRoundRepliedMsgs)
		CheckErrToLog(err)
	}

	unShuffledCurRoundRepliedMsgs := reverseShufflingReplyMsgs(curRoundRepliedMsgs, curRoundPerm)[:roundMsgsLen]

	for index, msg := range unShuffledCurRoundRepliedMsgs {
		unShuffledCurRoundRepliedMsgs[index], err = symmetricEncryption(msg, curRoundSymKeys[index])
		CheckErrAndPanic(err)
	}
	//fmt.Printf("unshuffle len: %v", len(unShuffledCurRoundRepliedMsgs))
	return unShuffledCurRoundRepliedMsgs
}


func (l *MediatorListener) GetRoundMsgs(msgs []OnionMessage, replies *[]EncryptedMsg) error {
	for _, msg := range msgs {
		l.GeneralListener.readMessage(msg)
	}
	decMsgs, symKeys := l.GeneralListener.readRoundMsgs()
	*replies = l.GeneralListener.sendRoundMessagesToNextHop(l.GeneralListener.nextHop, decMsgs, symKeys)

	time.Sleep(100*time.Millisecond)
	return nil
}


func (l *GeneralListener) appendMsgToRound(msg OnionMessage, msgSymKey []byte) int {
	l.msgMutex.Lock()
	msgIndex := len(l.roundMsgs)
	l.roundMsgs = append(l.roundMsgs, msg)
	l.roundSymKeys = append(l.roundSymKeys, msgSymKey)
	l.msgMutex.Unlock()
	return msgIndex
}


func (l *GeneralListener) readRoundMsgs() ([]OnionMessage, []SecretKey) {
	curRoundMsgs := l.roundMsgs
	l.roundMsgs = make([]OnionMessage, 0)
	curRoundSymKeys := l.roundSymKeys
	l.roundSymKeys = make([]SecretKey, 0)
	return curRoundMsgs, curRoundSymKeys
}


func (l *MediatorListener) listenToMediatorAddress() {
	address := userAddressesMap[l.GeneralListener.name]
	fmt.Printf("name: %v. listen to address: %v\n", l.GeneralListener.name, address)
	addy, err := net.ResolveTCPAddr("tcp", address)
	CheckErrToLog(err)
	inbound, err := net.ListenTCP("tcp", addy)
	CheckErrToLog(err)
	rpc.Register(l)
	go rpc.Accept(inbound)
}


func StartMediator(name string, num int, nextHopName string) {
	fmt.Printf("Starting Mediator %v...\n", name)
	var nextHop *rpc.Client
	var err error
	nextHopAddress := userAddressesMap[nextHopName]
	nextHop, err = rpc.Dial("tcp", nextHopAddress)
	CheckErrToLog(err)


	// listen to address
	listener := MediatorListener{GeneralListener{
		name,
		num,
		&sync.Mutex{},
		make([]OnionMessage, 0),
		make([]SecretKey, 0),
		&sync.Cond{},
		make([]EncryptedMsg, 0),
		nextHop,
		false,
		false,
		nil,
	},
	}

	listener.listenToMediatorAddress()
	for {
		continue
	}
}
