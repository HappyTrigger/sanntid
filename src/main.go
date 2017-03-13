package main

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

	DriverEvent := make(chan driver.OrderEvent)
	SensorEvent := make(chan int)
	StopButton := make(chan bool)

	SendOrderToElevator := make(chan driver.OrderEvent)
	ElevatorOrderComplete := make(chan driver.OrderEvent)
	ElevatorEmergency := make(chan bool)
	ElevatorState := make(chan utilities.State)

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
