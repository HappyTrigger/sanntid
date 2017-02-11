package elevator

import (
	".././utilities"
	"log"
	"time"
	".././mydriver"
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
var floor int

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
				if sensor !=-1{
					driver.Elev_set_floor_indicator(sensor)
					driver.Elev_set_motor_direction(driver.MotorStop)
				}






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
				log.Println("Initializing")
				state=State_OnFloor

			

			case State_Idle:

			case State_OnFloor:
				log.Println("Stopping")
				driver.Elev_set_motor_direction(driver.MotorStop)
				time.Sleep(3*time.Second)

				if floor >=3 {
					driver.Elev_set_motor_direction(driver.MotorDown)
				}else{
					driver.Elev_set_motor_direction(driver.MotorUp)
				}
				log.Println("Driving")






			

			case State_Failiure:


		}
	}
}


