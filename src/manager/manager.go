package manager

import (
	".././utilities"
	"log"
	//"math"
	"time"
	".././mydriver"
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
	DriverEvent <-chan driver.OrderEvent,
	SendOrderToElevator chan<- driver.OrderEvent,
	DoorOpen <- chan bool,
	DoorClosed <-chan bool,
	ElevatorEmergency <-chan bool) {


	time.Sleep(2 * time.Second)


	LocalIp := networking.GetLocalIp()
	StateMap := make(map[string]utilities.State)
	ConnectionMap := make(map[string]bool)



	for {
		select {
		case msg := <-reciveFromNetwork:
			switch msg.MessageType {
			case utilities.MESSAGE_ORDER:

				order:= driver.OrderEvent{Floor:msg.NewOrder.Floor,
					Button:msg.NewOrder.Button,
					OrderId:msg.Message_Id}

				log.Println("Recived order from network")

				SendOrderToElevator<-order


			case utilities.MESSAGE_STATE:
				StateMap[msg.Message_sender]=msg.State

			case utilities.MESSAGE_ORDER_COMPLETE:
				log.Println("Order complete")

			default:
				//Do nothing
			}










		case comMsg := <-ConnectionStatus:
			log.Println("ConnectionStatus has changed")
			if comMsg.Connection != true {
				log.Println("Connection with Ip:", comMsg.Ip, " has been lost")
				ConnectionMap[comMsg.Message_sender]=false
				
			} else {
				log.Println("Connection with Ip:", comMsg.Ip, " has been astablished")
				ConnectionMap[comMsg.Message_sender]=true

			}


		case event:=<-DriverEvent:


			newOrder:= utilities.NewOrder{Floor: event.Floor,Button: event.Button }
			sendMsg:= utilities.Message{NewOrder:newOrder, MessageType: utilities.MESSAGE_ORDER }
			sendToNetwork<-sendMsg
			log.Println("Sending orders")
		}
	}
}

