package main

//Pull from https://github.com/johhat/sanntidsheis if needed

import (	
	"./networking"
	"./utilities"
	"./manager"
	"./elevator"
	"log"
	"./mydriver"
)

func main() {
	
	fromManager := make(chan utilities.Message)
	toManager := make(chan utilities.Message)
	connectionStatus := make(chan utilities.ConnectionStatus)


	NewState := make(chan utilities.State)

	DriverEvent := make(chan utilities.NewOrder)

	SendOrderToElevator := make(chan utilities.NewOrder)
	



	OrderEvent := make(chan driver.OrderEvent)
	SensorEvent := make(chan int)

	ButtonStop := make(chan bool)
	OpenDoor := make(chan bool)





	go manager.Run(fromManager,
		toManager, 
		connectionStatus,
		NewState,
		DriverEvent,
		SendOrderToElevator)
	

	go networking.Run(fromManager,
		toManager,
		connectionStatus)	


	go elevator.Run(NewState,SendOrderToElevator,SensorEvent)


	go driver.Init(OrderEvent,SensorEvent,ButtonStop)


	for{

	}
	log.Println("Exiting program")
}