package mixnet

import (
	"fmt"
	"github.com/EladCoding/HideMetaData/scripts"
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
	roundSymKeys [][]byte
	roundFinished *sync.Cond
	roundRepliedMsgs []scripts.EncryptedMsg
}


func (l *CoordinatorListener) GetMessage(msg OnionMessage, reply *scripts.EncryptedMsg) error {
	_, msgIndex := l.readMessage(msg)
	l.roundFinished.Wait()
	*reply = l.roundRepliedMsgs[msgIndex]
	time.Sleep(100*time.Millisecond)
	return nil
}


func (l *CoordinatorListener) readMessage(msg OnionMessage) (OnionMessage, int) {
	encData := msg.Data
	from := msg.From
	msg, symKey := DecryptOnionLayer(msg, scripts.DecodePrivateKey(userPrivKeyMap[l.name]))
	to := msg.To
	fmt.Printf("Mediator %v Received OnionMessage:\nFrom: %v, To: %v, Data: %x\n", l.name, from, to, encData)
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


func (l *CoordinatorListener) readRoundMsgs() ([]OnionMessage, [][]byte) {
	curRoundMsgs := l.roundMsgs
	l.roundMsgs = make([]OnionMessage, 0)
	curRoundSymKeys := l.roundSymKeys
	l.roundSymKeys = make([][]byte, 0)
	return curRoundMsgs, curRoundSymKeys
}


func (l *CoordinatorListener) listenToMyAddress() {
	address := userAddressesMap[l.name]
	fmt.Printf("name: %v. address: %v\n", l.name, address)
	addy, err := net.ResolveTCPAddr("tcp", address)
	scripts.CheckErrToLog(err)
	inbound, err := net.ListenTCP("tcp", addy)
	scripts.CheckErrToLog(err)
	rpc.Register(l)
	go rpc.Accept(inbound)
}


func (l *CoordinatorListener) coordinateRounds(client *rpc.Client) {
	nextRound := time.Now().Add(roundSlotTime) // TODO check what define round
	roundNumber := 1
	for { // for each round
		roundNumber += 1
		if timeUntilNextRound := time.Until(nextRound); timeUntilNextRound > 0 {
			time.Sleep(timeUntilNextRound)
			//continue // TODO check if needed
		}
		l.msgMutex.Lock()
		nextRound = time.Now().Add(roundSlotTime)
		curRoundMsgs, curRoundSymKeys := l.readRoundMsgs() //TODO return this and use keys
		//curRoundMsgs = appendFakeMsgs(curRoundMsgs, fakeMsgsToAppend, l.name, scripts.MediatorNames[l.num:])
		// TODO random number of fakes
		if len(curRoundMsgs) == 0 { // TODO remove from here
			l.msgMutex.Unlock()
			continue
		}

		curRoundShuffledMsgs, curRoundPerm := shuffleMsgs(curRoundMsgs)

		var curRoundRepliedMsgs []scripts.EncryptedMsg
		err := client.Call("DistributorListener.GetBunchOfMessages", curRoundShuffledMsgs, &curRoundRepliedMsgs)
		scripts.CheckErrToLog(err)
		unShuffledCurRoundRepliedMsgs := reverseShufflingReplyMsgs(curRoundRepliedMsgs, curRoundPerm) // TODO use curRoundPerm for reply

		for index, msg := range unShuffledCurRoundRepliedMsgs {
			unShuffledCurRoundRepliedMsgs[index], err = symmetricEncryption(msg, curRoundSymKeys[index])
			scripts.CheckErrAndPanic(err)
		}
		l.roundRepliedMsgs = unShuffledCurRoundRepliedMsgs
		l.roundFinished.Broadcast()
		l.msgMutex.Unlock()
	}
}


func StartCoordinator(name string, num int, nextHopName string) {
	fmt.Printf("Starting Mediator %v...\n", name)
	var client *rpc.Client
	var err error
	mycond := sync.NewCond(&sync.Mutex{})
	mycond.L.Lock()
	nextHopAddress := userAddressesMap[nextHopName]
	client, err = rpc.Dial("tcp", nextHopAddress)
	scripts.CheckErrToLog(err)

	// listen to address
	listener := CoordinatorListener{
		name,
		num,
		&sync.Mutex{},
		make([]OnionMessage, 0),
		make([][]byte, 0),
		mycond,
		make([]scripts.EncryptedMsg, 0),
	}
	listener.listenToMyAddress()
	listener.coordinateRounds(client)
}
