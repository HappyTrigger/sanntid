package utilities



const (
	MESSAGE_ACKNOLEDGE  	= 0
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
	Message_Id int
	MessageType int 
	Acknoledge Acknoledge
	State State
	NewOrder NewOrder
	Heartbeat Heartbeat
}

type Acknoledge struct{
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

const Heartbeat_code = "Gruppe 45" 

func CreateHeartbeat(counter int) Heartbeat {
	return Heartbeat{Code: Heartbeat_code, Counter: counter}
}

type Heartbeat struct{
	Counter int
	Code string
}




