package elevator

import (
	".././utilities"
	"log"
	"time"
	".././driver"
	//".././dummydriver"
	"reflect"
)

const(
	DoorOpenTime = 2*time.Second
	No_order = -1
	//StatechangeInterval = 4*time.Second
) 

var ElevatorState utilities.State

func Run(
	NewOrder <-chan driver.OrderEvent,
	SensorEvent <-chan int,
	ElevatorEmergency chan<- bool,
	OrderComplete chan<-driver.OrderEvent,
	ElevatorStateToManager chan<- utilities.State,
	StopButton<-chan bool) {


	Direction 		:= &ElevatorState.Direction
	lastPassedFloor := &ElevatorState.LastPassedFloor 
	DoorState 		:= &ElevatorState.DoorState 
	BetweenFloors 	:= &ElevatorState.BetweenFloors 

	var doorClose <-chan time.Time
	//var tempElevatorState utilities.State
	//StateChangeTimer := time.Tick(StatechangeInterval)
	Orders := make(map[int]driver.OrderEvent)

	*DoorState		= false
	*BetweenFloors 	= false
	*Direction 		= driver.Down
	*lastPassedFloor = driver.Elev_get_floor_sensor_signal()
	//tempElevatorState = ElevatorState

	if *lastPassedFloor == -1{
		log.Fatal("[FATAL]\tElevator initialized between floors")
	}




	for{
		select{

		case order:=<-NewOrder:
			Orders[order.Checksum] = order
			driver.Elev_set_button_lamp(order.Button,order.Floor,true)
			log.Println("Order delegated to this elevator")

			

			if !*BetweenFloors && !*DoorState {
				
				elevatorControl(
					DoorState ,
					BetweenFloors ,
					Direction,
					&doorClose,
					Orders,
					lastPassedFloor,
					OrderComplete)

				ElevatorStateToManager<-sendState()
				}




		case *lastPassedFloor = <-SensorEvent:
			log.Println("On Floor ", *lastPassedFloor)
			driver.Elev_set_floor_indicator(*lastPassedFloor)
			if *lastPassedFloor !=-1 {
				elevatorControl(
					DoorState,
					BetweenFloors,
					Direction,
					&doorClose,
					Orders,
					lastPassedFloor,
					OrderComplete)	

				ElevatorStateToManager<-sendState()
				}
				

		case <-doorClose:
			driver.Elev_set_door_open_lamp(false)
			*DoorState=false
			elevatorControl(
					DoorState,
					BetweenFloors,
					Direction,
					&doorClose,
					Orders,
					lastPassedFloor,
					OrderComplete)

			ElevatorStateToManager<-sendState()



		case <-StopButton:
			ElevatorEmergency<-true

/*
		case <-StateChangeTimer:
			if !*BetweenFloors && !*DoorState{
				//system is idle
			}else{
				if reflect.DeepEqual(ElevatorState,tempElevatorState){
					ElevatorEmergency<-true
				}
				tempElevatorState=ElevatorState
			}

*/


			

		}
	}
}








//This function can probably be rewritten to half its length, if you just irierate over one segment twice, but change
//the direction for each iteration if no order in the current direction is found.
func OrderOnTheFloor(orders map[int]driver.OrderEvent,
	Direction* driver.ButtonType,
	LastPassedFloor* int,
	OrderComplete chan<-driver.OrderEvent)(int,bool){
	
	var orderOnNextFloors bool
	var orderOnFloor int 

	orderOnNextFloors = false
	orderOnFloor = -1

	// k and v should we written to checksum and order
	for k,v:= range orders{
		if v.Button == *Direction || v.Button == driver.Internal{
			if *LastPassedFloor == v.Floor{
				orderOnFloor=k //Send orderNumber back
				driver.Elev_set_button_lamp(v.Button,v.Floor,false)
				log.Println("Order on floor confirmed")
				OrderComplete<-v
				delete(orders,k)
			}
			if *Direction == driver.Up{
				if v.Floor>*LastPassedFloor{
					orderOnNextFloors=true
					log.Println("Order on the floor above")
				}
			}else{
				if v.Floor<*LastPassedFloor{
					orderOnNextFloors=true
					log.Println("Order on the floor below")
				}
			}
		}else{
			
			switch *Direction{
				case driver.Up:
					if v.Floor > *LastPassedFloor{
							orderOnNextFloors=true
					}

				case driver.Down:
					if v.Floor < *LastPassedFloor{
						orderOnNextFloors=true
							
					}
				}
			}
		}

	if orderOnFloor == No_order && orderOnNextFloors==false{
		if *Direction == driver.Up{
			*Direction = driver.Down
		}else{
			*Direction = driver.Up
		}
		for k,v:= range orders{
			if v.Button == *Direction || v.Button == driver.Internal{
				if *LastPassedFloor == v.Floor{
					orderOnFloor=k //Send orderNumber back
					driver.Elev_set_button_lamp(v.Button,v.Floor,false)
					OrderComplete<-v
					delete(orders,k)
					log.Println("Order on floor confirmed")
				}
				if *Direction == driver.Up{
					if v.Floor>*LastPassedFloor{
						orderOnNextFloors=true
						log.Println("Order on the floor above")
					}
				}else{
					if v.Floor<*LastPassedFloor{
						orderOnNextFloors=true
						log.Println("Order on the floor below")
					}
				}
			}else{

			switch *Direction{
				case driver.Up:
					if v.Floor > *LastPassedFloor{
							orderOnNextFloors=true
					}

				case driver.Down:
					if v.Floor < *LastPassedFloor{
						orderOnNextFloors=true
							
					}
				}
			}
		}
	}
	return orderOnFloor,orderOnNextFloors
}


func elevatorControl(DoorState* bool,
	BetweenFloors* bool,
	Direction* driver.ButtonType,
	doorClose* <-chan time.Time,
	Orders map[int]driver.OrderEvent,
	lastPassedFloor* int,
	OrderComplete chan<-driver.OrderEvent){
	
	
	orderOnFloor,orderOnNextFloors := OrderOnTheFloor(Orders,Direction,lastPassedFloor,OrderComplete)
	
	if orderOnFloor !=No_order {
		driver.Elev_set_motor_direction(driver.MotorStop)
		driver.Elev_set_door_open_lamp(true)
		*doorClose = time.After(DoorOpenTime)
		*DoorState=true //Should probably create a type for this, so that you can write Doorstate = DoorOpen or something like that
	} else {
		if orderOnNextFloors{
			if *Direction == driver.Up{
				driver.Elev_set_motor_direction(driver.MotorUp)
				log.Println("Driving up")
			}else{
				driver.Elev_set_motor_direction(driver.MotorDown)
				log.Println("Driving Down")
			}
			*BetweenFloors=true
		
		}else{
			//Between floors in only set to false when it is idle. This has to do with the delegation-alogrithm. 
			*BetweenFloors=false

		}
	}
}






func sendState() utilities.State{
	temp_state := utilities.State{}
	temp_state = ElevatorState
	return temp_state

}


