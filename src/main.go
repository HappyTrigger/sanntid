package main

//Pull from https://github.com/johhat/sanntidsheis if needed

import (	
	"./networking"
	"./utilities"
	"./manager"
	"./elevator"
	//"log"
	"./mydriver"
	//"os"
	//"os/signal"
)

func main() {
	
	FromManager := make(chan utilities.Message)
	ToManager := make(chan utilities.Message)
	ConnectionStatus := make(chan utilities.ConnectionStatus)


	NewState := make(chan utilities.State)


	SendOrderToElevator := make(chan driver.OrderEvent)
	

	DriverEvent := make(chan driver.OrderEvent)
	SensorEvent := make(chan int)

	ButtonStop := make(chan bool)
	DoorOpen := make(chan bool)
	DoorClosed := make(chan bool)

	//ReachedNewFloor := make(chan int)
	ElevatorEmergency := make(chan bool)






	go manager.Run(
		FromManager,
		ToManager, 
		ConnectionStatus,
		DriverEvent,
		SendOrderToElevator,
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



	go networking.Run(FromManager,
		ToManager,
		ConnectionStatus)	

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