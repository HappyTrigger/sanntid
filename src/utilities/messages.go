package utilities

import(
	//".././driver"
	".././dummydriver"
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





