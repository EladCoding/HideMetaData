package main

import (
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
	switch runningChoice {
	case "1":
		scripts.RunningExample()
	default:
		msg := "Please insert your running choice:\n1 for running example\n"
		fmt.Printf(msg)
	}
}
