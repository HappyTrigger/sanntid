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






func init() {

	var err error

	localIp, err = getLocalIp()

	if err != nil {
		log.Fatal("Refusing to start - ",err)
	}
}

func Run(fromManager <-chan utilities.Message,
	toManager chan<- utilities.Message,
	connection_status chan<- utilities.ConnectionStatus){
	

	log.Println("---Starting network loop---")
	log.Println("The ip of this computer is: ", localIp)


	//Channels
	udpBroadcastMsg,udpRecvMsg:=udp.Init(localIp)

	achnowledge := make(chan utilities.Message)
	connectionLost := make(chan utilities.ConnectionStatus)

	

	//Testing system ////////////////
//	udpBroadcastMsg,udpRecvMsg := make(chan []byte), make(chan udp.RawMessage)
//	go func(){
//		for{
//			select{
//			case msg:=<-udpBroadcastMsg:
//				 udpRecvMsg<-udp.RawMessage{Data:msg,Ip:localIp}
//			}
//		}
//	}()
	////////////////////
	go SendHeartBeat(udpBroadcastMsg)
	go send_udp_message(udpBroadcastMsg,fromManager,achnowledge,toManager,connection_status,connectionLost)
	heartbeatChan:=Heartbeat_recieved(udpBroadcastMsg,connection_status,connectionLost)
	


	for{		
		select {
			case raw_m := <-udpRecvMsg:
				msg:=utilities.Decoder(raw_m.Data)
				msg.Message_sender=raw_m.Ip


				switch msg.MessageType{
					case utilities.MESSAGE_ACKNOWLEDGE: 

						if msg.Message_origin == localIp{

							achnowledge<-msg
							
						}else{
							log.Println("Achnowledgement from another elevator")
						}
						

					case utilities.MESSAGE_HEARTBEAT: 
						//log.Println("Heartbeat recieved")
						heartbeatChan<-msg


					default:
						toManager<-msg //Sends the message to the manager
						
						//Task Send achnolwedge back to sender
						msg.MessageType = utilities.MESSAGE_ACKNOWLEDGE
				 		//log.Println("udpBroadcast achnowledgement")
						udpBroadcastMsg<-utilities.Encoder(msg)
				}
		}
	}

}





func send_udp_message(udpBroadCast chan<-[]byte,
	fromManager <-chan utilities.Message,
	achnowledge_chan <-chan utilities.Message,
	sendToManager chan<- utilities.Message,
	connectionStatusChan chan<- utilities.ConnectionStatus,
	connectionLost chan<- utilities.ConnectionStatus ){
	


	achnowledgement_confirmed := make(map[string]bool)
	
	for{
		select{
			case msg:=<-fromManager:
				msg.Message_origin = localIp
				encoded_msg:=utilities.Encoder(msg)
				for i:=0;i<2;i++{ 
					udpBroadCast<-encoded_msg
					forloop:
					for{
						select{
						case ach:=<-achnowledge_chan:
//							log.Println("Got achnowledgement")
							//log.Println("Achnowledgement recived")
							if msg.Message_Id == ach.Message_Id{
//								log.Println("got achnowledgement from right ip")
								achnowledgement_confirmed[ach.Message_sender] = true
								//log.Println("Achnowledgement for message :",ach.Message_Id, "Origing message_id: ",msg.Message_Id)
								}
							break forloop
						case <-time.After(20*time.Millisecond):
							break forloop
						}
					}
				}
				//log.Println("Checking if all recived")
//				log.Println("Map:",achnowledgement_confirmed)
				for k, v := range achnowledgement_confirmed { 
    				if v != true{
    					//K is now inactive/ not responding
    					//log.Println("Transfer of files failed, new connection_status")
    					connectionLost<-utilities.ConnectionStatus{Ip:k, Connection:false}
    				}else{
    					achnowledgement_confirmed[k]=false
    				}
    			}
    			
    			//log.Println("Achnowledgment done")
  //  			log.Println("Sendtomanager")
				sendToManager<-msg

    		case <-achnowledge_chan:
    			log.Println("Achnowledgement came after timeout")
    			

    		default:
    				//Do nothing


		}

	}
}





