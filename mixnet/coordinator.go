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


type CoordinatorListener struct {
	name      string
	num int
	msgMutex  *sync.Mutex
	roundMsgs []OnionMessage
	roundSymKeys []scripts.SecretKey
	roundFinished *sync.Cond
	roundRepliedMsgs []scripts.EncryptedMsg
	nextHop *rpc.Client
}


func (l *CoordinatorListener) GetMessage(msg OnionMessage, reply *scripts.EncryptedMsg) error {
	_, msgIndex := l.readMessage(msg)
	l.roundFinished.Wait()
	*reply = l.roundRepliedMsgs[msgIndex]
	time.Sleep(100*time.Millisecond)
	return nil
}


func (l *CoordinatorListener) readMessage(msg OnionMessage) (OnionMessage, int) {
	from := msg.From
	encMsg := msg
	msg, symKey := DecryptOnionLayer(msg, scripts.DecodePrivateKey(userPrivKeyMap[l.name]))
	to := msg.To
	log.Printf("Coordinator %v Received OnionMessage:\nFrom: %v, To: %v, len: %v\n", l.name, from, to, len(encMsg.Data))
	msgIndex := l.appendMsgToRound(msg, symKey)
	return msg, msgIndex
}


func (l *CoordinatorListener) appendMsgToRound(msg OnionMessage, msgSymKey []byte) int {
	l.msgMutex.Lock()
	msgIndex := len(l.roundMsgs)
	l.roundMsgs = append(l.roundMsgs, msg)
	l.roundSymKeys = append(l.roundSymKeys, msgSymKey)
	l.msgMutex.Unlock()
	return msgIndex
}


func (l *CoordinatorListener) readRoundMsgs() ([]OnionMessage, []scripts.SecretKey) {
	curRoundMsgs := l.roundMsgs
	l.roundMsgs = make([]OnionMessage, 0)
	curRoundSymKeys := l.roundSymKeys
	l.roundSymKeys = make([]scripts.SecretKey, 0)
	return curRoundMsgs, curRoundSymKeys
}


func (l *CoordinatorListener) listenToMyAddress() {
	address := userAddressesMap[l.name]
	fmt.Printf("name: %v. listen to address: %v\n", l.name, address)
	addy, err := net.ResolveTCPAddr("tcp", address)
	scripts.CheckErrToLog(err)
	inbound, err := net.ListenTCP("tcp", addy)
	scripts.CheckErrToLog(err)
	rpc.Register(l)
	go rpc.Accept(inbound)
}


func (l *CoordinatorListener) sendRoundMessagesToNextHop(nextHop *rpc.Client, curRoundMsgs []OnionMessage, curRoundSymKeys []scripts.SecretKey) []scripts.EncryptedMsg {
	var curRoundRepliedMsgs []scripts.EncryptedMsg
	roundMsgsLen := len(curRoundMsgs)
	if len(curRoundMsgs) == 0 { // TODO remove from here
		return nil
	}

	// TODO random number of fakes
	curRoundMsgs = appendFakeMsgs(curRoundMsgs, fakeMsgsToAppend, l.name, scripts.MediatorNames[l.num:])

	curRoundShuffledMsgs, curRoundPerm := shuffleMsgs(curRoundMsgs)

	err := nextHop.Call("MediatorListener.GetBunchOfMessages", curRoundShuffledMsgs, &curRoundRepliedMsgs)
	scripts.CheckErrToLog(err)

	unShuffledCurRoundRepliedMsgs := reverseShufflingReplyMsgs(curRoundRepliedMsgs, curRoundPerm)[:roundMsgsLen]

	for index, msg := range unShuffledCurRoundRepliedMsgs {
		unShuffledCurRoundRepliedMsgs[index], err = symmetricEncryption(msg, curRoundSymKeys[index])
		scripts.CheckErrAndPanic(err)
	}
	return unShuffledCurRoundRepliedMsgs
}


func (l *CoordinatorListener) coordinateRounds() {
	nextRound := time.Now().Add(roundSlotTime) // TODO check what define round
	roundNumber := 1
	for { // for each round
		roundNumber += 1
		if timeUntilNextRound := time.Until(nextRound); timeUntilNextRound > 0 {
			time.Sleep(timeUntilNextRound)
			//continue // TODO check if needed
		}
		nextRound = time.Now().Add(roundSlotTime)
		l.msgMutex.Lock()
		curRoundMsgs, curRoundSymKeys := l.readRoundMsgs()
		l.roundRepliedMsgs = l.sendRoundMessagesToNextHop(l.nextHop, curRoundMsgs, curRoundSymKeys)
		l.roundFinished.Broadcast()
		l.msgMutex.Unlock()
	}
}


func StartCoordinator(name string, num int, nextHopName string) {
	fmt.Printf("Starting Coordinator %v...\n", name)
	var nextHop *rpc.Client
	var err error
	mycond := sync.NewCond(&sync.Mutex{})
	mycond.L.Lock()
	nextHopAddress := userAddressesMap[nextHopName]
	nextHop, err = rpc.Dial("tcp", nextHopAddress)
	scripts.CheckErrToLog(err)

	// listen to address
	listener := CoordinatorListener{
		name,
		num,
		&sync.Mutex{},
		make([]OnionMessage, 0),
		make([]scripts.SecretKey, 0),
		mycond,
		make([]scripts.EncryptedMsg, 0),
		nextHop,
	}
	listener.listenToMyAddress()
	listener.coordinateRounds()
}
