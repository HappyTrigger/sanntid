package networking


import (

	"./udp"
    "fmt"
    //"time"
	//"strconv"
	"log"
	".././utilities"

)

var localIp string


func init() {

	var err error

	localIp, err = getLocalIp()

	if err != nil {
		log.Fatal("Refusing to start - ",err)
	}
}

func Run(sendMsg <-chan utilities.Message,recMsg chan<- utilities.Message){


	log.Println("---Starting network loop---")
	log.Println("The ip of this computer is: ", localIp)	

	recievedMessage := utilities.Message{}
	udpBroadcastMsg,udpRecvMsg:=udp.Init(localIp)
	processChan := make(chan utilities.Message)
	
	

	go processUDPMsg(processChan, recMsg)


	for{		
		select {
			case msg := <-sendMsg: 
				buf,err:=utilities.Encoder(msg)
				if err==nil{
					udpBroadcastMsg<-buf
				}else {
        			fmt.Println("Error while decoding")
    			}


	


			case raw_m := <-udpRecvMsg:
				utilities.Decoder(raw_m.Data,&recievedMessage)	
				processChan<-recievedMessage
		}
	}

}


func processUDPMsg(processChan <-chan utilities.Message, Send_msg chan<-utilities.Message ) {
	for{
		select{
			case msg:=<-processChan:
				switch msg.MessageType{

					case utilities.MESSAGE_ACKNOLEDGE: 



					case utilities.MESSAGE_HEARTBEAT: 

					default:
						Send_msg<-msg

			}
		}
	}
}



