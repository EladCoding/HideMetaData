package main

import (
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"github.com/EladCoding/HideMetaData/mixnet"
	"github.com/EladCoding/HideMetaData/scripts"
	"os"
)

func prepareRun() {
	scripts.CreateUsersMap()
	mixnet.UserAddressesMap = mixnet.ReadUserAddressMap()
	mixnet.UserPubKeyMap = mixnet.ReadUserPubKeyMap()
	mixnet.UserPrivKeyMap = mixnet.ReadUserPrivKeyMap()
	gob.Register(elliptic.P256())
}

func main() {
	prepareRun()
	args := os.Args
	var runningChoice string
	if len(args) > 1 {
		runningChoice = args[1]
	}
	switch runningChoice {
	case "1":
		fmt.Printf("--------------------Run Running example--------------------\n")
		scripts.RunningExample()
	case "2":
		fmt.Printf("--------------------Run Statistics--------------------\n")
		scripts.RunStatistics()
	case "3":
		fmt.Printf("--------------------Run Automatic Tests--------------------\n")
		scripts.RunAutomaticTests()
	case "4":
		fmt.Printf("--------------------Run InfraStructure--------------------\n")
		scripts.RunInfrastructure()
	case "5":
		fmt.Printf("--------------------Run One node--------------------\n")
		if len(os.Args) < 4 {
			msg := "please choose mode and name\n"
			fmt.Printf(msg)
		} else {
			mode := os.Args[2]
			name := os.Args[3]
			scripts.RunNode(mode, name)
		}
	default:
		msg := "Please insert your running choice:\n" +
			"1 Run example\n" +
			"2 Run statistics\n" +
			"3 Run Automatic Tests\n" +
			"4 Run infraStructure\n" +
			"5 Run one node\n"
		fmt.Printf(msg)
	}
}
