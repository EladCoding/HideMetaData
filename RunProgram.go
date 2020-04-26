package main

import (
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"github.com/EladCoding/HideMetaData/scripts"
	"os"
)

func main() {
	args := os.Args
	var runningChoice string
	if len(args) > 1 {
		runningChoice = args[1]
	}
	gob.Register(elliptic.P256())
	switch runningChoice {
	case "1":
		scripts.RunningExample()
	case "2":
		scripts.RunStatistics()
	case "3":
		scripts.RunInfrastructure()
	case "4":
		if len(os.Args) < 4 {
			msg := "please choose mode and name\n"
			fmt.Printf(msg)
		} else {
			mode := os.Args[2]
			name := os.Args[3]
			scripts.RunNode(mode, name)
		}
	default:
		msg := "Please insert your running choice:\n1 for running example\n2 for running statistics\n" +
			"3 to run infraStructure\n4 to run one node\n"
		fmt.Printf(msg)
	}
}
