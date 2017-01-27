package udp

import (
	"bytes"
	"log"
	"net"
)

const (
	broadcastAddress = "255.255.255.255:10001"
	listenPort       = ":10002"
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

func Init(localIp string) (chan<- []byte, <-chan RawMessage) {
	fmt.Println("Inne i gorutine")
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

	log.Println("UDP initialized")

	udpBroadcastMsg, udpRecvMsg := make(chan []byte), make(chan RawMessage)

	go func() {
		fmt.Println("Gorutine started in udp.init")
		for {
			select {
			case msg := <-udpBroadcastMsg:
				broadcastChan <- msg
				fmt.Println("sender beskjed")
			case rawMsg := <-recieveChan:
				fmt.Println("Beskjed mottat fra:",rawMsg.Ip)
				if rawMsg.Ip != localIp {
					fmt.Println("Beskjed sendt videre")
					udpRecvMsg <- rawMsg

				}
			}
		}
	}()

	return udpBroadcastMsg, udpRecvMsg
}
