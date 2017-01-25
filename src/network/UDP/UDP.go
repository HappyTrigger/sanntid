package main

import (
	"net"
	"strings"
	"encoding/json"
	"fmt"

	)
const (
	PORT = ":20022"
	PORTLISTENER = ":30000"
	)

type Message struct {
	Source int
	Id int
	Floor int
	Target int
	Checksum int
}

func UDPSender(channel<-chan Message) {
	broadcastAddr := []string{"129.241.187.255", PORT}
	broadcastUDP, _ := net.ResolveUDPAddr("udp", strings.Join(broadcastAddr, ""))
	broadcastConn, _ := net.DialUDP("udp", nil, broadcastUDP) //returns the UDP connection interface which supports reading and writing
	defer broadcastConn.Close() // Plan the connection to be closed soon
	for {
		buf, err := json.Marshal(<- channel) //codes the message
		if err == nil {
			broadcastConn.Write(buf)
		}
	}
}

func UDPListener(channel<-chan Message) {
	UDPReceiveAddr, err := net.ResolveUDPAddr("udp", PORTLISTENER);
	if err != nil { fmt.Println(err) }

	UDPConn, err := net.ListenUDP("udp", UDPReceiveAddr);
	if err != nil { fmt.Println(err) }
	defer UDPConn.Close()

	buf := make([]byte, 2048)
	trimmed_buf := make([]byte, 1)
	var received_message Message

	for {
		n, _, _ := UDPConn.ReadFromUDP(buf)
		trimmed_buf = buf[:n]
		err := json.Unmarshal(trimmed_buf, &received_message) //decodes the message
		if err == nil {
			channel <- received_message
		}
	}
}


func main() {
	message := make(chan Message)
	

	

	go UDPSender(message)
}