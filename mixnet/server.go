package mixnet

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"
)


type ServerListener struct {
	name string
}


func (l *ServerListener) GetMessage(msg OnionMessage, reply *EncryptedMsg) error {
	from := msg.From
	log.Printf("Server %v Received OnionMessage:\nFrom: %v, len: %v\n", l.name, from, len(msg.Data))
	symKey := DecryptKeyForKeyExchange(msg.PubKeyForSecret, DecodePrivateKey(userPrivKeyMap[l.name]))
	decryptedData, err := symmetricDecryption(msg.Data, symKey)
	CheckErrToLog(err)
	decryptedData, err = pkcs7strip(decryptedData, MsgBytes)
	CheckErrToLog(err)
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
	CheckErrToLog(err)
	inbound, err := net.ListenTCP("tcp", addy)
	CheckErrToLog(err)
	rpc.Register(l)
	rpc.Accept(inbound)
}


func StartServer(name string) {
	fmt.Printf("Starting Server %v...\n", name)
	listener := ServerListener{name}
	listener.listenToMyAddress()
}
