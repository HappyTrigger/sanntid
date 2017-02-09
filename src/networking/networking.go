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

func Run(sendMsg <-chan utilities.Message,
	recMsg chan<- utilities.Message,
	connection_status chan<-utilities.ConnectionStatus){
	

	log.Println("---Starting network loop---")
	log.Println("The ip of this computer is: ", localIp)


	//status:=true

	//Channels
	//udpBroadcastMsg,udpRecvMsg:=udp.Init(localIp)
	recivedMsg := make(chan utilities.Message)
	achnowledge := make(chan utilities.Message)
	heartbeatmap :=make(map[string]int)
	

	//Testing system ////////////////
	udpBroadcastMsg,udpRecvMsg := make(chan []byte), make(chan udp.RawMessage)
	go func(){
		for{
			select{
			case msg:=<-udpBroadcastMsg:
				 udpRecvMsg<-udp.RawMessage{Data:msg,Ip:localIp}
			}
		}
	}()
	////////////////////
	go SendHeartBeat(udpBroadcastMsg)
	heartbeatChan:=Heartbeat_recieved(udpBroadcastMsg,connection_status,heartbeatmap)

	go send_udp_message(udpBroadcastMsg,sendMsg,achnowledge,recMsg,connection_status)
	go handel_UDP_message(recivedMsg, recMsg, achnowledge, udpBroadcastMsg,connection_status,heartbeatmap,heartbeatChan)


	for{		
		select {
			case raw_m := <-udpRecvMsg:
				msg:=utilities.Decoder(raw_m.Data)
				msg.Message_sender=raw_m.Ip
				//log.Println("UDPRecvMsg")
				recivedMsg<-msg
		}
	}

}





func send_udp_message(udpBroadCast chan<-[]byte,
	sendMsg <-chan utilities.Message,
	achnowledge_chan <-chan utilities.Message,
	sendToManager chan<- utilities.Message,
	connectionStatusChan chan<- utilities.ConnectionStatus){
	


	achnowledgement_confirmed := make(map[string]bool)
	
	for{
		select{
			case msg:=<-sendMsg:
				//log.Println("Sending message")

				msg.Message_origin = localIp
//				log.Println("Local ip assigned")
				encoded_msg:=utilities.Encoder(msg)
//				log.Println("Encoded message")
				for i:=0;i<2;i++{ 
					udpBroadCast<-encoded_msg
//					log.Println("broadcast encoded message")
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
//							log.Println("Time after triggerd")
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
    					connectionStatusChan<-utilities.ConnectionStatus{Ip:k, Connection:false}
    				}else{
    					achnowledgement_confirmed[k]=false
    				}
    			}
    			
    			//log.Println("Achnowledgment done")
  //  			log.Println("Sendtomanager")
				sendToManager<-msg


    			case <-achnowledge_chan:
    			//Dump for trivial achnowledgements
    				log.Println("Achnowledgement came after timeout")
    			

    			default:
    				//Do nothing


		}

	}
}


func handel_UDP_message(recivedMsg <-chan utilities.Message,
	sendToManager chan<-utilities.Message,
	achnowledge_chan chan<- utilities.Message,
	udpBroadCast chan <-[]byte,
	connectionStatus chan<- utilities.ConnectionStatus,
	heartbeat_map map[string]int,
	heartbeat_chan chan utilities.Message){
	

	for{
		select{
			case msg:=<-recivedMsg:
			//	log.Println("Recieved message")
				switch msg.MessageType{

					case utilities.MESSAGE_ACKNOWLEDGE: 
						//log.Println("Achnowledgement for message :",msg.Message_Id)
						if msg.Message_origin == localIp{
							//log.Println("Is it locked here")
							achnowledge_chan<-msg
							//log.Println("Achnowledgement recived from :",msg.Message_sender)
						}else{
							log.Println("Achnowledgement from another elevator")
						}
						

					case utilities.MESSAGE_HEARTBEAT: 
						log.Println("Heartbeat recieved")
						heartbeat_chan<-msg


					default:
						//sendToManager<-msg //Sends the message to the manager
						
						//Task Send achnolwedge back to sender
						msg.MessageType = utilities.MESSAGE_ACKNOWLEDGE
				 		log.Println("udpBroadcast achnowledgement")
						udpBroadCast<-utilities.Encoder(msg)

			}
		}
	}
}



