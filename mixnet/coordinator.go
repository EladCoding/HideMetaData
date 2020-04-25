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
	GeneralListener GeneralListener
}


func (l *CoordinatorListener) GetMessageFromClient(msg OnionMessage, reply *scripts.EncryptedMsg) error {
	_, msgIndex := l.GeneralListener.readMessage(msg)
	l.GeneralListener.roundFinished.Wait()
	*reply = l.GeneralListener.roundRepliedMsgs[msgIndex]
	time.Sleep(100*time.Millisecond)
	return nil
}


func (l *CoordinatorListener) listenToMyAddress() {
	address := userAddressesMap[l.GeneralListener.name]
	fmt.Printf("name: %v. listen to address: %v\n", l.GeneralListener.name, address)
	addy, err := net.ResolveTCPAddr("tcp", address)
	scripts.CheckErrToLog(err)
	inbound, err := net.ListenTCP("tcp", addy)
	scripts.CheckErrToLog(err)
	rpc.Register(l)
	go rpc.Accept(inbound)
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
		l.GeneralListener.msgMutex.Lock()
		curRoundMsgs, curRoundSymKeys := l.GeneralListener.readRoundMsgs()
		l.GeneralListener.roundRepliedMsgs = l.GeneralListener.sendRoundMessagesToNextHop(l.GeneralListener.nextHop, curRoundMsgs, curRoundSymKeys)
		l.GeneralListener.roundFinished.Broadcast()
		l.GeneralListener.msgMutex.Unlock()
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
	listener := CoordinatorListener{GeneralListener{
		name,
		num,
		&sync.Mutex{},
		make([]OnionMessage, 0),
		make([]scripts.SecretKey, 0),
		mycond,
		make([]scripts.EncryptedMsg, 0),
		nextHop,
		true,
		false,
		nil,
	},
	}
	listener.listenToMyAddress()
	listener.coordinateRounds()
}
