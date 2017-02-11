package elevator

import (
	".././utilities"
	"log"
)


type State int

const(
	State_Init State = iota
	State_OnFloor
	State_Moving
	State_Idle
	State_Failiure
)


func Init() {
	// get his position
	// "let the manager write the InternalOrders in the case where rebooting and
	//another manager gives back the memory of interalorders"
}


var state State

func Run(NewState chan<-utilities.State,
	NewOrder <-chan utilities.NewOrder,
	SensorEvent <-chan int,
	StopButton <-chan bool) {



	state = State_Init
	log.Println("Starting elevator")

	go StateMachine()

	for{

		select{

			case order:=<-NewOrder:
				log.Println("New order from manager:",order)




			case sensor:=<-SensorEvent:
				log.Println("Elevator has reached new floor:",sensor)

			case stop:=<-StopButton:
				log.Println("Stop butten has been pressed:",stop)
				state = State_Failiure



			default:
				//Do nothing 

			}






	}





}


func StateMachine() {
	for{

		switch state{

			case State_Init:
				
			

			case 	State_Moving:
			

			case 	State_Idle:
			

			case 	State_Failiure:

		}
	}
}


