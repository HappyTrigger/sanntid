package main


import (	
	"./utilities"
	"./manager"
	"./elevator"
	//"log"
	"./mydriver"
	//"./dummydriver"
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



	NewState := make(chan utilities.State)

	//manager
	DriverEvent := make(chan driver.OrderEvent)
	SendOrderToElevator := make(chan driver.OrderEvent)
	elevatorOrderComplete :=make(chan driver.OrderEvent)


	//Elevator
	ButtonStop := make(chan bool)
	DoorOpen := make(chan bool)
	DoorClosed := make(chan bool)
	SensorEvent := make(chan int)
	ElevatorEmergency := make(chan bool)





	go manager.Run(
		SendOrderToElevator,
		DriverEvent,
		DoorOpen,
		DoorClosed,
		ElevatorEmergency,
		elevatorOrderComplete)
	



	driver.Init(DriverEvent,SensorEvent,ButtonStop)
	

	go elevator.Run(NewState,
		SendOrderToElevator,
		SensorEvent,
		ButtonStop,
		DoorOpen,
		DoorClosed,
		ElevatorEmergency,
		elevatorOrderComplete)


	select {

	}
}