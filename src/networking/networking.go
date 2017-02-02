package networking


import (

	"./udp"
    //"fmt"
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
				buf:=utilities.Encoder(msg)
				udpBroadcastMsg<-buf
				


	


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
						log.Println("Message achnoledge")


					case utilities.MESSAGE_HEARTBEAT: 
						log.Println("Heartbeat recieved")
						Heartbeat_recieved(msg)
					default:
						Send_msg<-msg

			}
		}
	}
}

func SendHeartBeat(counter int) {
	
}

func Heartbeat_recieved(m utilities.Message) {
	
}


