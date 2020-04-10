package main

import (
	"crypto/elliptic"
	"encoding/gob"
	"github.com/EladCoding/HideMetaData/mixnet"
	"os"
	"os/exec"
	"time"
)


func runWholeMixNet() {
	go mixnet.StartUser("server", "001")
	time.Sleep(100*time.Millisecond)
	go mixnet.StartUser("server", "002")
	time.Sleep(100*time.Millisecond)
	go mixnet.StartUser("server", "003")
	time.Sleep(100*time.Millisecond)
	go mixnet.StartUser("mediator", "103")
	time.Sleep(100*time.Millisecond)
	go mixnet.StartUser("mediator", "102")
	time.Sleep(100*time.Millisecond)
	go mixnet.StartUser("mediator", "101")
	time.Sleep(100*time.Millisecond)
	mixnet.StartUser("client", "201")
}


func runWholeMixNetInADiffrentProcesses() {
	var cmd *exec.Cmd
	cmd = exec.Command("runningExample.exe", "server 001")
	go cmd.Run()
	time.Sleep(time.Second)
	cmd = exec.Command("runningExample.exe", "server 002")
	go cmd.Run()
	time.Sleep(time.Second)
	cmd = exec.Command("runningExample.exe", "server 003")
	go cmd.Run()
	time.Sleep(time.Second)
	cmd = exec.Command("runningExample.exe", "mediator 103")
	go cmd.Run()
	time.Sleep(time.Second)
	cmd = exec.Command("runningExample.exe", "mediator 102")
	go cmd.Run()
	time.Sleep(time.Second)
	cmd = exec.Command("runningExample.exe", "mediator 101")
	go cmd.Run()
	time.Sleep(time.Second)
	mixnet.StartUser("client", "201")
}


func runOneNode(mode string, name string) {
	mixnet.StartUser(mode, name)
}


func main() {
	gob.Register(elliptic.P256())
	//convertStructToBytes()
	mode := os.Args[1]
	name := os.Args[2]
	runOneNode(mode, name)
	//runWholeMixNet()
	//scripts.Main()
}
