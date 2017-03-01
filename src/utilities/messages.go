package utilities

import(
	".././driver"
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
	DoorState bool
	BetweenFloors bool
	InternalOrders []driver.OrderEvent
	StateSentFromIp string
	
}



// Between floors can probably be rewritten to Elevator-Status or something like that
//Where it is active when serving any kind or order, and Idle when not doing anything.


