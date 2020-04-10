package mixnet

import (
	"fmt"
	"github.com/EladCoding/HideMetaData/scripts"
	"net"
	"net/rpc"
)


type ServerListener struct {
	name string
}
type Reply struct {
	From string
	Data []byte
}


func (l *ServerListener) GetMessage(msg OnionMessage, reply *Reply) error {
	encData := msg.Data
	from := msg.From
	*reply = Reply{l.name, encData}
	symKey := DecryptKeyForKeyExchange(msg.PubKeyForSecret[0], scripts.DecodePrivateKey(userPrivKeyMap[l.name]))
	decryptedData, err := symmetricDecryption(msg.Data, symKey)
	scripts.CheckErrToLog(err)
	fmt.Printf("Server %v Received Message:\nFrom: %v, Data: %v\n", l.name, from, string(decryptedData))
	return nil
}


func StartServer(name string) {
	fmt.Printf("Starting Server %v...\n", name)
	address := userAddressesMap[name]

	addy, err := net.ResolveTCPAddr("tcp", address)
	scripts.CheckErrToLog(err)
	inbound, err := net.ListenTCP("tcp", addy)
	scripts.CheckErrToLog(err)
	listener := ServerListener{name}
	rpc.Register(&listener)
	rpc.Accept(inbound)
}
