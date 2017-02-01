package networking


import (

	"./udp"
    "fmt"
    "time"
	//"strconv"
	"log"
	".././communication"

)

var localIp string


func init() {

	var err error

	localIp, err = getLocalIp()

	if err != nil {
		log.Fatal("Refusing to start - ",err)
	}
}

func Run(){


	log.Println("---Starting network loop---")
	log.Println("The ip of this computer is: ", localIp)	

	udpBroadcastMsg,udpRecvMsg:=udp.Init(localIp)

	//i := 3
	msg2 := Communication.NewOrder{Floor:1,Direction:1}
	msg3 := Communication.NewOrder{Floor:0,Direction:0}	
	//raw_m := udp.RawMessage{}

	for{

		buf,err:=Communication.Encoder(msg2)
		if err==nil{
			udpBroadcastMsg<-buf
		}
		msg2.Floor = msg2.Floor + 1

		//buf := []byte(msg)
		time.Sleep(3*time.Second)
		raw_m:=<-udpRecvMsg
		Communication.Decoder(raw_m.Data,&msg3)
		fmt.Printf("%+v\n", msg3.Floor)

		}


}




