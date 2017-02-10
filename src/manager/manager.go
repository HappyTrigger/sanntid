package manager

import (
	//".././elevator"
	".././utilities"
	"log"
	"time"
)

var AddOrder utilities.NewOrder
var messageId int

func Init(ExternalOrdersMap map[utilities.NewOrder]int) {
	//Initializes the Map of external orders
}

//network channels
func Run(sendMsg chan<- utilities.Message,
	recMsg <-chan utilities.Message,
	ConnectionStatus <-chan utilities.ConnectionStatus,
	//Elevator
	NewState <-chan utilities.State,
	ExtOrderRaised <-chan utilities.NewOrder,
	TakesExtOrd chan<- utilities.NewOrder,
	//
	ExternalOrdersMap map[utilities.NewOrder]int,
	state_map map[string]utilities.State) {

	ExternalOrders := make(map[utilities.NewOrder]int)

	time.Sleep(2 * time.Second)
	msg_map := make(map[int]utilities.Message)
	msg_map[1] = utilities.Message{MessageType: utilities.MESSAGE_ORDER}
	msg_map[2] = utilities.Message{MessageType: utilities.MESSAGE_ORDER}
	msg_map[3] = utilities.Message{MessageType: utilities.MESSAGE_ORDER}
	msg_map[4] = utilities.Message{MessageType: utilities.MESSAGE_ORDER}

	//sendMsg<-msg2

	go func() {
		for {
			for _, v := range msg_map {
				v.Message_Id = messageId + 1
				messageId++
				sendMsg <- v
				log.Println("Sent Order nr", v.Message_Id)
				time.Sleep(2000 * time.Millisecond)
			}
		}
	}()

	//Manager to Elevator
	go DistributeOrder(TakesExtOrd)

	for {
		select {
		case msg := <-recMsg:
			switch msg.MessageType {
			case utilities.MESSAGE_ORDER:
				//log.Println("New order from", msg.Message_origin, ". Message-Id = ", msg.Message_Id)
				msg.MessageType = utilities.MESSAGE_ORDER_COMPLETE

				//sendMsg<-msg

			case utilities.MESSAGE_STATE:
				log.Println("New State")

			case utilities.MESSAGE_ORDER_COMPLETE:
				log.Println("Order complete")

			default:
				//Do nothing
			}
		case comMsg := <-ConnectionStatus:
			log.Println("ConnectionStatus has changed")
			if comMsg.Connection != true {
				log.Println("Connection with Ip:", comMsg.Ip, " has been lost")
				//Send new connectionstate to elevator for further proccesssing
			} else {
				log.Println("Connection with Ip:", comMsg.Ip, " has been astablished")
				//Send new connectionstate to elevator
			}
		default:
			//

		// Elevator to Manager
		case MyElevator := <-NewState:

		case MyExtOrd := <-ExtOrderRaised:

		}
	}

}

func Cost() {

}

func DistributeOrder(TakesExtOrd chan<- utilities.NewOrder) {
	if false {
		TakesExtOrd <- AddOrder
	}
}
