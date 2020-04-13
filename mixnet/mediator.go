package mixnet

import (
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
	num int
	msgMutex  *sync.Mutex
	roundMsgs []OnionMessage
	roundSymKeys []scripts.SecretKey
	roundFinished *sync.Cond
	roundRepliedMsgs []scripts.EncryptedMsg
	nextHop *rpc.Client
}


func (l *MediatorListener) readMessage(msg OnionMessage) (OnionMessage, int, scripts.SecretKey) {
	from := msg.From
	msg, symKey := DecryptOnionLayer(msg, scripts.DecodePrivateKey(userPrivKeyMap[l.name]))
	to := msg.To
	log.Printf("Mediator %v Received OnionMessage:\nFrom: %v, To: %v\n", l.name, from, to)
	msgIndex := l.appendMsgToRound(msg, symKey)
	return msg, msgIndex, symKey
}


func (l *MediatorListener) sendRoundMessagesToNextHop(nextHop *rpc.Client, curRoundMsgs []OnionMessage, curRoundSymKeys []scripts.SecretKey) []scripts.EncryptedMsg {
	var curRoundRepliedMsgs []scripts.EncryptedMsg
	roundMsgsLen := len(curRoundMsgs)
	if len(curRoundMsgs) == 0 { // TODO remove from here
		return nil
	}

	// TODO random number of fakes
	curRoundMsgs = appendFakeMsgs(curRoundMsgs, fakeMsgsToAppend, l.name, scripts.MediatorNames[l.num:])

	curRoundShuffledMsgs, curRoundPerm := shuffleMsgs(curRoundMsgs)

	err := nextHop.Call("DistributorListener.GetBunchOfMessages", curRoundShuffledMsgs, &curRoundRepliedMsgs)
	scripts.CheckErrToLog(err)

	unShuffledCurRoundRepliedMsgs := reverseShufflingReplyMsgs(curRoundRepliedMsgs, curRoundPerm)[:roundMsgsLen]

	for index, msg := range unShuffledCurRoundRepliedMsgs {
		unShuffledCurRoundRepliedMsgs[index], err = symmetricEncryption(msg, curRoundSymKeys[index])
		scripts.CheckErrAndPanic(err)
	}
	return unShuffledCurRoundRepliedMsgs
}


func (l *MediatorListener) GetBunchOfMessages(msgs []OnionMessage, replies *[]scripts.EncryptedMsg) error {
	decMsgs := make([]OnionMessage, 0)
	symKeys := make([]scripts.SecretKey, 0)
	for _, msg := range msgs {
		decMsg, _, symKey := l.readMessage(msg)
		decMsgs = append(decMsgs, decMsg)
		symKeys = append(symKeys, symKey)
	}
	*replies = l.sendRoundMessagesToNextHop(l.nextHop, decMsgs, symKeys)

	time.Sleep(100*time.Millisecond)
	return nil
}



func (l *MediatorListener) appendMsgToRound(msg OnionMessage, msgSymKey []byte) int {
	l.msgMutex.Lock()
	msgIndex := len(l.roundMsgs)
	l.roundMsgs = append(l.roundMsgs, msg)
	l.roundSymKeys = append(l.roundSymKeys, msgSymKey)
	l.msgMutex.Unlock()
	return msgIndex
}


func (l *MediatorListener) readRoundMsgs() ([]OnionMessage, []scripts.SecretKey) {
	curRoundMsgs := l.roundMsgs
	l.roundMsgs = make([]OnionMessage, 0)
	curRoundSymKeys := l.roundSymKeys
	l.roundSymKeys = make([]scripts.SecretKey, 0)
	return curRoundMsgs, curRoundSymKeys
}


func (l *MediatorListener) listenToMyAddress() {
	address := userAddressesMap[l.name]
	fmt.Printf("name: %v. listen to address: %v\n", l.name, address)
	addy, err := net.ResolveTCPAddr("tcp", address)
	scripts.CheckErrToLog(err)
	inbound, err := net.ListenTCP("tcp", addy)
	scripts.CheckErrToLog(err)
	rpc.Register(l)
	go rpc.Accept(inbound)
}


func StartMediator(name string, num int, nextHopName string) {
	fmt.Printf("Starting Mediator %v...\n", name)
	var nextHop *rpc.Client
	var err error
	nextHopAddress := userAddressesMap[nextHopName]
	nextHop, err = rpc.Dial("tcp", nextHopAddress)
	scripts.CheckErrToLog(err)


	// listen to address
	listener := MediatorListener{
		name,
		num,
		&sync.Mutex{},
		make([]OnionMessage, 0),
		make([]scripts.SecretKey, 0),
		&sync.Cond{},
		make([]scripts.EncryptedMsg, 0),
		nextHop,
	}

	listener.listenToMyAddress()
	for {
		continue
	}
}
