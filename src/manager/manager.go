package manager
/*
The manager handels new and old information and executes action based on the information given. 
It tracks other connected elevators, and handels order-delegation based on the states of every single elevator.
The manager is also responsible for sending states, order-events and achnowledge messages from other elevators. 
*/

import (
	"../utilities"
	"log"
	"time"
	"../driver"
	"../network/bcast"
	"../network/localip"
	"../network/peers"
	"os"
	"fmt"
	"math"


)


const(
	OrderResendInterval = 50*time.Millisecond
	StateResendInterval = 100*time.Millisecond
	orderNotCompletedInterval = 2*time.Second

)
	var localId string
	var currentElevatorState utilities.State


// Some init might be necessary to check if the elevator is re-initializing

func Init(reciveStateFromPeers <-chan utilities.State,
	sendStateToPeers chan <- utilities.State,
	ElevatorStateFromElevator <-chan utilities.State,
	SendOrderToElevator chan<- driver.OrderEvent,
	internalOrderMap map[int]driver.OrderEvent, ) {
	
	var internalOrderSlice []driver.OrderEvent
	timeout := time.After(200*time.Millisecond)
	
	loop:
	for{
		select{
			case state:= <-reciveStateFromPeers:
				if state.Id == localId{ 
					
					for _,internalOrder := range state.InternalOrders{
						internalOrderSlice=append(internalOrderSlice, internalOrder)
						internalOrderMap[internalOrder.Checksum]=internalOrder
					}
				}
			case <-timeout:
				log.Println("First init of elevator")
				break loop
			case state:=<-ElevatorStateFromElevator:
				currentElevatorState = state 


			
		}
	}
	go func(){
		for _,orders := range internalOrderSlice{
			SendOrderToElevator<-orders
		}
	}()
	sendNewStateToPeers(sendStateToPeers,internalOrderMap)

}





func Run(SendOrderToElevator chan<- driver.OrderEvent,
	DriverEvent <-chan driver.OrderEvent,
	ElevatorEmergency <-chan bool,
	ElevatorOrderComplete<-chan driver.OrderEvent,
	ElevatorStateFromElevator <-chan utilities.State,
	Id string) {
	

	var currentPeers []string

	localId = Id

	if localId == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		localId = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}
	log.Println("Local Id : ", localId)


	orderMap 				:= make(map[int]driver.OrderEvent)
	internalOrderMap		:= make(map[int]driver.OrderEvent)
	unconfirmedOrderMap 	:= make(map[int]driver.OrderEvent)

	orderAssignedToMap 		:= make(map[int]string) 
	achnowledgementMap		:= make(map[int][]string)
	stateMap 				:= make(map[string]utilities.State)
	orderRecievedAtTime		:= make(map[int]time.Time)

	orderResend 			:= time.Tick(OrderResendInterval)
	stateResend 			:= time.Tick(StateResendInterval)
	orderNotCompleted		:= time.Tick(orderNotCompletedInterval)



	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	go peers.Transmitter(30201, localId, peerTxEnable)
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



	
	Init(recvStateFromPeers,
		sendStateToPeers,
		ElevatorStateFromElevator,
		SendOrderToElevator,
		internalOrderMap)



	for {

		select {

		case msg := <-reciveOrderFromPeers:
			sendAckToPeers<-utilities.Achnowledgement{Id:localId, Checksum: msg.Checksum }
			if _,orderExist := orderMap[msg.Checksum]; !orderExist{ //If order alread exist, dont process it
				orderRecievedAtTime[msg.Checksum]=time.Now()
				orderMap[msg.Checksum]=msg
				driver.Elev_set_button_lamp(msg.Button,msg.Floor,true) // This must be set on every elevator as it is a requirement
				if ok := OrderDelegator(stateMap,msg,currentPeers,orderAssignedToMap); ok{
					SendOrderToElevator<-msg
				}
			}
			

		case state:= <-recvStateFromPeers:
			stateMap[state.Id]=state

				
				
				

		case p := <-peerUpdateCh:
			log.Printf("Peer update:\n")
			log.Printf("  Peers:    %q\n", p.Peers)
			log.Printf("  New:      %q\n", p.New)
			log.Printf("  Lost:     %q\n", p.Lost)
			
			currentPeers = p.Peers
			
			if state, ok := stateMap[p.New]; ok {//&& state.StateSentFromId != localId{ 
				state.StateSentFromId = localId
				log.Println("Sending state to reconnecting elevator")
				log.Println(state)
    			sendStateToPeers<-state
			}
			log.Println(currentElevatorState)
			

			for _,lostId:= range p.Lost {
					for checksum,Id:= range orderAssignedToMap{
						if Id==lostId{
							msg := orderMap[checksum]
							if ok := OrderDelegator(stateMap,msg,currentPeers,orderAssignedToMap); ok{
								SendOrderToElevator<-msg
							}

						}
					}
				}


		case state:= <-ElevatorStateFromElevator:
			stateMap[localId]=state
			currentElevatorState = state
			sendNewStateToPeers(sendStateToPeers,internalOrderMap)
			


		case orderComplete:=<- recOrderCompleteFromPeers:
			log.Println("Order at Floor:",orderComplete.Floor," completed")
			delete(orderAssignedToMap,orderComplete.Checksum)
			delete(orderMap,orderComplete.Checksum)
			delete (orderRecievedAtTime, orderComplete.Checksum)
			driver.Elev_set_button_lamp(orderComplete.Button, orderComplete.Floor, false)


		case orderComplete:=<-ElevatorOrderComplete:
			switch orderComplete.Button{
				case driver.Internal:
					log.Println("Internal Order complete")
					delete(internalOrderMap,orderComplete.Checksum)

					sendNewStateToPeers(sendStateToPeers,internalOrderMap)

				default: 
					sendOrderCompleteToPeers<-orderComplete
					delete(orderAssignedToMap,orderComplete.Checksum)
					delete(orderMap,orderComplete.Checksum)
					delete (orderRecievedAtTime,orderComplete.Checksum)
			}





		case event:=<-DriverEvent:
			event.Checksum = event.Floor*10 + int(event.Button)
			switch event.Button{
				case driver.Internal:
					SendOrderToElevator<-event
					internalOrderMap[event.Checksum]=event
					sendNewStateToPeers(sendStateToPeers,internalOrderMap) 

				default: 
					sendOrderToPeers<-event
					unconfirmedOrderMap[event.Checksum]=event
			}


		case ack:=<-recvAckFromPeers:
			var achnowledgeIteration int
			alreadyExist := false

			achnowledgelist := achnowledgementMap[ack.Checksum]
			for _,Id:= range achnowledgelist{
				if Id == ack.Id {
					alreadyExist=true
				}
			}
			if !alreadyExist{
				achnowledgelist = append(achnowledgelist, ack.Id)
				achnowledgementMap[ack.Checksum] = achnowledgelist
			}

			for _,Id:= range achnowledgelist{
				for _,peer:= range currentPeers{
					if Id==peer{
						achnowledgeIteration++
					}
				}
			}
			if achnowledgeIteration>=len(currentPeers){
				log.Println("Recieced achnowledge from all active elevators")
				delete(unconfirmedOrderMap,ack.Checksum)
				delete(achnowledgementMap, ack.Checksum)
			}



			

		case <-ElevatorEmergency:
			log.Println("-------------------------------")
			log.Println("ElevatorEmergency is now active")
			log.Println("-------------------------------")
			peerTxEnable <- false

			//Should start some re-initalize of the elevator here.


		case <-orderResend:
			for _,v:=range unconfirmedOrderMap{
				sendOrderToPeers<-v
			}
		case <-stateResend:
			//log.Println("CurrentStateResend")
			//Println(currentElevatorState)
			sendStateToPeers<-currentElevatorState
		
		case <- orderNotCompleted:
			for checksum,timeSince := range orderRecievedAtTime{
				if time.Since(timeSince)>20*time.Second {
					if ok := OrderDelegator(stateMap,orderMap[checksum],currentPeers,orderAssignedToMap); ok{
						SendOrderToElevator<-orderMap[checksum]
					}
				}			
			}



		} 
	}
}




func sendNewStateToPeers(stateToPeers chan<- utilities.State,
	internalOrders map[int]driver.OrderEvent) {
	
	var internalOrderSlice []driver.OrderEvent

	currentElevatorState.Id,currentElevatorState.StateSentFromId = localId,localId
	
	for _,order := range internalOrders{
		internalOrderSlice = append(internalOrderSlice, order)
	}
	
	currentElevatorState.InternalOrders = internalOrderSlice
	
	stateToPeers<-currentElevatorState
}

//Might have to change some of the values
//This functions as it should. Can see if we want to rewrite it
// into several small functions instead, splitting up the cost-functions
// and combining them in a larger function, just to clean up the code.
// The cost in general is up for discussion. This algorithm does not take into account
// the time used opening and closing doors. 


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
	var currentId string
	for elevatorId, fitness := range fitnessMap{
		if fitness == minFitness{
			if elevatorId > currentId{
				currentId = elevatorId
			}
		}
		if fitness < minFitness{
			minFitness = fitness
			currentId = elevatorId
		}
	}

	orderAssignedToMap[orderEvent.Checksum]=currentId
	log.Println("Fitness")
	log.Println(fitnessMap)

	log.Println("States")
	log.Println(stateMap)
	if currentId == localId{
		return true
	}else{
		log.Println("Did not take order")
		return false
	}
}







