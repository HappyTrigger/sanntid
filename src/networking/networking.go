package networking


import (

	"./udp"
    //"fmt"
    "time"
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

	udpBroadcastMsg,udpRecvMsg:=udp.Init(localIp)
	processChan := make(chan utilities.Message)
	msgSentMap := make(map[int][]byte)
	var MsgAchnowledgeNumber int
	

	MsgAchnowledgeNumber = 0

	go processUDPmsg(processChan, recMsg,msgSentMap)
	go resendMsg(msgSentMap,udpBroadcastMsg)


	for{		
		select {
			case msg := <-sendMsg: 
				if msg.Message_Id == 0{
					msg.Message_Id=MsgAchnowledgeNumber+1
				}
				buf:=utilities.Encoder(msg)
				udpBroadcastMsg<-buf
				
				if (msg.MessageType != utilities.MESSAGE_ACKNOWLEDGE) && 
				(msg.MessageType!=utilities.MESSAGE_HEARTBEAT){
				
				msgSentMap[msg.Message_Id]=buf
				}


	


			case raw_m := <-udpRecvMsg:
				recievedMessage:=utilities.Decoder(raw_m.Data)
				recievedMessage.Message_sender=raw_m.Ip
				
				processChan<-recievedMessage

				//This needs to be rewritten, just want to test the consept out

				if (recievedMessage.MessageType != utilities.MESSAGE_ACKNOWLEDGE) && 
				(recievedMessage.MessageType!=utilities.MESSAGE_HEARTBEAT){
					ach :=utilities.Acknowledge{Message_recieved_from: recievedMessage.Message_sender,Message_Id:recievedMessage.Message_Id}
					achmsg:= utilities.Message{Acknowledge: ach}
					buf:=utilities.Encoder(achmsg)
					udpBroadcastMsg<-buf
				}

		}
	}

}


func processUDPmsg(processChan <-chan utilities.Message, rec_msg chan<-utilities.Message,  msgMap map[int][]byte ) {
	for{
		select{
			case msg:=<-processChan:
				switch msg.MessageType{

					case utilities.MESSAGE_ACKNOWLEDGE: 
						log.Println("Message achnoledge from", msg.Message_sender)
						delete(msgMap,msg.Message_Id) //Achnowledgement recieved, No need to resend message anymore

					case utilities.MESSAGE_HEARTBEAT: 
						log.Println("Heartbeat recieved")
						Heartbeat_recieved(msg)
					default:
						rec_msg<-msg //Sends the message to the manager



			}
		}
	}
}


func resendMsg(msgMap map[int][]byte, msg chan<-[]byte ) {
	for{
		time.Sleep(1*time.Second)
	for _, v := range msgMap{
		msg<-v
		}
	}
}



