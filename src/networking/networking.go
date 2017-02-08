package networking


import (

	"./udp"
    //"fmt"
    "time"
	//"strconv"
	"log"
	".././utilities"
	//"sync"

)

var localIp string
var msg_Id int


func init() {
	msg_Id = 0
	var err error

	localIp, err = getLocalIp()

	if err != nil {
		log.Fatal("Refusing to start - ",err)
	}
}

func Run(sendMsg <-chan utilities.Message,recMsg chan<- utilities.Message, connection_status chan utilities.ConnectionStatus){


	log.Println("---Starting network loop---")
	log.Println("The ip of this computer is: ", localIp)



	//Channels
	udpBroadcastMsg,udpRecvMsg:=udp.Init(localIp)
	processChan := make(chan utilities.Message)
	achnowledge := make(chan utilities.Message)
	


//	udpBroadcastMsg,udpRecvMsg := make(chan []byte,50), make(chan udp.RawMessage,50)
//	go func(){
//		for{
//			select{
//			case msg:=<-udpBroadcastMsg: 
//				 udpRecvMsg<-udp.RawMessage{Data:msg,Ip:localIp}
//			}
//		}
//	}()

	go send_udp_message(udpBroadcastMsg,sendMsg,achnowledge,recMsg,connection_status)
	go handel_UDP_message(processChan, recMsg, achnowledge, udpBroadcastMsg)




	for{		
		select {
			case raw_m := <-udpRecvMsg:
				msg:=utilities.Decoder(raw_m.Data)
				msg.Message_sender=raw_m.Ip
				processChan<-msg
		}
	}

}

func send_udp_message(udpBroadCast chan<-[]byte,
	sendMsg <-chan utilities.Message,
	achnowledge_chan <-chan utilities.Message,
	sendToManager chan<- utilities.Message,
	connectionStatusChan chan<- utilities.ConnectionStatus){
	


	achnowledge_map := make(map[string]bool)
 


	for{
		select{
			case msg:=<-sendMsg:


				encoded_msg:=utilities.Encoder(msg)
				for i:=0;i<2;i++{
					udpBroadCast<-encoded_msg
					time.Sleep(5*time.Millisecond)
					forloop:
					for{
						select{
						case ach:=<-achnowledge_chan:
							if msg.Message_Id == ach.Message_Id{
								achnowledge_map[ach.Message_sender] = true
								}
							break forloop
						case <-time.After(10*time.Millisecond):
							break forloop
						}
					}
				}

				for k, v := range achnowledge_map { 
    				if v != true{
    					//K is now inactive/ not responding
    					connectionStatusChan<-utilities.ConnectionStatus{Ip:k, Connection:false}
    				}else{
    					v = false
    				}
    			}


				//sendToManager<-msg

		}

	}
}


func handel_UDP_message(recivedMsg <-chan utilities.Message,
	sendToManager chan<-utilities.Message,
	achnowledge_chan chan<- utilities.Message,
	udpBroadCast chan <-[]byte){
	



	for{
		select{
			case msg:=<-recivedMsg:
				switch msg.MessageType{

					case utilities.MESSAGE_ACKNOWLEDGE: 
						log.Println("Achnowledgement for message :",msg.Message_Id)
						achnowledge_chan<-msg

					case utilities.MESSAGE_HEARTBEAT: 
						log.Println("Heartbeat recieved")
						Heartbeat_recieved(msg)

					default:
						sendToManager<-msg //Sends the message to the manager
						
						//Task Send achnolwedge back to sender
						msg.MessageType = utilities.MESSAGE_ACKNOWLEDGE
						udpBroadCast<-utilities.Encoder(msg)

			}
		}
	}
}



