package networking

import(
	"log"
	".././utilities"
	"time"
)


func SendHeartBeat(udpBroadcastMsg chan<-[]byte){
	
	udpHeartBeatNum := 0
	udpHeartBeatTick := time.Tick(1*time.Second)

	for{
		select{
			case<-udpHeartBeatTick:
					
				data:=utilities.Message{MessageType: utilities.MESSAGE_HEARTBEAT,
					Heartbeat: utilities.CreateHeartbeat(udpHeartBeatNum)}
					
				msg:=utilities.Encoder(data)

				udpBroadcastMsg<-msg

				udpHeartBeatNum++
				//case<-stop:
				//	return

			}
		}
	

	return
	
}

func Heartbeat_recieved(udpBroadcastMsg chan<-[]byte,
	connectionStatus chan<- utilities.ConnectionStatus,
	heartbeat_map map[string]int) chan utilities.Message {
	
	newHeartbeat := make(chan utilities.Message)


	go func() {
		heartbeatTimer := time.Tick(1*time.Second)
		failed_heartbeats := make(map[string]int)
		connection_map :=  make(map[string]string)


		
		var val string

		for{

			select{
			case heartbeat:=<-newHeartbeat:

				if _, ok := connection_map[heartbeat.Message_sender]; !ok {
    				connection_map[heartbeat.Message_sender]=heartbeat.Message_sender
    				connectionStatus<-utilities.ConnectionStatus{Ip:heartbeat.Message_sender,Connection: true}
    				
    					log.Println("New elevator found",ok)
				}

				prev:= heartbeat_map[heartbeat.Message_sender]
				current:=heartbeat.Heartbeat.Counter
				if prev+1 == current{
					heartbeat_map[heartbeat.Message_sender]=current
					log.Println("Recived heartbeat matches")
					failed_heartbeats[heartbeat.Message_sender]=0
				}else{
				 	
					heartbeat_map[heartbeat.Message_sender]=current
					log.Println("Recived heartbeat is lagging")
				}


			case <-heartbeatTimer: 
				for _, val := range connection_map { 
    				failed_heartbeats[val]=failed_heartbeats[val]+1
					}

					if failed_heartbeats[val]>3{
						connectionStatus<-utilities.ConnectionStatus{Ip:val,Connection: false}
						delete(connection_map,val)
					}





				
				




			}
		}
	}()


	return newHeartbeat



}
