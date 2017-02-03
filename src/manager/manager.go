package manager

import(
	".././utilities"
	"time"
	"log"

)

func Init(){

}


func Run(sendMsg chan<- utilities.Message, recMsg <-chan  utilities.Message ){

	msg2 := utilities.Message{MessageType: utilities.MESSAGE_ORDER}
	
	go func () {
		sendMsg<-msg2
		time.Sleep(3*time.Second)
	}()


	for{
		select{
			case msg:=<-recMsg:
				switch msg.MessageType{
					case utilities.MESSAGE_ORDER:
						log.Println("New order")


					case utilities.MESSAGE_STATE:
						log.Println("New State")


					case utilities.MESSAGE_ORDER_COMPLETE:
						log.Println("Order complete")


				}

		}
	}

}


