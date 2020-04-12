package mixnet

import (
	"fmt"
	"github.com/EladCoding/HideMetaData/scripts"
	"net"
	"net/rpc"
	"time"
)


type ServerListener struct {
	name string
}
type ReplyMessage struct {
	From string
	To string
	Data []byte
}


func (l *ServerListener) GetMessage(msg OnionMessage, reply *scripts.EncryptedMsg) error {
	//encData := msg.Data
	from := msg.From
	symKey := DecryptKeyForKeyExchange(msg.PubKeyForSecret, scripts.DecodePrivateKey(userPrivKeyMap[l.name]))
	decryptedData, err := symmetricDecryption(msg.Data, symKey)
	scripts.CheckErrToLog(err)
	fmt.Printf("Server %v Received Message:\nFrom: %v, Data: %v\n", l.name, from, string(decryptedData))
	replyMsg := ReplyMessage{
		l.name, // TODO check what about from
		from,
		[]byte(fmt.Sprintf("I got ur msg: %s", string(decryptedData))),
	}
	*reply, err = symmetricEncryption(ConvertMsgToBytes(replyMsg), symKey)
	time.Sleep(100*time.Millisecond)
	return nil
}


func (l *ServerListener) listenToMyAddress() {
	address := userAddressesMap[l.name]
	addy, err := net.ResolveTCPAddr("tcp", address)
	scripts.CheckErrToLog(err)
	inbound, err := net.ListenTCP("tcp", addy)
	scripts.CheckErrToLog(err)
	rpc.Register(l)
	rpc.Accept(inbound)
}


func StartServer(name string) {
	fmt.Printf("Starting Server %v...\n", name)
	listener := ServerListener{name}
	listener.listenToMyAddress()
}
