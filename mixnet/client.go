package mixnet

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
)


type Message struct {
	From string
	Data string
}


func StartClient(name string) {
	fmt.Printf("Starting Client %v...\n", name)
	serverAddress := usersMap["101"][AddressSpot]
	client, err := rpc.Dial("tcp", serverAddress)
	checkErr(err)
	in := bufio.NewReader(os.Stdin)
	for {
		line, _, err := in.ReadLine()
		checkErr(err)
		msg := Message{name, string(line)}
		var reply Reply
		err = client.Call("MediatorListener.GetMessage", msg, &reply)
		checkErr(err)
		log.Printf("Reply: %v, From: %v, Data: %v", reply, reply.From, reply.Data)
	}
}
