package networking


import (

	"./udp"
    "fmt"
    "time"
	"strconv"
	"log"

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

	localIp = "129.241.187.38"
	udpBroadcastMsg,udpRecvMsg:=udp.Init(localIp)

	i := 3
	for{
		msg := strconv.Itoa(i)
		i++
		buf := []byte(msg)
		udpBroadcastMsg<-buf
		time.Sleep(3*time.Second)

		fmt.Printf("%+v\n", <-udpRecvMsg)
		}



		

	
}




