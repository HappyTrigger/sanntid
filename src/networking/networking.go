package networking


import (

	"./udp"
    //"fmt"
    "time"
	//"strconv"
	"log"
	".././utilities"
	"sync"

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

func Run(sendMsg <-chan utilities.Message,recMsg chan<- utilities.Message){


	log.Println("---Starting network loop---")
	log.Println("The ip of this computer is: ", localIp)


	//Channels
	//udpBroadcastMsg,udpRecvMsg:=udp.Init(localIp)
	processChan := make(chan utilities.Message,50)
	achnowledge := make(chan utilities.Message,50)
	
	msgSentMap := make(map[int]utilities.Message)

	udpBroadcastMsg,udpRecvMsg := make(chan []byte,50), make(chan udp.RawMessage,50)
	go func(){
		for{
			select{
			case msg:=<-udpBroadcastMsg: 
				 udpRecvMsg<-udp.RawMessage{Data:msg,Ip:localIp}
			}
		}
	}()

	var mutex = &sync.Mutex{}
	var idMutex = &sync.Mutex{}

	go handel_UDP_message(processChan, recMsg,msgSentMap, achnowledge, mutex, idMutex)
	go resend_Msg(msgSentMap,udpBroadcastMsg, mutex)




	for{		
		select {
			case msg := <-sendMsg: 
				if msg.Message_Id == 0{
					idMutex.Lock()
					msg.Message_Id=msg_Id+1
					msg_Id+=1
					log.Println("Msg_id sent message: ", msg_Id)
					idMutex.Unlock()
				}
				buf:=utilities.Encoder(msg)
				udpBroadcastMsg<-buf
				
				if (msg.MessageType != utilities.MESSAGE_ACKNOWLEDGE) && (msg.MessageType!=utilities.MESSAGE_HEARTBEAT){
				mutex.Lock()
					msgSentMap[msg.Message_Id]=msg // Comfirmation map if message has been recieved, delete order when order is achnowledged
				mutex.Unlock()
				}



			case raw_m := <-udpRecvMsg:
				msg:=utilities.Decoder(raw_m.Data)
				msg.Message_sender=raw_m.Ip
				processChan<-msg



			case ach_m := <-achnowledge:
				ach_m.MessageType = utilities.MESSAGE_ACKNOWLEDGE
				udpBroadcastMsg<-utilities.Encoder(ach_m)
				

		}
	}

}


func handel_UDP_message(msg <-chan utilities.Message,
	rec_msg chan<-utilities.Message,
	msgMap map[int]utilities.Message,
	achnowledge chan<- utilities.Message, mutex * sync.Mutex,
	id_mutex * sync.Mutex ) {
	

	//s := make([])


	for{
		select{
			case msg:=<-msg:
				switch msg.MessageType{

					case utilities.MESSAGE_ACKNOWLEDGE: 
						log.Println("Achnowledgement for message :",msg.Message_Idq)
						mutex.Lock()
						delete(msgMap,msg.Message_Id) //Achnowledgement recieved, No need to resend message anymore
						mutex.Unlock()

					case utilities.MESSAGE_HEARTBEAT: 
						log.Println("Heartbeat recieved")
						Heartbeat_recieved(msg)
					



					default:
						achnowledge<-msg
						rec_msg<-msg //Sends the message to the manager
						

						id_mutex.Lock()
						msg_Id+=1 
						id_mutex.Unlock()





			}
		}
	}
}


func resend_Msg(msgMap map[int]utilities.Message, msg chan<-[]byte, mutex * sync.Mutex ) {
	for{
		time.Sleep(20*time.Millisecond)
		mutex.Lock()
		for _, v := range msgMap{
			msg<-utilities.Encoder(v)

		}
		mutex.Unlock()
	}
}



