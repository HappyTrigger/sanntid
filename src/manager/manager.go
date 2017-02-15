package manager

import (
	".././utilities"
	"log"
	//"time"
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
	ElevatorEmergency <-chan bool) {
	
	var id string

	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			log.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}



	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)


	sendOrderToPeers :=make(chan utilities.NewOrder)
	reciveOrderFromPeers := make(chan utilities.NewOrder)

	sendOrderComplete := make(chan utilities.OrderComplete)
	recOrderComplete := make(chan utilities.OrderComplete)


	go bcast.Transmitter(16569, sendOrderToPeers)
	go bcast.Receiver(16569, reciveOrderFromPeers)

	go bcast.Transmitter(16570, sendOrderComplete)
	go bcast.Receiver(16570, recOrderComplete)


	for {
		select {
		case msg := <-reciveOrderFromPeers:
			
			order:= driver.OrderEvent{Floor:msg.Floor,
				Button:msg.Button}

				log.Println("Recived order from network")
				SendOrderToElevator<-order

		case orderComplete:=<- recOrderComplete:
			log.Println("Order at Floor:",orderComplete.Floor," Complete")

		case p := <-peerUpdateCh:
			log.Printf("Peer update:\n")
			log.Printf("  Peers:    %q\n", p.Peers)
			log.Printf("  New:      %q\n", p.New)
			log.Printf("  Lost:     %q\n", p.Lost)
			


		case event:=<-DriverEvent:
			newOrder:= utilities.NewOrder{Floor: event.Floor,Button: event.Button }
			sendOrderToPeers<-newOrder
			
		}
	}
}

