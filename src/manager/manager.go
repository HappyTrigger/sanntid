package manager

import (
	"fmt"
	"log"
	"os"
	"time"

	"../driver"
	"../network/bcast"
	"../network/localip"
	"../network/peers"
	"../utilities"
)

const (
	orderResendInterval       = 30 * time.Millisecond
	stateResendInterval       = 100 * time.Millisecond
	orderNotCompletedInterval = 2 * time.Second
)

var localId string
var currentElevatorState utilities.State

func Init(reciveStateFromPeers <-chan utilities.State,
	SendStateToPeers chan<- utilities.State,
	ElevatorStateFromElevator <-chan utilities.State,
	SendOrderToElevator chan<- driver.OrderEvent,
	InternalOrderMap map[int]driver.OrderEvent) {

	var internalOrderSlice []driver.OrderEvent
	timeout := time.After(200 * time.Millisecond)

loop:
	for {
		select {
		case state := <-reciveStateFromPeers:
			if state.Id == localId {
				for _, internalOrder := range state.InternalOrders {
					internalOrderSlice = append(internalOrderSlice, internalOrder)
					InternalOrderMap[internalOrder.Checksum] = internalOrder
				}
			}
		case state := <-ElevatorStateFromElevator:
			currentElevatorState = state
		case <-timeout:
			break loop

		}
	}
	go func() {
		for _, orders := range internalOrderSlice {
			SendOrderToElevator <- orders
			currentElevatorState.InternalOrders = append(currentElevatorState.InternalOrders, orders)
		}
	}()
	sendNewStateToPeers(SendStateToPeers, InternalOrderMap)
}

func Run(SendOrderToElevator chan<- driver.OrderEvent,
	DriverEvent <-chan driver.OrderEvent,
	ElevatorEmergency <-chan bool,
	ElevatorOrderComplete <-chan driver.OrderEvent,
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

	externalOrderMap := make(map[int]driver.OrderEvent)
	internalOrderMap := make(map[int]driver.OrderEvent)
	unconfirmedOrderMap := make(map[int]driver.OrderEvent)

	orderAssignedToMap := make(map[int]string)
	achnowledgementMap := make(map[int][]string)
	stateMap := make(map[string]utilities.State)
	orderRecievedAtTime := make(map[int]time.Time)

	orderResend := time.Tick(orderResendInterval)
	stateResend := time.Tick(stateResendInterval)
	orderNotCompleted := time.Tick(orderNotCompletedInterval)

	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	go peers.Transmitter(30201, localId, peerTxEnable)
	go peers.Receiver(30201, peerUpdateCh)

	sendOrderToPeers := make(chan driver.OrderEvent)
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

	sendAckToPeers := make(chan utilities.Achnowledgement)
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
			sendAckToPeers <- utilities.Achnowledgement{Id: localId, Checksum: msg.Checksum}
			if _, orderExist := externalOrderMap[msg.Checksum]; !orderExist {
				orderRecievedAtTime[msg.Checksum] = time.Now()
				externalOrderMap[msg.Checksum] = msg
				driver.Elev_set_button_lamp(msg.Button, msg.Floor, true)
				if ok := orderDelegator(stateMap, msg, currentPeers, orderAssignedToMap); ok {
					SendOrderToElevator <- msg
				}
			}

		case state := <-recvStateFromPeers:
			stateMap[state.Id] = state

		case p := <-peerUpdateCh:
			log.Printf("Peer update:\n")
			log.Printf("  Peers:    %q\n", p.Peers)
			log.Printf("  New:      %q\n", p.New)
			log.Printf("  Lost:     %q\n", p.Lost)

			currentPeers = p.Peers

			if state, ok := stateMap[p.New]; ok {
				state.StateSentFromId = localId
				log.Println("Sending state to reconnecting elevator")
				for i := 0; i < 10; i++ {
					sendStateToPeers <- state

				}
			}

			for _, lostId := range p.Lost {
				for checksum, Id := range orderAssignedToMap {
					if Id == lostId {
						msg := externalOrderMap[checksum]
						orderRecievedAtTime[msg.Checksum] = time.Now()
						if ok := orderDelegator(stateMap, msg, currentPeers, orderAssignedToMap); ok {
							SendOrderToElevator <- msg

						}
					}
				}
			}

		case state := <-ElevatorStateFromElevator:
			stateMap[localId] = state
			currentElevatorState = state
			sendNewStateToPeers(sendStateToPeers, internalOrderMap)

		case orderComplete := <-recOrderCompleteFromPeers:
			log.Println("Order at Floor:", orderComplete.Floor+1, " completed")
			delete(orderAssignedToMap, orderComplete.Checksum)
			delete(externalOrderMap, orderComplete.Checksum)
			delete(orderRecievedAtTime, orderComplete.Checksum)
			driver.Elev_set_button_lamp(orderComplete.Button, orderComplete.Floor, false)

		case orderComplete := <-ElevatorOrderComplete:
			switch orderComplete.Button {
			case driver.Internal:
				log.Println("Internal Order completed")
				delete(internalOrderMap, orderComplete.Checksum)
				sendNewStateToPeers(sendStateToPeers, internalOrderMap)

			default:
				sendOrderCompleteToPeers <- orderComplete
				delete(orderAssignedToMap, orderComplete.Checksum)
				delete(externalOrderMap, orderComplete.Checksum)
				delete(orderRecievedAtTime, orderComplete.Checksum)
			}

		case event := <-DriverEvent:
			event.Checksum = event.Floor*10 + int(event.Button)
			switch event.Button {
			case driver.Internal:
				SendOrderToElevator <- event
				internalOrderMap[event.Checksum] = event
				sendNewStateToPeers(sendStateToPeers, internalOrderMap)

			default:
				sendOrderToPeers <- event
				unconfirmedOrderMap[event.Checksum] = event
			}

		case ack := <-recvAckFromPeers:
			var achnowledgeIteration int
			alreadyExist := false

			achnowledgelist := achnowledgementMap[ack.Checksum]
			for _, Id := range achnowledgelist {
				if Id == ack.Id {
					alreadyExist = true
				}
			}
			if !alreadyExist {
				achnowledgelist = append(achnowledgelist, ack.Id)
				achnowledgementMap[ack.Checksum] = achnowledgelist
			}

			for _, Id := range achnowledgelist {
				for _, peer := range currentPeers {
					if Id == peer {
						achnowledgeIteration++
					}
				}
			}
			if achnowledgeIteration >= len(currentPeers) {
				delete(unconfirmedOrderMap, ack.Checksum)
				delete(achnowledgementMap, ack.Checksum)
			}

		case <-ElevatorEmergency:
			log.Println("-------------------------------")
			log.Println("ElevatorEmergency is now active")
			log.Println("-------------------------------")
			peerTxEnable <- false

			driver.Elev_set_motor_direction(driver.MotorStop)

			err := utilities.Reboot(localId)

			if err != nil {
				panic("Major malfunction, call technical assistance")
			}
			panic("Something went wrong, restarting")

		case <-orderResend:
			for _, v := range unconfirmedOrderMap {
				sendOrderToPeers <- v
			}
		case <-stateResend:
			sendStateToPeers <- currentElevatorState

		case <-orderNotCompleted:
			for checksum, timeSince := range orderRecievedAtTime {
				if time.Since(timeSince) > 20*time.Second {
					if orderAssignedToMap[checksum] == localId && len(currentPeers) > 1 {
						err := utilities.Reboot(localId)
						if err != nil {
							panic("Major malfunction, call technical assistance")
						}
						panic("The elevator is not working properly, rebooting and redistributing orders")
					}
					if ok := orderDelegator(stateMap, externalOrderMap[checksum], currentPeers, orderAssignedToMap); ok {
						SendOrderToElevator <- externalOrderMap[checksum]
						orderRecievedAtTime[checksum] = time.Now()
					}
				}
			}
		}
	}
}

func sendNewStateToPeers(stateToPeers chan<- utilities.State,
	internalOrders map[int]driver.OrderEvent) {

	var internalOrderSlice []driver.OrderEvent
	currentElevatorState.Id, currentElevatorState.StateSentFromId = localId, localId

	for _, order := range internalOrders {
		internalOrderSlice = append(internalOrderSlice, order)
	}

	currentElevatorState.InternalOrders = internalOrderSlice
	stateToPeers <- currentElevatorState
}
