package utilities

import(
	".././driver"
	//".././dummydriver"
	)


type Achnowledgement struct {
	Checksum int
	Id string
}

type State struct{
	Id string
	LastPassedFloor int
	Direction driver.ButtonType
	DoorState bool
	BetweenFloors bool
	InternalOrders []driver.OrderEvent
	StateSentFromId string
	
}



// Between floors can probably be rewritten to Elevator-Status or something like that
//Where it is active when serving any kind or order, and Idle when not doing anything.


