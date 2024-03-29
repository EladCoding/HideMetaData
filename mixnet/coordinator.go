package mixnet

import (
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)

// Coordinator listener main object.
type CoordinatorListener struct {
	GeneralListener GeneralListener
	startReading         *sync.Mutex
	wg *sync.WaitGroup
	coordinatorWaiting *sync.WaitGroup
	startRound bool
	startRoundMutex         *sync.RWMutex
}

// Rpc method that receiving a message from the user, and pass it to the next step.
func (l *CoordinatorListener) GetMessageFromClient(encMsg OnionMessage, reply *EncryptedMsg) error {
	decMsg, symKey := DecryptOnionLayer(encMsg, UserPrivKeyMap[l.GeneralListener.name])
	for {
		l.startReading.Lock()
		l.startRoundMutex.RLock()
		coordinatorWaiting := l.startRound
		l.startRoundMutex.RUnlock()
		if coordinatorWaiting {
			l.startReading.Unlock()
			continue
		}
		break
	}
	l.wg.Add(1)
	l.GeneralListener.roundFinished.L.Lock()
	msgIndex := l.GeneralListener.appendMsgToRound(decMsg, symKey)
	l.startReading.Unlock()
	l.GeneralListener.roundFinished.Wait()
	*reply = l.GeneralListener.roundRepliedMsgs[msgIndex]
	l.GeneralListener.roundFinished.L.Unlock()
	l.wg.Done()
	return nil
}

// Listen to a TCP local socket, as the coordinator.
func (l *CoordinatorListener) listenToCoordinatorAddress() {
	address := UserAddressesMap[l.GeneralListener.name]
	log.Printf("name: %v. listen to address: %v\n", l.GeneralListener.name, address)
	addy, err := net.ResolveTCPAddr("tcp", address)
	CheckErrToLog(err)
	inbound, err := net.ListenTCP("tcp", addy)
	CheckErrToLog(err)
	rpc.Register(l)
	go rpc.Accept(inbound)
}

// Manage the rounds as the first mediator at the mixnet architecture.
func (l *CoordinatorListener) coordinateRounds() {
	startTime := time.Now()
	nextRound := startTime
	round := 1
	totalMsgs := 0
	for { // for each round
		log.Printf("Coordinator: ready for round: %v\n", round)
		nextRound = nextRound.Add(RoundSlotTime)
		if timeUntilNextRound := time.Until(nextRound); timeUntilNextRound > 0 {
			time.Sleep(timeUntilNextRound)
		}
		l.startRoundMutex.Lock()
		l.startRound = true
		l.startRoundMutex.Unlock()
		log.Printf("Coordinator: starting round: %v\n", round)
		round += 1
		l.startReading.Lock()
		l.GeneralListener.roundFinished.L.Lock()
		l.GeneralListener.msgMutex.Lock()
		curRoundMsgs, curRoundSymKeys := l.GeneralListener.readRoundMsgs()
		totalMsgs += len(curRoundMsgs)
		log.Printf("Coordinator: Got: %v msgs\n", totalMsgs)
		l.GeneralListener.msgMutex.Unlock()
		l.GeneralListener.roundRepliedMsgs = l.GeneralListener.sendRoundMessagesToNextHop(l.GeneralListener.nextHop, curRoundMsgs, curRoundSymKeys)
		l.GeneralListener.roundFinished.Broadcast()
		l.GeneralListener.roundFinished.L.Unlock()
		l.wg.Wait()
		l.startRoundMutex.Lock()
		l.startRound = false
		l.startRoundMutex.Unlock()
		l.startReading.Unlock()
	}
}

// Start a mediator node as the Coordinator.
func StartCoordinator(name string, num int, nextHopName string) {
	log.Printf("Starting Coordinator %v...\n", name)
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
		&sync.WaitGroup{},
		false,
		&sync.RWMutex{},
	}
	listener.listenToCoordinatorAddress()
	listener.coordinateRounds()
}
