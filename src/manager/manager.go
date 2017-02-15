package manager

import (
	".././utilities"
	"log"
	"time"
	".././mydriver"
	".././network/bcast"
	".././network/localip"
	".././network/peers"
	"os"
	"fmt"

)

var AddOrder utilities.NewOrder
var messageId int

func Init(ExternalOrdersMap map[utilities.NewOrder]int) {
	//Initializes the Map of external orders
}

//network channels
func Run(SendOrderToElevator chan<- driver.OrderEvent,
	DriverEvent <-chan driver.OrderEvent,
	DoorOpen <- chan bool,
	DoorClosed <-chan bool,
	ElevatorEmergency <-chan bool,
	elevatorOrderComplete<-chan utilities.NewOrder) {
	


	var id string
	var checksum int
	var localIP string

	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			log.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}






	orderMap := make(map[int]utilities.NewOrder)
	unconfirmedOrderMap := make(map[int]utilities.NewOrder)
	stateMap := make(map[string]utilities.State)
	orderResend := time.Tick(2*time.Second)




	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	go peers.Transmitter(30201, id, peerTxEnable)
	go peers.Receiver(30201, peerUpdateCh)

	sendOrderToPeers :=make(chan utilities.NewOrder)
	reciveOrderFromPeers := make(chan utilities.NewOrder)
	go bcast.Transmitter(30202, sendOrderToPeers)
	go bcast.Receiver(30202, reciveOrderFromPeers)

	sendOrderCompleteToPeers := make(chan utilities.NewOrder)
	recOrderCompleteFromPeers := make(chan utilities.NewOrder)
	go bcast.Transmitter(30203, sendOrderCompleteToPeers)
	go bcast.Receiver(30203, recOrderCompleteFromPeers)

	sendStateToPeers := make(chan utilities.State)
	recvStateFromPeers := make(chan utilities.State)
	go bcast.Transmitter(30204, sendStateToPeers)
	go bcast.Receiver(30204, recvStateFromPeers)

	sendAckToPeers:= make(chan utilities.Achnowledgement)
	recvAckFromPeers := make(chan utilities.Achnowledgement)
	go bcast.Transmitter(30205, sendAckToPeers)
	go bcast.Receiver(30205, recvAckFromPeers)






	for {

		select {
		case msg := <-reciveOrderFromPeers:
			order:= driver.OrderEvent{Floor:msg.Floor,
				Button:msg.Button, OrderId: msg.Checksum}

			orderMap[msg.Checksum]=msg
			sendAckToPeers<-utilities.Achnowledgement{Ip:localIP, Checksum: msg.Checksum }

			log.Println("Recived order from network")
			
			SendOrderToElevator<-order


		case state:= <-recvStateFromPeers:
				stateMap[state.Ip]=state


		case orderComplete:=<- recOrderCompleteFromPeers:
			log.Println("Order at Floor:",orderComplete.Floor," Complete")
			delete(orderMap,orderComplete.Checksum)
			log.Println(orderMap)


		case order:=<-elevatorOrderComplete:
			sendOrderCompleteToPeers<-order




		case p := <-peerUpdateCh:
			log.Printf("Peer update:\n")
			log.Printf("  Peers:    %q\n", p.Peers)
			log.Printf("  New:      %q\n", p.New)
			log.Printf("  Lost:     %q\n", p.Lost)
			


		case event:=<-DriverEvent:
			checksum = event.Floor*10 + int(event.Button)
			newOrder:= utilities.NewOrder{Floor: event.Floor,Button: event.Button, Checksum:checksum }
			sendOrderToPeers<-newOrder
			unconfirmedOrderMap[newOrder.Checksum]=newOrder


		case ack:=<-recvAckFromPeers:
			delete(unconfirmedOrderMap,ack.Checksum)


		case <-orderResend:
			for _,v:=range unconfirmedOrderMap{
				sendOrderToPeers<-v
			}
			
		}
	}
}

