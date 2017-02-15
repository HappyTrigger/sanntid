package main

//Pull from https://github.com/johhat/sanntidsheis if needed

import (	
	"./utilities"
	"./manager"
	"./elevator"
	//"log"
	"./mydriver"
	//"os"
	//"os/signal"
)

func main() {
	



	NewState := make(chan utilities.State)


	DriverEvent := make(chan driver.OrderEvent)
	SensorEvent := make(chan int)

	ButtonStop := make(chan bool)
	DoorOpen := make(chan bool)
	DoorClosed := make(chan bool)

	//ReachedNewFloor := make(chan int)
	ElevatorEmergency := make(chan bool)

	SendOrderToElevator := make(chan driver.OrderEvent)





	go manager.Run(
		SendOrderToElevator,
		DriverEvent,
		DoorOpen,
		DoorClosed,
		ElevatorEmergency)
	



	driver.Init(DriverEvent,SensorEvent,ButtonStop)
	

	go elevator.Run(NewState,
		SendOrderToElevator,
		SensorEvent,
		ButtonStop,
		DoorOpen,
		DoorClosed,
		ElevatorEmergency)



	/*c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		driver.Elev_set_motor_direction(driver.MotorStop)
		log.Fatal("[FATAL]\tUser terminated program")

		}()
*/
	for{

	}
}