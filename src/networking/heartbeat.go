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


			}
		}
	

	return
	
}

func Heartbeat_recieved(udpBroadcastMsg chan<-[]byte,
	connectionStatus chan<- utilities.ConnectionStatus,
	connectionLost <-chan utilities.ConnectionStatus) chan utilities.Message {
	
	newHeartbeat := make(chan utilities.Message)


	go func() {
		heartbeatTimer := time.Tick(2*time.Second)
		failed_heartbeats := make(map[string]int)
		connection_map :=  make(map[string]string)
		heartbeat_map := make(map[string]int)

		for{

			select{
			case heartbeat:=<-newHeartbeat:

				if _, ok := connection_map[heartbeat.Message_sender]; !ok {
    				connection_map[heartbeat.Message_sender]=heartbeat.Message_sender
    				heartbeat_map[heartbeat.Message_sender] = heartbeat.Heartbeat.Counter
    				connectionStatus<-utilities.ConnectionStatus{Ip:heartbeat.Message_sender,Connection: true}
    				
    					log.Println("New elevator found")
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
    				log.Println("Ip: ",val," Failed heartbeats: ",failed_heartbeats[val])
					

					if failed_heartbeats[val]>3{
						connectionStatus<-utilities.ConnectionStatus{Ip:val,Connection: false}
						delete(connection_map,val)
						failed_heartbeats[val]=0

					}
				}
			case conMsg:=<-connectionLost:
				connectionStatus<-conMsg
				delete(connection_map,conMsg.Ip)

			default:
				//
			}
		}
	}()


	return newHeartbeat



}
