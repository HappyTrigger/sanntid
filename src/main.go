package main

//Pull from https://github.com/johhat/sanntidsheis if needed

import (	
	"./utilities"
	"./manager"
	"./elevator"
	//"log"
	"./mydriver"
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
	orderComplete :=make(chan utilities.NewOrder)


	//Elevator
	ButtonStop := make(chan bool)
	DoorOpen := make(chan bool)
	DoorClosed := make(chan bool)
	SensorEvent := make(chan int)
	//ReachedNewFloor := make(chan int)
	ElevatorEmergency := make(chan bool)





	go manager.Run(
		SendOrderToElevator,
		DriverEvent,
		DoorOpen,
		DoorClosed,
		ElevatorEmergency,
		orderComplete)
	



	driver.Init(DriverEvent,SensorEvent,ButtonStop)
	

	go elevator.Run(NewState,
		SendOrderToElevator,
		SensorEvent,
		ButtonStop,
		DoorOpen,
		DoorClosed,
		ElevatorEmergency,
		orderComplete)



	/*c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		driver.Elev_set_motor_direction(driver.MotorStop)
		log.Fatal("[FATAL]\tUser terminated program")

		}()
*/
	select {

	}
}