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





