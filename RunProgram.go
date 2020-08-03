package main

import (
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"github.com/EladCoding/HideMetaData/mixnet"
	"github.com/EladCoding/HideMetaData/scripts"
	"log"
	"os"
	"strconv"
	"time"
)

// Initialize mixnet architecture and external files to prepare for a simulation.
func prepareRun() {
	mixnet.InitLogFile()
	scripts.CreateNodesMap()
	mixnet.UserAddressesMap = mixnet.ReadUserAddressMap()
	mixnet.UserPubKeyMap = mixnet.ReadUserPubKeyMap()
	mixnet.UserPrivKeyMap = mixnet.ReadUserPrivKeyMap()
	gob.Register(elliptic.P256())
}

// Run the whole mixnet architecture simulation.
func main() {
	prepareRun()
	args := os.Args
	var runningChoice string
	if len(args) > 1 {
		runningChoice = args[1]
	}
	if len(args) > 2 {
		if i, err := strconv.ParseInt(args[2], 10, 64); err == nil {
			mixnet.RoundSlotTime = time.Duration(i) * time.Second
		} else {
			fmt.Printf("Round slot time is not valid!\n")
			return
		}
	}
	switch runningChoice {
	case "1":
		if len(args) <= 2 {
			mixnet.RoundSlotTime = time.Second
		}
		fmt.Printf("--------------------Run Automatic Tests--------------------\n")
		log.Printf("--------------------Run Automatic Tests--------------------\n")
		scripts.RunAutomaticTests()
	case "2":
		fmt.Printf("--------------------Run Statistics--------------------\n")
		log.Printf("--------------------Run Statistics--------------------\n")
		scripts.RunStatistics()
	case "3":
		fmt.Printf("--------------------Run Playing example--------------------\n")
		log.Printf("--------------------Run Playing example--------------------\n")
		scripts.PlayingExample()
	case "4":
		fmt.Printf("--------------------Run InfraStructure--------------------\n")
		log.Printf("--------------------Run InfraStructure--------------------\n")
		scripts.RunInfrastructure()
	case "5":
		fmt.Printf("--------------------Run One node--------------------\n")
		log.Printf("--------------------Run One node--------------------\n")
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
			"1 Run Automatic Tests.\n" +
			"2 Run statistics.\n" +
			"3 Run Playing example.\n" +
			"4 Run infraStructure.\n" +
			"5 Run one node.\n\n" +
			"For changing round slot time - pass the new round slot seconds as a third parameter.\n"
		fmt.Printf(msg)
	}
}
