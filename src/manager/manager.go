package manager

import(
	".././utilities"
	"time"
	"fmt"

)

func Init(){

}


func Run(recMsg <-chan utilities.Message, sendMsg chan <- utilities.Message ){

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
						fmt.Println("New order")


					case utilities.MESSAGE_STATE:
						fmt.Println("New State")


					case utilities.MESSAGE_ORDER_COMPLETE:
						fmt.Println("Order complete")


				}

		}
	}

}


