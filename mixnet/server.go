package mixnet

import (
	"fmt"
	"net"
	"net/rpc"
)


type ServerListener struct {
	name string
}
type Reply struct {
	From string
	Data string
}


func (l *ServerListener) GetMessage(msg Message, reply *Reply) error {
	from := msg.From
	data := msg.Data
	fmt.Printf("Server %v Received Message:\nFrom: %v, Data: %v\n", l.name, from, data)
	*reply = Reply{l.name, data}
	return nil
}


func StartServer(name string) {
	fmt.Printf("Starting Server %v...\n", name)
	address := usersMap[name][AddressSpot]

	addy, err := net.ResolveTCPAddr("tcp", address)
	checkErr(err)
	inbound, err := net.ListenTCP("tcp", addy)
	checkErr(err)
	listener := ServerListener{name}
	rpc.Register(&listener)
	rpc.Accept(inbound)
}
