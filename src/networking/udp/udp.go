package udp

import (
	"bytes"
	"log"
	"net"
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

func Init(localIp string) (chan<- []byte, <-chan RawMessage){

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

	broadcastChan := make(chan []byte,50)
	go broadcast(broadcastChan, localListener)

	recieveChan := make(chan RawMessage,50)
	go recieve(recieveChan, broadcastListener)

	
	log.Println("UDP initialized")
   	
	udpBroadcastMsg, udpRecvMsg := make(chan []byte,50), make(chan RawMessage,50)


	go func() {
		for {
			select {
			case msg := <-udpBroadcastMsg:
			//	log.Println("Trying to broadcast within UDP")
				broadcastChan <- msg
			//	log.Println("Completed broadcast within UDP")

			case rawMsg := <-recieveChan:
					if rawMsg.Ip != localIp {
						udpRecvMsg <- rawMsg
					}
			default:
				// 
				
			}
		}
	}()

	return udpBroadcastMsg, udpRecvMsg
}



