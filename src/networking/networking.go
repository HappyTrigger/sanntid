package networking


import (

	"./udp"
    //"fmt"
    //"time"
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
	udBroadcastHeartBeat := make(chan []byte)
	//udpBroadCastAchnowledge := make(chan []byte)

	

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
	go SendHeartBeat(udBroadcastHeartBeat)
	
/*
	go send_udp_message(udpBroadcastMsg,
		fromManager,
		achnowledge,
		toManager,
		connection_status,
		connectionLost,
		udBroadcastHeartBeat,
		udpBroadCastAchnowledge)
*/	
	heartbeatChan:=Heartbeat_recieved(udpBroadcastMsg,
		connection_status,
		connectionLost)
	


	for{		
		select {
			case msg:=<-fromManager:
				encodedMsg:=utilities.Encoder(msg)
				udpBroadcastMsg<-encodedMsg


			case <-achnowledge:
				log.Println("achnowledge")
    		

    		case heartbeat:=<-udBroadcastHeartBeat:
    			//log.Println("Heartbeat sent")
    			udpBroadcastMsg<-heartbeat

			

			case raw_m := <-udpRecvMsg:
				msg:=utilities.Decoder(raw_m.Data)
				msg.Message_sender=raw_m.Ip


				switch msg.MessageType{
					case utilities.MESSAGE_ACKNOWLEDGE: 

						if msg.Message_origin == localIp{
							achnowledge<-msg
						}
						

					case utilities.MESSAGE_HEARTBEAT: 
						heartbeatChan<-msg

					//case utilities.MESSAGE_ORDER_COMPLETE:
						//toManager<-msg


					default:
						toManager<-msg //Sends the message to the manager
					
						
						msg.MessageType = utilities.MESSAGE_ACKNOWLEDGE
						encodedMsg:=utilities.Encoder(msg)
						udpBroadcastMsg<-encodedMsg

				}
		}
	}

}



/*

func send_udp_message(udpBroadCast chan<-[]byte,
	fromManager <-chan utilities.Message,
	achnowledge_chan <-chan utilities.Message,
	sendToManager chan<- utilities.Message,
	connectionStatusChan chan<- utilities.ConnectionStatus,
	connectionLost chan<- utilities.ConnectionStatus,
	udBroadcastHeartBeat<-chan []byte,
	udpBroadCastAchnowledge<-chan []byte){
	


	achnowledgement_confirmed := make(map[string]bool)
	
	for{
		select{
			case msg:=<-fromManager:
				msg.Message_origin = localIp
				encoded_msg:=utilities.Encoder(msg)
				for i:=0;i<5;i++{

					udpBroadCast<-encoded_msg
					forloop:
					for{
						select{
						case ach:=<-achnowledge_chan:
							if (msg.Message_Id == ach.Message_Id && 
								achnowledgement_confirmed[ach.Message_sender] == false){
								
								achnowledgement_confirmed[ach.Message_sender] = true
								break forloop
							}
						case <-time.After(40*time.Millisecond):
							break forloop
						}
					}
				}
				for k, v := range achnowledgement_confirmed { 
    				if v != true{
    					//K is now inactive/ not responding
    					log.Println("Transfer of files failed, new connection_status")
    					connectionLost<-utilities.ConnectionStatus{Ip:k, Connection:false}
    					delete(achnowledgement_confirmed,k)
    				}else{
    					achnowledgement_confirmed[k]=false
    					sendToManager<-msg
    					log.Println("Achnowledgement confirmed")
    				}
    			}
    			
 
				

    		case <-achnowledge_chan:
    			//acknowledge-dump



    		case heartbeat:=<-udBroadcastHeartBeat:
    			udpBroadCast<-heartbeat

   


		}

	}
}


*/


