package manager

import(
	".././utilities"
	"time"
	"log"

)

func Init(){

}

	var messageId int
func Run(sendMsg chan<- utilities.Message, recMsg <-chan  utilities.Message, ConnectionStatus chan utilities.ConnectionStatus ){


	msg2 := utilities.Message{MessageType: utilities.MESSAGE_ORDER}
	
	go func () {
		for{
		time.Sleep(5*time.Second)
		msg2.Message_Id = messageId+1
		messageId+=1
		sendMsg<-msg2
		
		log.Println("Sending Message")
		}
	}()


	for{
		select{
			case msg:=<-recMsg:
				switch msg.MessageType{
					case utilities.MESSAGE_ORDER:
						log.Println("New order, number ", msg.Message_Id)



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


