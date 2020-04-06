package mixnet

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)


type MediatorListener struct {
	name string
	mediatorChannel chan Message
}


func (l *MediatorListener) GetMessage(msg Message, reply *Reply) error {
	from := msg.From
	data := msg.Data
	fmt.Printf("Mediator %v Received Message:\nFrom: %v, Data: %v\n", l.name, from, data)
	*reply = Reply{l.name, data}
	l.mediatorChannel <- msg
	return nil
}


func StartMediator(name string) {
	fmt.Printf("Starting Mediator %v...\n", name)

	nextHopAddress := usersMap["001"][AddressSpot]
	client, err := rpc.Dial("tcp", nextHopAddress)
	checkErr(err)
	mediatorChannel := make(chan Message)

	address := usersMap[name][AddressSpot]
	addy, err := net.ResolveTCPAddr("tcp", address)
	checkErr(err)
	inbound, err := net.ListenTCP("tcp", addy)
	checkErr(err)
	listener := MediatorListener{name, mediatorChannel}
	rpc.Register(&listener)
	go rpc.Accept(inbound)

	for {
		msg := <- mediatorChannel
		var reply Reply
		err = client.Call("ServerListener.GetMessage", msg, &reply)
		checkErr(err)
		log.Printf("Reply: %v, From: %v, Data: %v", reply, reply.From, reply.Data)
	}
}
