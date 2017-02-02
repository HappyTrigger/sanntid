package networking

import(
	"log"
	".././utilities"
)


func SendHeartBeat(counter int) {
	log.Println("Sent heartbeat")
}

func Heartbeat_recieved(m utilities.Message) {
	log.Println("Recieved Heartbeat")
}
