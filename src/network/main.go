package main


import (

	"./UDP"
	"net"
    "os"
    "fmt"


)
var name string
var localip string
var err bool

func main() {
	name, err := os.Hostname()
	localip, err := net.LookupHost(name)
	udpBroadcastMsg,udpRecvMsg:=udp.Init(localip)
	fmt.Println(localip)
	if err!=nil{
		fmt.Println("jei")
	}
	
}