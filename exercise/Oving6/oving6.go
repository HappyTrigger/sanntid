package main


import(
	"bytes"
	"log"
	"net"
	"errors"
	//"strconv"
	//"strings"
	"time"
	 "encoding/binary"
)

type RawMessage struct {
	Data []byte
	Ip   string
}

const (
	broadcastAddress = "255.255.255.255:10701"
	listenPort       = ":30700"
)

func getLocalIp() (string, error) {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", errors.New("Could not get local ip")
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", errors.New("Could not get local ip")
}


func broadcast(broadcastChan <-chan []byte, localListener *net.UDPConn) {

	addr, _ := net.ResolveUDPAddr("udp", broadcastAddress)

	var b bytes.Buffer



	for msg := range broadcastChan {
		b.Write(msg)
		b.WriteRune('\n')

		_, err := localListener.WriteToUDP(b.Bytes(), addr)

		b.Reset()

		if err != nil {
			log.Println(err)
		}
	}
}


func recieve(recieveChan chan<- RawMessage, broadcastListener *net.UDPConn) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error in UDP recieve: %s \n Closing connection.", r)
			broadcastListener.Close()
		}
	}()

	buffer := make([]byte, 1024)

	for {
		n, address, err := broadcastListener.ReadFromUDP(buffer)

		if err != nil || n < 0 {
			log.Printf("Error in UDP recieve\n")
			panic(err)
		}

		data, err := bytes.NewBuffer(buffer).ReadBytes('\n')

		if err != nil {
			log.Println("Error when reading UDP message buffer:", err)
		}

		recieveChan <- RawMessage{Data: data, Ip: address.IP.String()}
	}
}


func Init(localIp string) (chan<- []byte, <-chan RawMessage){

	addr, _ := net.ResolveUDPAddr("udp", listenPort)

	localListener, err := net.ListenUDP("udp", addr)

	if err != nil {
		log.Fatal(err)
	}

	addr, _ = net.ResolveUDPAddr("udp", broadcastAddress)

	broadcastListener, err := net.ListenUDP("udp", addr)

	if err != nil {
		log.Fatal(err)
	}

	broadcastChan := make(chan []byte)
	go broadcast(broadcastChan, localListener)

	recieveChan := make(chan RawMessage)
	go recieve(recieveChan, broadcastListener)

	
	log.Println("UDP initialized")
   	
	


	return broadcastChan, recieveChan
}






func main() {

	localIp, err := getLocalIp() 
	if err != nil {
		log.Fatal("Refusing to start - ",err)
	}
	
	broadcast, recievechan := Init(localIp)

	amIMaster:=false
	var count uint64

	count = 0
	bs := make([]byte, 8)

	for{
		select{

		case recieve:=<-recievechan:
			if recieve.Ip != localIp{
				//count, _ := strconv.Atoi(string(recieve.Data))
				count  = binary.BigEndian.Uint64(recieve.Data)

				log.Println("Recieved count: ",count)
				amIMaster = false
			}
			break

			
		case <-time.After(4*time.Second):
			log.Println("Timed out \n")
			amIMaster = true
			break

		}
		if(amIMaster){
				binary.Write(bs,binary.LittleEndian ,count)
				broadcast<-bs
				log.Println("Broadcast Count:",count)
				count ++ 
				time.Sleep(1*time.Second)
				amIMaster = true
				
		}


	}







}




