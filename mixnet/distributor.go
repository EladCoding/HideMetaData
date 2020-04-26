package mixnet

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
	"time"
)


type DistributorListener struct {
	GeneralListener GeneralListener
}


func (l *DistributorListener) GetRoundMsgs(msgs []OnionMessage, replies *[]EncryptedMsg) error {
	for _, msg := range msgs {
		l.GeneralListener.readMessage(msg)
	}
	decMsgs, symKeys := l.GeneralListener.readRoundMsgs()
	*replies = l.GeneralListener.sendRoundMessagesToNextHop(l.GeneralListener.nextHop, decMsgs, symKeys)

	time.Sleep(100*time.Millisecond)
	return nil
}


func (l *DistributorListener) listenToDistributorAddress() {
	address := userAddressesMap[l.GeneralListener.name]
	fmt.Printf("name: %v. listen to address: %v\n", l.GeneralListener.name, address)
	addy, err := net.ResolveTCPAddr("tcp", address)
	CheckErrToLog(err)
	inbound, err := net.ListenTCP("tcp", addy)
	CheckErrToLog(err)
	rpc.Register(l)
	go rpc.Accept(inbound)
}


func StartDistributor(name string, num int) {
	fmt.Printf("Starting Distributor %v...\n", name)
	var client *rpc.Client
	var clients map[string]*rpc.Client
	var err error
	// dial to all servers
	clients = make(map[string]*rpc.Client, 0)
	for _, serverName := range ServerNames {
		serverAddress := userAddressesMap[serverName]
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
	},
	}
	listener.listenToDistributorAddress()
	for {
		continue
	}
}
