package elevator

import (
	".././utilities"
	"log"
	"time"
	".././mydriver"
	//".././dummydriver"
)


const(
	DoorOpenTime = 2*time.Second
) 

var ElevatorState utilities.State //:= utilities.State{}


func Run(
	NewOrder <-chan driver.OrderEvent,
	SensorEvent <-chan int,
	ElevatorEmergency <-chan bool,
	OrderComplete chan<-driver.OrderEvent,
	ElevatorStateToManager chan<- utilities.State,
	StopButton<-chan bool) {


	var State State
	var doorClose <-chan time.Time

	Direction := &ElevatorState.Direction
	lastPassedFloor := &ElevatorState.LastPassedFloor 
	DoorState := &ElevatorState.DoorState 
	BetweenFloors := &ElevatorState.BetweenFloors 

	Orders := make(map[int]driver.OrderEvent)
	StateChange := make(chan bool)




	*lastPassedFloor = driver.Elev_get_floor_sensor_signal()
	if *lastPassedFloor == -1{
		log.Fatal("[FATAL]\tElevator initialized between floors")
	}




	//Regarding the state-changes that must be sent to the manager. Those can either be implemented as
	// channels, seperating every single event into different channels, and updating the state of 
	// the elevator when reciveving new messages on those channels, or sending the complete state from the elevator modul


	for{

		select{

		case order:=<-NewOrder:
			Orders[order.Checksum] = order
			driver.Elev_set_button_lamp(order.Button,order.Floor,true)
			log.Println("New order in elevator")


			switch State {

				case State_idle:
				orderOnFloor,orderOnNextFloors:=OrderOnTheFloor(Orders,
				*Direction,*lastPassedFloor,OrderComplete)
				if orderOnFloor !=-1 {
					driver.Elev_set_motor_direction(driver.MotorStop)
					driver.Elev_set_door_open_lamp(true)
					doorClose = time.After(DoorOpenTime)
					*DoorState=true
					*BetweenFloors=false
				}else{
					if orderOnNextFloors{
						if *Direction==driver.Down{
							driver.Elev_set_motor_direction(driver.MotorDown)
						}else{
							driver.Elev_set_motor_direction(driver.MotorUp)
						}
						*BetweenFloors=true


						//Send new state here including current floor, direction, active
						
					}
				}
				StateChange<-true


				case State_moving:
					//Do nothing
				}

			


		case *lastPassedFloor = <-SensorEvent:
			driver.Elev_set_floor_indicator(*lastPassedFloor)
			if *lastPassedFloor !=-1 {

				orderOnFloor,orderOnNextFloors:=OrderOnTheFloor(Orders,
				*Direction,*lastPassedFloor,OrderComplete)
			
				if orderOnFloor !=-1 {
					driver.Elev_set_motor_direction(driver.MotorStop)
					driver.Elev_set_door_open_lamp(true)
					doorClose = time.After(3*time.Second)
					*DoorState=true
					*BetweenFloors=false
				}else{
					if orderOnNextFloors{
						if *Direction==driver.Down{
							driver.Elev_set_motor_direction(driver.MotorDown)
						}else{
							driver.Elev_set_motor_direction(driver.MotorUp)
						}
						
					}
					*BetweenFloors=true
				}
			}
			StateChange<-true




		case stop:=<-StopButton:
			log.Println("Stop butten has been pressed:",stop)

		case <-StateChange:
			temp:=ElevatorState
			ElevatorStateToManager<-temp



		case <-doorClose:
			driver.Elev_set_door_open_lamp(false)
			*DoorState=false
			orderOnFloor, orderOnNextFloors := 
				OrderOnTheFloor(Orders, *Direction,*lastPassedFloor,OrderComplete)
			
			if orderOnFloor !=-1 {
				driver.Elev_set_motor_direction(driver.MotorStop)
				driver.Elev_set_door_open_lamp(true)
				doorClose = time.After(DoorOpenTime)
				*DoorState=true
				*BetweenFloors=false
			} else {
				if orderOnNextFloors{
					if *Direction == driver.Up{
						driver.Elev_set_motor_direction(driver.MotorUp)
					}else{
						driver.Elev_set_motor_direction(driver.MotorDown)
					}
					State=State_moving
					*BetweenFloors=true
				}else{
					State=State_idle
				}
			}
		}
		StateChange<-true
	}
}

func OrderOnTheFloor(orders map[int]driver.OrderEvent,
	Direction driver.ButtonType,
	CurrentFloor int,
	OrderComplete chan<-driver.OrderEvent)(int,bool){


	//log.Println("Order Checking func started")
	No_order := -1
	var orderOnNextFloors bool
	var orderOnFloor int 

	orderOnNextFloors = false
	orderOnFloor = -1

	for k,v:= range orders{
		//log.Println("Checking")
		if v.Button == Direction || v.Button == driver.Internal{
			if CurrentFloor == v.Floor{
				orderOnFloor=k //Send orderNumber back
				driver.Elev_set_button_lamp(v.Button,v.Floor,false)
				log.Println("Order on floor confirmed")
				OrderComplete<-v
				delete(orders,k)
			}
			
			if Direction == driver.Up{
				if v.Floor>CurrentFloor{
					orderOnNextFloors=true
					log.Println("Order on the floor above")
				}
			}else{
				if v.Floor<CurrentFloor{
					orderOnNextFloors=true
					log.Println("Order on the floor below")
				}
			}
		}else{

			switch Direction{
				case driver.Up:
					if v.Floor > CurrentFloor{
							orderOnNextFloors=true
					}

				case driver.Down:
					if v.Floor < CurrentFloor{
						orderOnNextFloors=true
							
					}
				}
			}
		}

	if orderOnFloor == No_order && orderOnNextFloors==false{
		if Direction == driver.Up{
			Direction = driver.Down
		}else{
			Direction = driver.Up
		}
		for k,v:= range orders{
			if v.Button == Direction || v.Button == driver.Internal{
				if CurrentFloor == v.Floor{
					orderOnFloor=k //Send orderNumber back
					driver.Elev_set_button_lamp(v.Button,v.Floor,false)
					OrderComplete<-v
					delete(orders,k)
					log.Println("Order on floor confirmed")
				}
				if Direction == driver.Up{
					if v.Floor>CurrentFloor{
						orderOnNextFloors=true
						log.Println("Order on the floor above")
					}
				}else{
					if v.Floor<CurrentFloor{
						orderOnNextFloors=true
						log.Println("Order on the floor below")
					}
				}
			}else{

			switch Direction{
				case driver.Up:
					if v.Floor > CurrentFloor{
							orderOnNextFloors=true
					}

				case driver.Down:
					if v.Floor < CurrentFloor{
						orderOnNextFloors=true
							
					}
				}
			}
		}
	}
	return orderOnFloor,orderOnNextFloors
}


func elevatorControl(DoorState bool,
	BetweenFloors bool,
	Direction driver.ButtonType,
	StateChange chan<-bool,
	doorClose<-chan time.Time,
	Orders map[int]driver.OrderEvent,
	lastPassedFloor int,
	OrderComplete chan<-driver.OrderEvent,
	State State){
	
	
	orderOnFloor,orderOnNextFloors := OrderOnTheFloor(Orders,Direction,lastPassedFloor,OrderComplete)
	
	if orderOnFloor !=-1 {
		driver.Elev_set_motor_direction(driver.MotorStop)
		driver.Elev_set_door_open_lamp(true)
		doorClose = time.After(DoorOpenTime)
		DoorState=true
		BetweenFloors=false
	} else {
		if orderOnNextFloors{
			if Direction == driver.Up{
				driver.Elev_set_motor_direction(driver.MotorUp)
			}else{
				driver.Elev_set_motor_direction(driver.MotorDown)
			}
				State=State_moving
				BetweenFloors=true
		}else{
			State=State_idle
			BetweenFloors=false
		}
	}
	StateChange<-true
}









