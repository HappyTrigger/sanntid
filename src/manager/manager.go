package manager

import (
	".././utilities"
	"log"
	"time"
	//".././dummydriver"
	".././driver"
	".././network/bcast"
	".././network/localip"
	".././network/peers"
	"os"
	"fmt"
	"math"

)


const(
	OrderResendInterval = 5*time.Millisecond

)
	var localIP string
	var err error
	var currentElevatorState utilities.State


// Some init might be necessary to check if the elevator is re-initializing



func Run(SendOrderToElevator chan<- driver.OrderEvent,
	DriverEvent <-chan driver.OrderEvent,
	ElevatorEmergency <-chan bool,
	ElevatorOrderComplete<-chan driver.OrderEvent,
	ElevatorStateFromElevator <-chan utilities.State) {
	

	var id string
	var currentPeers []string

	localIP, err = localip.LocalIP()
	if err != nil {
		log.Println(err)
		localIP = "DISCONNECTED"
	}
	id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	orderMap 				:= make(map[int]driver.OrderEvent)
	orderAssignedToMap 		:= make(map[int]string) // combine the checksum and IP of the given elevator
	unconfirmedOrderMap 	:= make(map[int]driver.OrderEvent)
	achnowledgementMap		:= make(map[int][]string)
	stateMap 				:= make(map[string]utilities.State)
	orderResend 			:= time.Tick(OrderResendInterval)


	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	//PeerTxenable
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

	log.Println("Starting")
	log.Println("Local Ip : ", localIP)

//Test
	/*
	go func() {
		//time.Sleep(3*time.Second)
		//reciveOrderFromPeers <- driver.OrderEvent{3, driver.ButtonType(driver.Down),0}
		 for {
		 	time.Sleep(5*time.Second)
			sendOrderToPeers <- driver.OrderEvent{3, driver.ButtonType(driver.Down),0}

		 }
	}()
	currentPeers = append(currentPeers, localIP)
*/

	for {

		select {

		case msg := <-reciveOrderFromPeers:
			sendAckToPeers<-utilities.Achnowledgement{Ip:localIP, Checksum: msg.Checksum }
			if _,orderExist := orderMap[msg.Checksum]; !orderExist{ //If order alread exist, dont process it
				orderMap[msg.Checksum]=msg
				driver.Elev_set_button_lamp(msg.Button,msg.Floor,true) // This must be set on every elevator as it is a requirement
				if ok := OrderDelegator(stateMap,msg,currentPeers,orderAssignedToMap); ok{
					SendOrderToElevator<-msg
				}
			}
			

		case state:= <-recvStateFromPeers:
				if state.Ip == localIP && state.StateSentFromIp != localIP{ 
					log.Println("Internal orders recieved from : ", state.StateSentFromIp)
					for _,internalOrder := range state.InternalOrders{
						SendOrderToElevator<-internalOrder

					} 
				}else{
					stateMap[state.Ip]=state
				}
				
				

		case p := <-peerUpdateCh:
			log.Printf("Peer update:\n")
			log.Printf("  Peers:    %q\n", p.Peers)
			log.Printf("  New:      %q\n", p.New)
			log.Printf("  Lost:     %q\n", p.Lost)
			
			currentPeers = p.Peers
			sendStateToPeers<-currentElevatorState // In case reconnection, send state before new orders are recived
			
			if state, ok := stateMap[p.New]; ok { 
				log.Println("Reconnecting elevator, sending internal Orders")
				state.StateSentFromIp = localIP
    			sendStateToPeers<-state
			}
			

			for _,lostIp:= range p.Lost {
					for checksum,ip:= range orderAssignedToMap{
						if ip==lostIp{
							msg := orderMap[checksum]
							if ok := OrderDelegator(stateMap,msg,currentPeers,orderAssignedToMap); ok{
								SendOrderToElevator<-msg
							}

						}
					}
				}


		case state:= <-ElevatorStateFromElevator:
			state.Ip,state.StateSentFromIp = localIP,localIP
			stateMap[localIP]=state
			currentElevatorState = state
			sendStateToPeers<-state


		case orderComplete:=<- recOrderCompleteFromPeers:
			log.Println("Order at Floor:",orderComplete.Floor," completed")
			delete(orderAssignedToMap,orderComplete.Checksum)
			delete(orderMap,orderComplete.Checksum)


		case orderComplete:=<-ElevatorOrderComplete:
			delete(orderMap,orderComplete.Checksum)
			switch orderComplete.Button{
				case driver.Internal:
					log.Println("Internal Order complete")
					//could use a map here, iterate over it and send the internalOrder list to the other elevators

				default: 
					sendOrderCompleteToPeers<-orderComplete
					delete(orderAssignedToMap,orderComplete.Checksum)
			}





		case event:=<-DriverEvent:
			event.Checksum = event.Floor*10 + int(event.Button)
			switch event.Button{
				case driver.Internal:
					SendOrderToElevator<-event
					//Should probably rewrite this
					driver.Elev_set_button_lamp(event.Button,event.Floor,true)
					orderMap[event.Checksum]=event
					// implement for loop here, send internal orders



				default: 
					log.Println("sending order")
					sendOrderToPeers<-event
					unconfirmedOrderMap[event.Checksum]=event
			}


		case ack:=<-recvAckFromPeers:
			var achnowledgeIteration int
			alreadyExist := false

			achnowledgelist := achnowledgementMap[ack.Checksum]
			for _,ip:= range achnowledgelist{
				if ip == ack.Ip {
					alreadyExist=true
				}
			}
			if !alreadyExist{
				achnowledgelist = append(achnowledgelist, ack.Ip)
				achnowledgementMap[ack.Checksum] = achnowledgelist
			}

			for _,ip:= range achnowledgelist{
				for _,peer:= range currentPeers{
					if ip==peer{
						achnowledgeIteration++
					}
				}
			}
			if achnowledgeIteration>=len(currentPeers){
				log.Println("Recieced achnowledge")
				delete(unconfirmedOrderMap,ack.Checksum)
				delete(achnowledgementMap, ack.Checksum)
			}



			

		case <-ElevatorEmergency:
			log.Println("-------------------------------")
			log.Println("ElevatorEmergency is now active")
			log.Println("-------------------------------")
			fmt.Sprintf("peer-%s-%d", localIP, id)
			peerTxEnable <- false

			//Should start some re-initalize of the elevator here.


		case <-orderResend:
			for _,v:=range unconfirmedOrderMap{
				sendOrderToPeers<-v
			}
			
		}
	}
}




// Rather than relying on elevator-events to spread states, we send it continously.
//Need to create a new state channel that is directly linked to this function.
func continousStateSender(sendToPeers chan<- utilities.State, recieveNewStateFromManager <- chan utilities.State) {
	stateSender := time.Tick(1*time.Second)
	var elevatorState utilities.State 
	for{
		select{
		case <- stateSender:
			sendToPeers<- elevatorState
		case state :=<-recieveNewStateFromManager:
			elevatorState = state
			sendToPeers<- elevatorState

		}
	}
}




//Might have to change some of the values
//This functions as it should. Can see if we want to rewrite it
// into several small functions instead, splitting up the cost-functions
// and combining them in a larger function, just to clean up the code.


func OrderDelegator(stateMap map[string]utilities.State,
	orderEvent driver.OrderEvent,currentPeers[]string,orderAssignedToMap map[int]string) bool {

	//The elevator with the lowest fitness takes the order
	fitnessMap:= make(map[string]float64)

	for elevator,state := range stateMap{
		for _,peer:= range currentPeers{
			if elevator == peer{
				if !state.BetweenFloors{
					fitnessMap[elevator]=math.Abs(float64(state.LastPassedFloor-orderEvent.Floor))
				}else{
					floorDifference:=float64(orderEvent.Floor-state.LastPassedFloor)
					
					switch orderEvent.Button{
					case driver.Up:
						if state.Direction == driver.Up{							
							if floorDifference >=0 { //order Up above the elevator, and elavtor moving up
								fitnessMap[elevator]=float64(floorDifference)
							}else{
								//order Up below, and elevator moving up
								fitnessMap[elevator] = float64(math.Abs(float64(floorDifference)) + float64((driver.N_FLOORS-state.LastPassedFloor)*2))
							}
						

						}else{
							//Order Up bellow, and eleavtor moving down 
							if floorDifference <=0 {
								fitnessMap[elevator]= float64(state.LastPassedFloor*2) - math.Abs(float64(floorDifference))
							}else{
								//Order Up above and eleavtor moving down
								fitnessMap[elevator] = float64(math.Abs(float64(floorDifference)) + float64(state.LastPassedFloor*2))
							}
						}

					case driver.Down: 
						if state.Direction == driver.Down{
							
							if floorDifference <=0 { // Order downwards bellow, and elevator moving down
								fitnessMap[elevator]=math.Abs(float64(floorDifference))
							}else{
								//Order downwards above, and elevator moving down
								fitnessMap[elevator] = float64(math.Abs(float64(floorDifference)) + float64(state.LastPassedFloor*2))
							}
						}else{
							//Order downwards above and elevator moving up
							if floorDifference >=0 {
								fitnessMap[elevator]=float64((driver.N_FLOORS-state.LastPassedFloor)*2) - float64(floorDifference)
							}else{
								//Order downwards bellow and elevator moving up
								fitnessMap[elevator] = float64(math.Abs(float64(floorDifference)) + float64((driver.N_FLOORS-state.LastPassedFloor)*2))
							}
						}
					}
				}
			}
		}
	}
	var minFitness float64
	minFitness = 20
	var ip string
	for elevator, fitness := range fitnessMap{
		if fitness == minFitness{
			if elevator > ip{
				ip = elevator
			} // Elevator with highest IP takes the order
		}
		if fitness < minFitness{
			minFitness = fitness
			ip = elevator
		}
	}
	log.Println("Order assigned to :", ip)
	log.Println("Local Ip: ", localIP)

	log.Println("------CurrentPeers------")
	log.Println(currentPeers)

	log.Println("------FitnessMap------")
	log.Println(fitnessMap)
	for k,v := range fitnessMap{
		log.Println("Ip: ",k, " - Fitness: ",v)
	}
	orderAssignedToMap[orderEvent.Checksum]=ip

	if ip == localIP{
		return true
	}else{
		log.Println("Did not take order")
		return false
	}
}







