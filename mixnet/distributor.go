package mixnet

import (
	"fmt"
	"github.com/EladCoding/HideMetaData/scripts"
	"net"
	"net/rpc"
	"sync"
	"time"
)


type DistributorListener struct {
	name      string
	num int
	msgMutex  *sync.Mutex
	roundMsgs []OnionMessage
	roundSymKeys [][]byte
	roundFinished *sync.Cond
	roundRepliedMsgs []scripts.EncryptedMsg
	clients scripts.ClientsMap
}


func (l *DistributorListener) readMessage(msg OnionMessage) (OnionMessage, int, []byte) {
	encData := msg.Data
	from := msg.From
	msg, symKey := DecryptOnionLayer(msg, scripts.DecodePrivateKey(userPrivKeyMap[l.name]))
	to := msg.To
	fmt.Printf("Mediator %v Received OnionMessage:\nFrom: %v, To: %v, Data: %x\n", l.name, from, to, encData)
	msgIndex := l.appendMsgToRound(msg, symKey)
	return msg, msgIndex, symKey
}


func (l *DistributorListener) GetBunchOfMessages(msgs []OnionMessage, replies *[]scripts.EncryptedMsg) error {
	for _, msg := range msgs {
		decMsg, _, symKey := l.readMessage(msg)
		var reply *scripts.EncryptedMsg
		destinitionServer := decMsg.To
		client := l.clients[destinitionServer]
		err := client.Call("ServerListener.GetMessage", decMsg, &reply)
		scripts.CheckErrToLog(err)
		*reply, err = symmetricEncryption(*reply, symKey)
		*replies = append(*replies, *reply)
	}
	time.Sleep(100*time.Millisecond)
	return nil
}



func (l *DistributorListener) appendMsgToRound(msg OnionMessage, msgSymKey []byte) int {
	l.msgMutex.Lock()
	msgIndex := len(l.roundMsgs)
	l.roundMsgs = append(l.roundMsgs, msg)
	l.roundSymKeys = append(l.roundSymKeys, msgSymKey)
	l.msgMutex.Unlock()
	return msgIndex
}


func (l *DistributorListener) readRoundMsgs() ([]OnionMessage, [][]byte) {
	curRoundMsgs := l.roundMsgs
	l.roundMsgs = make([]OnionMessage, 0)
	curRoundSymKeys := l.roundSymKeys
	l.roundSymKeys = make([][]byte, 0)
	return curRoundMsgs, curRoundSymKeys
}


func (l *DistributorListener) listenToMyAddress() {
	address := userAddressesMap[l.name]
	fmt.Printf("name: %v. address: %v\n", l.name, address)
	addy, err := net.ResolveTCPAddr("tcp", address)
	scripts.CheckErrToLog(err)
	inbound, err := net.ListenTCP("tcp", addy)
	scripts.CheckErrToLog(err)
	rpc.Register(l)
	go rpc.Accept(inbound)
}


func StartDistributor(name string, num int) {
	fmt.Printf("Starting Mediator %v...\n", name)
	var client *rpc.Client
	var clients map[string]*rpc.Client
	var err error
	// dial to all servers
	clients = make(map[string]*rpc.Client, 0)
	for _, serverName := range scripts.ServerNames {
		serverAddress := userAddressesMap[serverName]
		client, err = rpc.Dial("tcp", serverAddress)
		scripts.CheckErrToLog(err)
		clients[serverName] = client
	}

	// listen to address
	listener := DistributorListener{
		name,
		num,
		&sync.Mutex{},
		make([]OnionMessage, 0),
		make([][]byte, 0),
		&sync.Cond{},
		make([]scripts.EncryptedMsg, 0),
		clients,
	}

	listener.listenToMyAddress()
	for {
		continue
	}
}
