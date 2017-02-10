package main

//Pull from https://github.com/johhat/sanntidsheis if needed

import (	
	"./networking"
	"./utilities"
	"./manager"
	//"./elevator"
	"log"
)

func main() {
	
	fromManager := make(chan utilities.Message)
	toManager := make(chan utilities.Message)
	connectionStatus := make(chan utilities.ConnectionStatus)


	NewState := make(chan utilities.State)

	DriverEvent := make(chan utilities.NewOrder)

	SendOrderToElevator := make(chan utilities.NewOrder)
	


	go manager.Run(fromManager,
		toManager, 
		connectionStatus,
		NewState,
		DriverEvent,
		SendOrderToElevator)
	

	go networking.Run(fromManager,
		toManager,
		connectionStatus)	
	for{

	}
	log.Println("Exiting program")
}