package manager

import(
	".././utilities"
	"time"
	"fmt"

)

func Init(){

}


func Run(recMsg <-chan utilities.Message, sendMsg chan <- utilities.Message ){

	msg2 := utilities.Message{MessageType: utilities.MESSAGE_ACKNOLEDGE}
	for{
		sendMsg<-msg2
		time.Sleep(3*time.Second)
		msg2:=<-recMsg
		fmt.Printf("%+v\n", msg2.MessageType)
		msg2.MessageType+=1



		
	}

}