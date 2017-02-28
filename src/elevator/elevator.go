package elevator

import (
	".././utilities"
	"log"
	"time"
	".././driver"
	//".././dummydriver"
)

const(
	DoorOpenTime = 2*time.Second
	No_order = -1
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

	*DoorState		= false
	*BetweenFloors 	= false


	var doorClose <-chan time.Time
	Orders := make(map[int]driver.OrderEvent)



	*lastPassedFloor = driver.Elev_get_floor_sensor_signal()
	if *lastPassedFloor == -1{
		log.Fatal("[FATAL]\tElevator initialized between floors")
	}

	for{
		select{

		case order:=<-NewOrder:
			Orders[order.Checksum] = order
			driver.Elev_set_button_lamp(order.Button,order.Floor,true)
			log.Println("Order delegated to this elevator")

			

			if !*BetweenFloors {
				
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
			

		}
	}
}








//This function can probably be rewritten to half its length, if you just irierate over one segment twice.
func OrderOnTheFloor(orders map[int]driver.OrderEvent,
	Direction* driver.ButtonType,
	LastPassedFloor* int,
	OrderComplete chan<-driver.OrderEvent)(int,bool){


	
	var orderOnNextFloors bool
	var orderOnFloor int 

	orderOnNextFloors = false
	orderOnFloor = -1

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
	
	if orderOnFloor !=-1 {
		driver.Elev_set_motor_direction(driver.MotorStop)
		driver.Elev_set_door_open_lamp(true)
		*doorClose = time.After(DoorOpenTime)
		*DoorState=true
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
			*BetweenFloors=false
		}
	}
}






func sendState() utilities.State{
	temp_state := utilities.State{}
	temp_state = ElevatorState
	return temp_state

}


