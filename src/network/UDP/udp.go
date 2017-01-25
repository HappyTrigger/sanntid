package main

import (
	"bytes"
	"log"
	"net"
	"time"
	"fmt"
	"strconv"
)

const (
	broadcastAddress = "255.255.255.255:10001"
	listenPort       = ":30000"
)

type RawMessage struct {
	Data []byte
	Ip   string
}

func recieve(recieveChan chan<- RawMessage, broadcastListener *net.UDPConn) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error in UDP recieve: %s \n Closing connection.", r)
			broadcastListener.Close()
		}
	}()

	buffer := make([]byte, 1024)

	for {
		n, address, err := broadcastListener.ReadFromUDP(buffer)

		if err != nil || n < 0 {
			log.Printf("Error in UDP recieve\n")
			panic(err)
		}

		data, err := bytes.NewBuffer(buffer).ReadBytes('\n')

		if err != nil {
			log.Println("Error when reading UDP message buffer:", err)
		}

		recieveChan <- RawMessage{Data: data, Ip: address.IP.String()}
	}
}

func broadcast(broadcastChan <-chan []byte, localListener *net.UDPConn) {

	addr, _ := net.ResolveUDPAddr("udp", broadcastAddress)

	var b bytes.Buffer

	for msg := range broadcastChan {
		b.Write(msg)
		b.WriteRune('\n')

		_, err := localListener.WriteToUDP(b.Bytes(), addr)

		b.Reset()

		if err != nil {
			log.Println(err)
		}
	}
}
var localIp string
func main(){
	localIp = "127.0.0.1"
	
	addr, _ := net.ResolveUDPAddr("udp", listenPort)

	localListener, err := net.ListenUDP("udp", addr)

	if err != nil {
		log.Fatal(err)
	}

	addr, _ = net.ResolveUDPAddr("udp", broadcastAddress)

	broadcastListener, err := net.ListenUDP("udp", addr)

	if err != nil {
		log.Fatal(err)
	}

	broadcastChan := make(chan []byte)
	go broadcast(broadcastChan, localListener)

	recieveChan := make(chan RawMessage)
	go recieve(recieveChan, broadcastListener)

	//msg2:=RawMessage{}
	log.Println("UDP initialized")
   	//buf22 := make(RawMessage)
	udpBroadcastMsg, udpRecvMsg := make(chan []byte), make(chan RawMessage)


	go func() {
		for {
			select {
			case msg := <-udpBroadcastMsg:
				broadcastChan <- msg
			case rawMsg := <-recieveChan:

				//if rawMsg.Ip != localIp {
				udpRecvMsg <- rawMsg

				//}
			}
		}
	}()
	i := 3
	for{
		msg := strconv.Itoa(i)
		i++
		buf := []byte(msg)

		time.Sleep(3*time.Second)

		udpBroadcastMsg<-buf

		msg2:=udpRecvMsg
		fmt.Printf("%+v\n", msg2)
		}
	

	//return udpBroadcastMsg, udpRecvMsg
}



