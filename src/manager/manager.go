package manager

import (
	//"./networking"
	//".././elevator"
	".././utilities"
	"log"
	//"math"
	"time"
)

var AddOrder utilities.NewOrder
var messageId int

func Init(ExternalOrdersMap map[utilities.NewOrder]int) {
	//Initializes the Map of external orders
}

//network channels
func Run(sendToNetwork chan<- utilities.Message,
	reciveFromNetwork <-chan utilities.Message,
	ConnectionStatus <-chan utilities.ConnectionStatus,
	NewState <-chan utilities.State,
	DriverEvent <-chan utilities.NewOrder,
	SendOrderToElevator chan<- utilities.NewOrder) {


	time.Sleep(2 * time.Second)
	msg_map := make(map[int]utilities.Message)
	msg_map[1] = utilities.Message{MessageType: utilities.MESSAGE_ORDER}
	msg_map[2] = utilities.Message{MessageType: utilities.MESSAGE_ORDER}
	msg_map[3] = utilities.Message{MessageType: utilities.MESSAGE_ORDER}
	msg_map[4] = utilities.Message{MessageType: utilities.MESSAGE_ORDER}

	//sendToNetwork<-msg2

	go func() {
		for {
			for _, v := range msg_map {
				v.Message_Id = messageId + 1
				messageId++
				sendToNetwork <- v
				log.Println("Sent Order nr", v.Message_Id)
				time.Sleep(2000 * time.Millisecond)
			}
		}
	}()


	for {
		select {
		case msg := <-reciveFromNetwork:
			switch msg.MessageType {
			case utilities.MESSAGE_ORDER:
				//log.Println("New order from", msg.Message_origin, ". Message-Id = ", msg.Message_Id)
				msg.MessageType = utilities.MESSAGE_ORDER_COMPLETE

				//sendToNetwork<-msg

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


		//---------------------------------Elevator to Manager----------------------------------
		//case MyElevator := <-NewState:
			//state_map map[string]utilities.State
		//	state_map[networking.localIp] = NewState
		//case MyExtOrd := <-ExtOrderRaised:
			//I give a key to the order
			//ExternalOrders[ExtOrderRaised] = 10*msg.Order.Floor + math.Abs(msg.Order.Direction)
		}
	}

}
/*
func Cost() {

}

func DistributeOrder(TakesExtOrd chan<- utilities.NewOrder) {
	if false {
		TakesExtOrd <- AddOrder
	}
}
*/