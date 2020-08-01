package mixnet

import (
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)

// General listener object.
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
	lastRoundTime time.Time
}

// Mediator listener main object.
type MediatorListener struct {
	GeneralListener GeneralListener
}

// Rpc method that read a message from an other mediator, and pass it to the next one.
func (l *GeneralListener) readMessageFromMediator(encMsg OnionMessage, msgIndex int, wg *sync.WaitGroup) (OnionMessage, int) {
	msg, symKey := DecryptOnionLayer(encMsg, UserPrivKeyMap[l.name])
	l.roundMsgs[msgIndex] = msg
	l.roundSymKeys[msgIndex] = symKey
	wg.Done()
	return msg, msgIndex
}

// Send a message to a specific server.
func (l *GeneralListener) sendMsgToServer(msg OnionMessage, msgIndex int, curRoundRepliedMsgs []EncryptedMsg,
	replyFromServerMutex *sync.Mutex, wg *sync.WaitGroup) {
	var reply *EncryptedMsg
	destinitionServer := msg.To
	client := l.clients[destinitionServer]
	err := client.Call("ServerListener.GetMessage", msg, &reply)
	CheckErrToLog(err)
	replyFromServerMutex.Lock()
	curRoundRepliedMsgs[msgIndex] = *reply
	replyFromServerMutex.Unlock()
	wg.Done()
}

// Pass all the round messages it to the next mediator hop.
func (l *GeneralListener) sendRoundMessagesToNextHop(nextHop *rpc.Client, curRoundMsgs []OnionMessage, curRoundSymKeys []SecretKey) []EncryptedMsg {
	var curRoundRepliedMsgs []EncryptedMsg
	var err error
	roundMsgsLen := len(curRoundMsgs)

	fakeMsgsToAppend := sampleFromLaplace()
	curRoundMsgs = appendFakeMsgs(curRoundMsgs, fakeMsgsToAppend, l.name, MediatorNames[l.num:])

	curRoundShuffledMsgs, curRoundPerm := shuffleMsgs(curRoundMsgs)

	if l.isDistributor {
		curRoundRepliedMsgs = make([]EncryptedMsg, len(curRoundShuffledMsgs))
		replyFromServerMutex := &sync.Mutex{}
		wg := &sync.WaitGroup{}
		wg.Add(len(curRoundShuffledMsgs))
		for msgIndex, msg := range curRoundShuffledMsgs {
			go l.sendMsgToServer(msg, msgIndex, curRoundRepliedMsgs, replyFromServerMutex, wg)
		}
		wg.Wait()
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
	return unShuffledCurRoundRepliedMsgs
}

// Start receiving messages for this round.
func (l *MediatorListener) GetRoundMsgs(msgs []OnionMessage, replies *[]EncryptedMsg) error {
	l.GeneralListener.lastRoundTime = time.Now()
	wg := &sync.WaitGroup{}
	wg.Add(len(msgs))

	l.GeneralListener.roundMsgs = make([]OnionMessage, len(msgs))
	l.GeneralListener.roundSymKeys = make([]SecretKey, len(msgs))
	for msgIndex, msg := range msgs {
		go l.GeneralListener.readMessageFromMediator(msg, msgIndex, wg)
	}
	wg.Wait()
	decMsgs, symKeys := l.GeneralListener.readRoundMsgs()
	*replies = l.GeneralListener.sendRoundMessagesToNextHop(l.GeneralListener.nextHop, decMsgs, symKeys)

	return nil
}

// Append a message to the messages received this round.
func (l *GeneralListener) appendMsgToRound(msg OnionMessage, msgSymKey []byte) int {
	l.msgMutex.Lock()
	msgIndex := len(l.roundMsgs)
	l.roundMsgs = append(l.roundMsgs, msg)
	l.roundSymKeys = append(l.roundSymKeys, msgSymKey)
	l.msgMutex.Unlock()
	return msgIndex
}

// Read the whole message that received this round.
func (l *GeneralListener) readRoundMsgs() ([]OnionMessage, []SecretKey) {
	curRoundMsgs := l.roundMsgs
	l.roundMsgs = make([]OnionMessage, 0)
	curRoundSymKeys := l.roundSymKeys
	l.roundSymKeys = make([]SecretKey, 0)
	return curRoundMsgs, curRoundSymKeys
}

// Listen to a TCP local socket, as a mediator.
func (l *MediatorListener) listenToMediatorAddress() {
	address := UserAddressesMap[l.GeneralListener.name]
	log.Printf("name: %v. listen to address: %v\n", l.GeneralListener.name, address)
	addy, err := net.ResolveTCPAddr("tcp", address)
	CheckErrToLog(err)
	inbound, err := net.ListenTCP("tcp", addy)
	CheckErrToLog(err)
	rpc.Register(l)
	rpc.Accept(inbound)
}

// Start a mediator node as a Mediator.
func StartMediator(name string, num int, nextHopName string) {
	log.Printf("Starting Mediator %v...\n", name)
	var nextHop *rpc.Client
	var err error
	nextHopAddress := UserAddressesMap[nextHopName]
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
		time.Now(),
	},
	}
	listener.listenToMediatorAddress()
}
