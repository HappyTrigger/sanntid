package utilities

import(
	".././mydriver"
)



const (
	MESSAGE_ACKNOWLEDGE  	= iota
	MESSAGE_ORDER 			
	MESSAGE_STATE 		   	
	MESSAGE_ORDER_COMPLETE 	
	MESSAGE_HEARTBEAT		

)


type Achnowledgement struct {
	Checksum int
	Ip string
}


type NewOrder struct{
	Floor int
	Button driver.ButtonType
	Checksum int
}


type State struct{
	CurrentFloor int
	Direction int
	InternalOrders []int
	Door_open bool
	Ip string
	
}


func CreateHeartbeat(counter int) Heartbeat {
	return Heartbeat{Counter: counter}
}

type Heartbeat struct{
	Counter int
	Ip string
}




