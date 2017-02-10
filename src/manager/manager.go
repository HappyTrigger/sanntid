package manager

import(
	".././utilities"
	"time"
	"log"

)

func Init(){

}

	var messageId int
func Run(sendMsg chan<- utilities.Message,
	recMsg <-chan  utilities.Message,
	ConnectionStatus <-chan utilities.ConnectionStatus){
/*
	time.Sleep(2*time.Second)
	msg_map := make(map[int]utilities.Message)
	msg_map[1] = utilities.Message{MessageType: utilities.MESSAGE_ORDER}
	msg_map[2] = utilities.Message{MessageType: utilities.MESSAGE_ORDER}
	msg_map[3] = utilities.Message{MessageType: utilities.MESSAGE_ORDER}
	msg_map[4] = utilities.Message{MessageType: utilities.MESSAGE_ORDER}
					
	//sendMsg<-msg2

	go func () {
		for{
			for _,v:=range msg_map{
				v.Message_Id = messageId +1
				messageId++
				sendMsg<-v
				time.Sleep(2000*time.Millisecond)
			}
		}
	}()

*/
	for{
		select{
			case msg:=<-recMsg:
				switch msg.MessageType{
					case utilities.MESSAGE_ORDER:
						log.Println("New order from", msg.Message_origin, ". Message-Id = ", msg.Message_Id)
						msg.MessageType=utilities.MESSAGE_ORDER_COMPLETE

						//sendMsg<-msg


					case utilities.MESSAGE_STATE:
						log.Println("New State")


					case utilities.MESSAGE_ORDER_COMPLETE:
						log.Println("Order complete")



					default:
						//Do nothing
				}
			case comMsg:=<-ConnectionStatus:
				log.Println("ConnectionStatus has changed")
				if comMsg.Connection != true{
					log.Println("Connection with Ip:",comMsg.Ip," has been lost")
					//Send new connectionstate to elevator for further proccesssing
				}else{
					log.Println("Connection with Ip:",comMsg.Ip," has been astablished")
					//Send new connectionstate to elevator
				}
			

		}
	}

}


