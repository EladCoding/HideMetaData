package mixnet
//
//import (
//	"fmt"
//	"github.com/EladCoding/HideMetaData/scripts"
//	"log"
//	"net"
//	"net/rpc"
//	"sync"
//	"time"
//)
//
//
//type MediatorListener struct {
//	name      string
//	num int
//	msgMutex  *sync.Mutex
//	roundMsgs []OnionMessage
//	roundSymKeys [][]byte
//	lastMediator bool
//}
//
//
//func (l *MediatorListener) GetMessageFromClient(msg OnionMessage, reply *Reply) error {
//	encData := msg.Data
//	from := msg.From
//	msg, symKey := DecryptOnionLayer(msg, scripts.DecodePrivateKey(userPrivKeyMap[l.name]))
//	to := msg.To
//	fmt.Printf("Mediator %v Received OnionMessage:\nFrom: %v, To: %v, Data: %x\n", l.name, from, to, encData)
//	err := l.appendMsgToRound(msg, symKey)
//	scripts.CheckErrToLog(err)
//	*reply = Reply{l.name, encData}
//	return nil
//}
//
//
//func (l *MediatorListener) appendMsgToRound(msg OnionMessage, msgSymKey []byte) error {
//	l.msgMutex.Lock()
//	l.roundMsgs = append(l.roundMsgs, msg)
//	l.roundSymKeys = append(l.roundSymKeys, msgSymKey)
//	l.msgMutex.Unlock()
//	return nil
//}
//
//
//func (l *MediatorListener) readRoundMsgs() ([]OnionMessage, [][]byte) {
//	l.msgMutex.Lock()
//	curRoundMsgs := l.roundMsgs
//	l.roundMsgs = make([]OnionMessage, 0)
//	curRoundSymKeys := l.roundSymKeys
//	l.roundSymKeys = make([][]byte, 0)
//	l.msgMutex.Unlock()
//	return curRoundMsgs, curRoundSymKeys
//}
//
//
//func (l *MediatorListener) listenToMyAddress() {
//	address := userAddressesMap[l.name]
//	fmt.Printf("name: %v. address: %v\n", l.name, address)
//	addy, err := net.ResolveTCPAddr("tcp", address)
//	scripts.CheckErrToLog(err)
//	inbound, err := net.ListenTCP("tcp", addy)
//	scripts.CheckErrToLog(err)
//	rpc.Register(l)
//	go rpc.Accept(inbound)
//}
//
//
//func (l *MediatorListener) coordinateRounds(client *rpc.Client, clients map[string]*rpc.Client) {
//	nextRound := time.Now().Add(roundSlotTime) // TODO check what define round
//	roundNumber := 1
//	for { // for each round
//		fmt.Printf("Round number: %v\n", roundNumber)
//		roundNumber += 1
//		if timeUntilNextRound := time.Until(nextRound); timeUntilNextRound > 0 {
//			time.Sleep(timeUntilNextRound)
//			//continue // TODO check if needed
//		}
//		nextRound = time.Now().Add(roundSlotTime)
//		curRoundMsgs, _ := l.readRoundMsgs()
//		//curRoundMsgs, curRoundSymKeys := l.readRoundMsgs() //TODO return this and use keys
//		curRoundMsgs = appendFakeMsgs(curRoundMsgs, fakeMsgsToAppend, l.name, scripts.MediatorNames[l.num:])
//
//		curRoundShuffledMsgs, curRoundPerm := shuffleMsgs(curRoundMsgs)
//		reverseShufflingMsgs(curRoundShuffledMsgs, curRoundPerm) // TODO use curRoundPerm for reply
//		for _, msg := range curRoundShuffledMsgs {
//			var reply Reply
//			if l.lastMediator {
//				destinitionServer := msg.To
//				client = clients[destinitionServer]
//				err := client.Call("ServerListener.GetMessageFromClient", msg, &reply)
//				scripts.CheckErrToLog(err)
//			} else {
//				err := client.Call("DistributorListener.GetMessageFromClient", msg, &reply)
//				scripts.CheckErrToLog(err)
//			}
//			log.Printf("Reply: %v, From: %v, Data: %v", reply, reply.From, reply.Data)
//		}
//	}
//}
//
//
//func StartMediator(name string, num int, nextHopName string, firstMediator bool, lastMediator bool) {
//	fmt.Printf("Starting Mediator %v...\n", name)
//	var client *rpc.Client
//	var clients map[string]*rpc.Client
//	var err error
//	if !lastMediator { // dial to the next hop
//		nextHopAddress := userAddressesMap[nextHopName]
//		client, err = rpc.Dial("tcp", nextHopAddress)
//		scripts.CheckErrToLog(err)
//	} else { // dial to all servers
//		clients = make(map[string]*rpc.Client, 0)
//		for _, serverName := range scripts.ServerNames {
//			serverAddress := userAddressesMap[serverName]
//			client, err = rpc.Dial("tcp", serverAddress)
//			scripts.CheckErrToLog(err)
//			clients[serverName] = client
//		}
//	}
//
//	// listen to address
//	listener := MediatorListener
//		name,
//		num,
//		&sync.Mutex{},
//		make([]OnionMessage, 0),
//		make([][]byte, 0),
//		lastMediator,
//	}
//	listener.listenToMyAddress()
//	listener.coordinateRounds(client, clients)
//}
