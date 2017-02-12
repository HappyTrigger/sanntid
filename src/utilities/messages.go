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



type Message struct{
	Message_origin string
	Message_sender string
	Message_Id int 
	MessageType int 
	State State
	NewOrder NewOrder
	Heartbeat Heartbeat
}


type NewOrder struct{
	Floor int
	Button driver.ButtonType 
}
type State struct{
	CurrentFloor int
	Direction int
	InternalOrders []int
	Door_open bool
	


}

type ConnectionStatus struct{
	Ip string
	Connection bool
}


func CreateHeartbeat(counter int) Heartbeat {
	return Heartbeat{Counter: counter}
}

type Heartbeat struct{
	Counter int
}




