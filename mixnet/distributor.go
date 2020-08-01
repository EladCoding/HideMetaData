package mixnet

import (
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)


type DistributorListener struct {
	GeneralListener GeneralListener
}


func (l *DistributorListener) GetRoundMsgs(msgs []OnionMessage, replies *[]EncryptedMsg) error {
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


func (l *DistributorListener) listenToDistributorAddress() {
	address := UserAddressesMap[l.GeneralListener.name]
	log.Printf("name: %v. listen to address: %v\n", l.GeneralListener.name, address)
	addy, err := net.ResolveTCPAddr("tcp", address)
	CheckErrToLog(err)
	inbound, err := net.ListenTCP("tcp", addy)
	CheckErrToLog(err)
	rpc.Register(l)
	rpc.Accept(inbound)
}


func StartDistributor(name string, num int) {
	log.Printf("Starting Distributor %v...\n", name)
	var client *rpc.Client
	var clients map[string]*rpc.Client
	var err error
	// dial to all servers
	clients = make(map[string]*rpc.Client, 0)
	for _, serverName := range ServerNames {
		serverAddress := UserAddressesMap[serverName]
		client, err = rpc.Dial("tcp", serverAddress)
		CheckErrToLog(err)
		clients[serverName] = client
	}

	// listen to address
	listener := DistributorListener{GeneralListener{
		name,
		num,
		&sync.Mutex{},
		make([]OnionMessage, 0),
		make([]SecretKey, 0),
		&sync.Cond{},
		make([]EncryptedMsg, 0),
		nil,
		false,
		true,
		clients,
		time.Now(),
	},
	}
	listener.listenToDistributorAddress()
}
