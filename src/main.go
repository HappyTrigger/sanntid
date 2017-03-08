package main

/*
Main functions as a connectionpoint between modules, creating and sharing channels based on their needs.
*/

import (
	"flag"
	"os"
	"os/signal"

	"./driver"
	"./elevator"
	"./manager"
	"./utilities"
)

func main() {
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		panic("Ctrl+C")
	}()

	//manager
	DriverEvent := make(chan driver.OrderEvent)
	SendOrderToElevator := make(chan driver.OrderEvent)
	ElevatorOrderComplete := make(chan driver.OrderEvent)

	//Elevator
	SensorEvent := make(chan int)
	ElevatorEmergency := make(chan bool)

	//statetransfer
	ElevatorState := make(chan utilities.State)
	StopButton := make(chan bool)

	driver.Init(DriverEvent, SensorEvent, StopButton)

	go elevator.Run(
		SendOrderToElevator,
		SensorEvent,
		ElevatorEmergency,
		ElevatorOrderComplete,
		ElevatorState,
		StopButton)

	go manager.Run(
		SendOrderToElevator,
		DriverEvent,
		ElevatorEmergency,
		ElevatorOrderComplete,
		ElevatorState,
		id)

	select {}
}
