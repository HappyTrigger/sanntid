package utilities

import (
	"../driver"
	//"../dummydriver"
)

type Achnowledgement struct {
	Checksum int
	Id       string
}

type State struct {
	Id                 string
	LastRegisterdFloor int
	Direction          driver.ButtonType
	DoorState          bool
	Idle               bool
	InternalOrders     []driver.OrderEvent
	StateSentFromId    string
}
