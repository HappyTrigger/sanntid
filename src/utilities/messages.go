package utilities

import(
	".././mydriver"
	//".././dummydriver"
	)


type Achnowledgement struct {
	Checksum int
	Ip string
}

type State struct{
	Ip string
	LastPassedFloor int
	Direction driver.ButtonType
	LastDestinationFloor int
	DoorState bool
	BetweenFloors bool
	EmergencyMode bool 
	InternalOrders []driver.OrderEvent
	CurrentExternalOrders [] driver.OrderEvent //Need this in cost function
	//To distinguis between the senders in special occasions
	StateSentFromIp string
	
}





