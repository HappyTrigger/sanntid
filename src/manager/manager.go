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
	orderAssignedToMap := make(map[int]string) // combine the checksum and IP of the given elevator
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



			takeorder:=false
			//Do some calculations on elevator states here, send order to elevator if this elevator is best suited.
			if takeorder,orderAssignedToMap = orderDelegated(stateMap,msg,currentPeers,orderAssignedToMap); takeorder{
				SendOrderToElevator<-msg
			}
			


		case state:= <-recvStateFromPeers:
				if state.Ip == localIP && state.StateSentFromIp != localIP{ // fix so that when you recive a state update from yourself, you dont iterate over the list
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


			if state, ok := stateMap[p.New]; ok {
						stateMap[p.New]=state
						state.StateSentFromIp = localIP
    					sendStateToPeers<-state
				}else{
					stateMap[p.New]=state
				}
			takeorder := false //rewrite this, looks ugly as fuck
			for _,lostIp:= range p.Lost {
					for checksum,ip:= range orderAssignedToMap{
						if ip==lostIp{
							msg := orderMap[checksum]
							if takeorder,orderAssignedToMap = orderDelegated(stateMap,msg,currentPeers,orderAssignedToMap); takeorder{
								SendOrderToElevator<-msg
							}

						}
					}
				}




		case orderComplete:=<- recOrderCompleteFromPeers:
			log.Println("Order at Floor:",orderComplete.Floor," Complete")

			delete(orderAssignedToMap,orderComplete.Checksum)
			delete(orderMap,orderComplete.Checksum)


			
		case orderComplete:=<-elevatorOrderComplete:
			switch orderComplete.Button{

				case driver.Internal:
					var internalOrders [] driver.OrderEvent
					delete(orderMap,orderComplete.Checksum)
					// Loop threw ordermap and put every internal-order into the elevator.internalorder slice then send the state update
					
					for _, order := range orderMap{
						if order.Button == driver.Internal{
							internalOrders = append(internalOrders,order)
						}

					}
					elevatorState.InternalOrders = internalOrders
					sendStateToPeers<-elevatorState


				default:
					sendOrderCompleteToPeers<-orderComplete
					delete(orderAssignedToMap,orderComplete.Checksum)
					delete(orderMap,orderComplete.Checksum)

				}

			
			



		case event:=<-DriverEvent:
			event.Checksum = event.Floor*10 + int(event.Button)
			
			switch event.Button{

				case driver.Internal:
					var internalOrders [] driver.OrderEvent
					log.Println("internal order")
					orderMap[event.Checksum]=event
					// Loop threw ordermap and put every internal-order into the elevator.internalorder slice then send the state update
					//sendStateToPeers<-elevatorState

					for _, order := range orderMap{
						if order.Button == driver.Internal{
							internalOrders = append(internalOrders,order)
						}

					}
					elevatorState.InternalOrders = internalOrders
					sendStateToPeers<-elevatorState
					SendOrderToElevator<-event




				default: 
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






//This function should totally be rewritten
func orderDelegated(stateMap map[string]utilities.State,
	orderEvent driver.OrderEvent,currentPeers[]string,orderAssignedToMap map[int]string) (bool,map[int]string) {


	var fitness int
	fitnessMap := make(map[string]int)
	var i int

	// The elevator with the lowest fitness takes the order

	for elevator,state := range stateMap{
		for _,peer:= range currentPeers{
			if elevator == peer{
				if state.BetweenFloors{
					fitness += i

				}else{
					fitness = 1000 + i
					i++
				}
				fitnessMap[state.Ip]=fitness
			}
		}
		fitness=0
	}

	//Create some sort of delegation-algorithm, which bases its decision on the current active peers, 
	// the status of the peers(Are they moving or idle) and the postition and direction of the active elevators
	var maxValue int
	var OrderGivenToIp string
	for ip,value := range fitnessMap{
		if value > maxValue{
			maxValue=value
			OrderGivenToIp = ip
		}
	}


	orderAssignedToMap[orderEvent.Checksum]=OrderGivenToIp

	if OrderGivenToIp == localIP{
		return true, orderAssignedToMap
	}
	return false, orderAssignedToMap
}

