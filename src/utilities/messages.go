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
	CurrentFloor int
	Direction int
	InternalOrders []driver.OrderEvent
	
	Door_open bool
	BetweenFloors bool
	ConnectionStatus bool



	StateSentFromIp string
	
}





