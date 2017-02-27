package main


import (	
	"./utilities"
	"./manager"
	"./elevator"
	//"log"
	//"./driver"
	"./dummydriver"
	"os"
	"os/signal"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
	    <-c
	    panic("Ctrl+C")
	}()



	//manager
	DriverEvent := make(chan driver.OrderEvent)
	SendOrderToElevator := make(chan driver.OrderEvent)
	ElevatorOrderComplete :=make(chan driver.OrderEvent)


	//Elevator
	SensorEvent := make(chan int)
	ElevatorEmergency := make(chan bool)


	//statetransfer
	ElevatorState := make(chan utilities.State)
	StopButton := make(chan bool)


	go manager.Run(
		SendOrderToElevator,
		DriverEvent,
		ElevatorEmergency,
		ElevatorOrderComplete,
		ElevatorState)
	



	driver.Init(DriverEvent,SensorEvent,StopButton)
	

	go elevator.Run(
		SendOrderToElevator,
		SensorEvent,
		ElevatorEmergency,
		ElevatorOrderComplete,
		ElevatorState,
		StopButton)


	select {

	}
}