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
	startReading         *sync.Mutex
	wg *sync.WaitGroup
}


func (l *CoordinatorListener) GetMessageFromClient(msg OnionMessage, reply *EncryptedMsg) error {
	l.startReading.Lock()
	l.wg.Add(1)
	l.GeneralListener.roundFinished.L.Lock()
	_, msgIndex := l.GeneralListener.readMessage(msg)
	l.startReading.Unlock()
	l.GeneralListener.roundFinished.Wait()
	*reply = l.GeneralListener.roundRepliedMsgs[msgIndex]
	l.GeneralListener.roundFinished.L.Unlock()
	l.wg.Done()
	return nil
}


func (l *CoordinatorListener) listenToCoordinatorAddress() {
	address := UserAddressesMap[l.GeneralListener.name]
	fmt.Printf("name: %v. listen to address: %v\n", l.GeneralListener.name, address)
	addy, err := net.ResolveTCPAddr("tcp", address)
	CheckErrToLog(err)
	inbound, err := net.ListenTCP("tcp", addy)
	CheckErrToLog(err)
	rpc.Register(l)
	go rpc.Accept(inbound)
}


func (l *CoordinatorListener) coordinateRounds() {
	nextRound := time.Now().Add(roundSlotTime)
	round := 1
	totalMsgs := 0
	for { // for each round
		if timeUntilNextRound := time.Until(nextRound); timeUntilNextRound > 0 {
			time.Sleep(timeUntilNextRound)
		}
		fmt.Printf("Coordinator: round: %v\n", round)
		round += 1
		nextRound = time.Now().Add(roundSlotTime) // TODO maybe change to slots instead of +1 from now
		l.startReading.Lock()
		l.GeneralListener.roundFinished.L.Lock()
		l.GeneralListener.msgMutex.Lock()
		curRoundMsgs, curRoundSymKeys := l.GeneralListener.readRoundMsgs()
		totalMsgs += len(curRoundMsgs)
		fmt.Printf("Coordinator: Got: %v msgs\n", totalMsgs)
		l.GeneralListener.msgMutex.Unlock()
		l.GeneralListener.roundRepliedMsgs = l.GeneralListener.sendRoundMessagesToNextHop(l.GeneralListener.nextHop, curRoundMsgs, curRoundSymKeys)
		l.GeneralListener.roundFinished.Broadcast()
		l.GeneralListener.roundFinished.L.Unlock()
		l.wg.Wait()
		l.startReading.Unlock()
	}
}


func StartCoordinator(name string, num int, nextHopName string) {
	fmt.Printf("Starting Coordinator %v...\n", name)
	var nextHop *rpc.Client
	var err error
	roundCond := sync.NewCond(&sync.Mutex{})
	nextHopAddress := UserAddressesMap[nextHopName]
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
		time.Now(),
	},
		&sync.Mutex{},
		&sync.WaitGroup{},
	}
	listener.listenToCoordinatorAddress()
	listener.coordinateRounds()
}
