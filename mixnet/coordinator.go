package mixnet

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
	"time"
)


type CoordinatorListener struct {
	GeneralListener GeneralListener
}


func (l *CoordinatorListener) GetMessageFromClient(msg OnionMessage, reply *EncryptedMsg) error {
	l.GeneralListener.roundFinished.L.Lock()
	_, msgIndex := l.GeneralListener.readMessage(msg)
	l.GeneralListener.roundFinished.Wait()
	*reply = l.GeneralListener.roundRepliedMsgs[msgIndex]
	l.GeneralListener.roundFinished.L.Unlock()
	time.Sleep(100*time.Millisecond)
	return nil
}


func (l *CoordinatorListener) listenToCoordinatorAddress() {
	address := userAddressesMap[l.GeneralListener.name]
	fmt.Printf("name: %v. listen to address: %v\n", l.GeneralListener.name, address)
	addy, err := net.ResolveTCPAddr("tcp", address)
	CheckErrToLog(err)
	inbound, err := net.ListenTCP("tcp", addy)
	CheckErrToLog(err)
	rpc.Register(l)
	go rpc.Accept(inbound)
}


func (l *CoordinatorListener) coordinateRounds() {
	nextRound := time.Now().Add(roundSlotTime) // TODO check what define round
	for { // for each round
		if timeUntilNextRound := time.Until(nextRound); timeUntilNextRound > 0 {
			time.Sleep(timeUntilNextRound)
			//continue // TODO check if needed
		}
		l.GeneralListener.roundFinished.L.Lock()
		nextRound = time.Now().Add(roundSlotTime)
		l.GeneralListener.msgMutex.Lock()
		curRoundMsgs, curRoundSymKeys := l.GeneralListener.readRoundMsgs()
		l.GeneralListener.roundRepliedMsgs = l.GeneralListener.sendRoundMessagesToNextHop(l.GeneralListener.nextHop, curRoundMsgs, curRoundSymKeys)
		l.GeneralListener.msgMutex.Unlock()
		l.GeneralListener.roundFinished.Broadcast()
		l.GeneralListener.roundFinished.L.Unlock()
	}
}


func StartCoordinator(name string, num int, nextHopName string) {
	fmt.Printf("Starting Coordinator %v...\n", name)
	var nextHop *rpc.Client
	var err error
	roundCond := sync.NewCond(&sync.Mutex{})
	nextHopAddress := userAddressesMap[nextHopName]
	nextHop, err = rpc.Dial("tcp", nextHopAddress)
	CheckErrToLog(err)

	// listen to address
	listener := CoordinatorListener{GeneralListener{
		name,
		num,
		&sync.Mutex{},
		make([]OnionMessage, 0),
		make([]SecretKey, 0),
		roundCond,
		make([]EncryptedMsg, 0),
		nextHop,
		true,
		false,
		nil,
	},
	}
	listener.listenToCoordinatorAddress()
	listener.coordinateRounds()
}
