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


const(
	OrderResendInterval = 500*time.Millisecond
)
	var localIP string



func Run(SendOrderToElevator chan<- driver.OrderEvent,
	DriverEvent <-chan driver.OrderEvent,
	DoorOpen <- chan bool,
	DoorClosed <-chan bool,
	ElevatorEmergency <-chan bool,
	elevatorOrderComplete<-chan utilities.NewOrder) {
	


	var id string
	var elevatorState utilities.State
	var currentPeers []string


	orderMap := make(map[int]driver.OrderEvent)
	unconfirmedOrderMap := make(map[int]driver.OrderEvent)
	stateMap := make(map[string]utilities.State)
	orderResend := time.Tick(OrderResendInterval)



	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			log.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}



	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	go peers.Transmitter(30201, id, peerTxEnable)
	go peers.Receiver(30201, peerUpdateCh)

	sendOrderToPeers :=make(chan driver.OrderEvent)
	reciveOrderFromPeers := make(chan driver.OrderEvent)
	go bcast.Transmitter(30202, sendOrderToPeers)
	go bcast.Receiver(30202, reciveOrderFromPeers)

	sendOrderCompleteToPeers := make(chan driver.OrderEvent)
	recOrderCompleteFromPeers := make(chan driver.OrderEvent)
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

			//Should probably rewrite this  so that we can just send and recive a driver.Eventype rather than converting it

		case msg := <-reciveOrderFromPeers:
			orderMap[msg.Checksum]=msg
			sendAckToPeers<-utilities.Achnowledgement{Ip:localIP, Checksum: msg.Checksum }
			log.Println("Recived order from network")

			//One should probably store the orders with the time that they were recived
			// so that we can itirate over them and see if any orders have not been completed 
			// after some time. This way we can 

			// Or we could have this check in each elevator which sends out an emergency signal if it should be active, 
			// but isnt registering any state changes




			//Do some calculations on elevator states here, send order to elevator if this elevator is best suited.
			if orderDelegated(stateMap,msg,currentPeers){
				SendOrderToElevator<-order
			}
			


		case state:= <-recvStateFromPeers:
				stateMap[state.Ip]=state


		case orderComplete:=<- recOrderCompleteFromPeers:
			log.Println("Order at Floor:",orderComplete.Floor," Complete")
			delete(orderMap,orderComplete.Checksum)
			log.Println(orderMap)


		case orderComplete:=<-elevatorOrderComplete:
			delete(orderMap,orderComplete.Checksum)
			sendOrderCompleteToPeers<-orderComplete




		case p := <-peerUpdateCh:
			log.Printf("Peer update:\n")
			log.Printf("  Peers:    %q\n", p.Peers)
			log.Printf("  New:      %q\n", p.New)
			log.Printf("  Lost:     %q\n", p.Lost)
			
			currentPeers = p.Peers

			




		case event:=<-DriverEvent:
			switch event.Button{

				case driver.Internal:	
					log.Println("internal order")

					//Append the new internal order to the internal order map, then send the new state

				default: 
					event.Checksum = event.Floor*10 + int(event.Button)
					sendOrderToPeers<-event
					unconfirmedOrderMap[newOrder.Checksum]=newOrder
			}


		case ack:=<-recvAckFromPeers:
			delete(unconfirmedOrderMap,ack.Checksum)


		case <-orderResend:
			for _,v:=range unconfirmedOrderMap{
				sendOrderToPeers<-v
			}
			
		}
	}
}

//order delegation function based on the states this and the other current active elevators
func orderDelegated(stateMap map[string]utilities.State,
	orderEvent driver.OrderEvent,currentPeers[]string) bool {



	// The elevator with the lowest fitness takes the order

/*
	fitness := make(map[string]int)
	for peer,state := range stateMap{




	}

*/	
	//Create some sort of delegation-algorithm, which bases its decision on the current active peers, 
	// the status of the peers(Are they moving or idle) and the postition and direction of the active elevators





	//returns true if this elevator should take the order, else return false
	return true
}

