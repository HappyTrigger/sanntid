package manager

import (
	".././utilities"
	"log"
	"time"
	//".././dummydriver"
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
	elevatorOrderComplete<-chan driver.OrderEvent) {
	


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

	elevatorState.Ip = localIP
	stateMap[elevatorState.Ip]=elevatorState


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



		case msg := <-reciveOrderFromPeers:
			orderMap[msg.Checksum]=msg
			sendAckToPeers<-utilities.Achnowledgement{Ip:localIP, Checksum: msg.Checksum }
			log.Println("Recived order from network")

			//One should probably store the orders with the time that they were recived
			// so that we can itirate over them and see if any orders have not been completed 
			// after some time. 

			// Or we could have this check in each elevator which sends out an emergency signal if it should be active, 
			// but isnt registering any state changes, some sort of currently active 




			//Do some calculations on elevator states here, send order to elevator if this elevator is best suited.
			if orderDelegated(stateMap,msg,currentPeers){
				SendOrderToElevator<-msg
			}
			


		case state:= <-recvStateFromPeers:
				if state.Ip == localIP && state.SenderIp != localIP{ // fix so that when you recive a state update from yourself, you dont iterate over the list
					log.Println("Recieved state from antoher elevator")

					for _,order := range state.InternalOrders{
						SendOrderToElevator<-order
						elevatorState.InternalOrders[order.Floor]=order

					} 
				}

				stateMap[state.Ip]=state







		case p := <-peerUpdateCh:
			log.Printf("Peer update:\n")
			log.Printf("  Peers:    %q\n", p.Peers)
			log.Printf("  New:      %q\n", p.New)
			log.Printf("  Lost:     %q\n", p.Lost)
			
			currentPeers = p.Peers

			// If elevator reconnects, send the saved state to it
			


			if val, ok := stateMap[p.New]; ok {
						val.ConnectionStatus = true
    					sendStateToPeers<-val
				}

			for _,v:= range p.Lost {
				//Then a connection has been lost and should be dealt with
				if elevator, ok := stateMap[v]; ok {
    				elevator.ConnectionStatus = false 	


    				//Here one should iterate over a map with all the current orders of the given elevator
					// and send each one into the distribution-func 
					}
				}



			// if peers are lost, iterate over the order map and send each order currently active on the lost elevator into
			//orderdelegation-func with the updated peers-map. 

			//if peer is new, see if it has already been in the system, if it has then push the state stored here. 



		case orderComplete:=<- recOrderCompleteFromPeers:
			log.Println("Order at Floor:",orderComplete.Floor," Complete")

			delete(orderMap,orderComplete.Checksum)


			
		case orderComplete:=<-elevatorOrderComplete:
			switch orderComplete.Button{

				case driver.Internal:
					//elevatorState.InternalOrders[orderComplete.Floor]=0


				default:
					sendOrderCompleteToPeers<-orderComplete
					delete(orderMap,orderComplete.Checksum)
				}

			
			



		case event:=<-DriverEvent:
			switch event.Button{

				case driver.Internal:	
					log.Println("internal order")
					SendOrderToElevator<-event

					elevatorState.InternalOrders[event.Floor]=event
					sendStateToPeers<-elevatorState



				default: 
					event.Checksum = event.Floor*10 + int(event.Button)
					sendOrderToPeers<-event
					unconfirmedOrderMap[event.Checksum]=event
			}


		case ack:=<-recvAckFromPeers:
			delete(unconfirmedOrderMap,ack.Checksum)
			// some check for both IPs must be implemented here before the order is deleted


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

