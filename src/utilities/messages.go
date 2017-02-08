package utilities



const (
	MESSAGE_ACKNOWLEDGE  	= 0
	MESSAGE_ORDER 			= 1
	MESSAGE_STATE 		   	= 2
	MESSAGE_ORDER_COMPLETE 	= 3
	MESSAGE_HEARTBEAT		= 4

)

const(
	DIR_UP 		= 1
	DIR_DOWN 	= -1
	STANDSTILL 	= 0
)




type Message struct{
	Message_sender string
	Message_Id int 
	MessageType int 
	Acknowledge Acknowledge
	State State
	NewOrder NewOrder
	Heartbeat Heartbeat
}

type Acknowledge struct{
	Message_Id int
}

type NewOrder struct{
	Floor int
	Direction int
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




