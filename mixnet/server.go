package mixnet

import (
	"fmt"
	"github.com/EladCoding/HideMetaData/scripts"
	"log"
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
	from := msg.From
	log.Printf("Server %v Received OnionMessage:\nFrom: %v, len: %v\n", l.name, from, len(msg.Data))
	symKey := DecryptKeyForKeyExchange(msg.PubKeyForSecret, scripts.DecodePrivateKey(userPrivKeyMap[l.name]))
	decryptedData, err := symmetricDecryption(msg.Data, symKey)
	scripts.CheckErrToLog(err)
	decryptedData, err = pkcs7strip(decryptedData, MsgBytes)
	scripts.CheckErrToLog(err)
	log.Printf("Server %v Received Message:\nFrom: %v, Data: %v\n", l.name, from, string(decryptedData))
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
	fmt.Printf("name: %v. listen to address: %v\n", l.name, address)
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
