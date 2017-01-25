package UDP

import (
	"net"
	"strings"
	"encoding/json"
	"fmt"
	"time"



	)
const (
	PORT = ":20022"
	PORTLISTENER = ":30000"
	)

type Message struct {
	Ip int
	Data int
}

func UDPSender(broadcast_chan chan int) {
	broadcastAddr := []string{"129.241.187.255", PORT}
	broadcastUDP, _ := net.ResolveUDPAddr("udp", strings.Join(broadcastAddr, ""))
	broadcastConn, _ := net.DialUDP("udp", nil, broadcastUDP) //returns the UDP connection interface which supports reading and writing
	defer broadcastConn.Close() // Plan the connection to be closed soon
	for {
		buf, err := json.Marshal(<- broadcast_chan) //codes the message
		if err == nil {
			broadcastConn.Write(buf)
		}
	}
}

func UDPListener(recieve_chan chan int) {
	UDPReceiveAddr, err := net.ResolveUDPAddr("udp", PORTLISTENER);
	if err != nil { fmt.Println(err) }

	UDPConn, err := net.ListenUDP("udp", UDPReceiveAddr);
	if err != nil { fmt.Println(err) }
	defer UDPConn.Close()

	buf := make([]byte, 2048)
	trimmed_buf := make([]byte, 1)
	var received_message int

	for {
		n, _, _ := UDPConn.ReadFromUDP(buf)
		trimmed_buf = buf[:n]
		err := json.Unmarshal(trimmed_buf, &received_message) //decodes the message
		if err == nil {
			recieve_chan <- received_message
		}
	}
}


func testUDP(broadcast_chan chan int, done_chan chan bool){
	var tick = 0
	for{
		time.Sleep(3)
		broadcast_chan<-5
		fmt.Println("hei")
		if tick==5{
			done_chan<-true
		}


	}
}


func main() {

	message_recieved := make(chan int,100)
	message_send := make(chan int,100)
	done :=make(chan bool)
	


	

	

	go UDPSender(message_send)
	go UDPListener(message_recieved)
	go testUDP(message_send, done)



	go func() {
		for {
			select {
			case msg := <-message_recieved:
				fmt.Printf("%d",msg)

			}

		}
	}()

	<-done
	
}