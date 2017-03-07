package elevator
/*
The elevator-module is responsible for everything regarding the elevator control and exectution.
It controls lighting, motorfunctioality, and keeps track of the current elevator-state. 
It also sends state-updates to the manager when a state-change has been registerd.

*/
import (
	".././utilities"
	"log"
	"time"
	".././driver"
	//".././dummydriver"

)

const(
	ElevatorEmergencyTimeInterval = 6*time.Second
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

	var doorClose <-chan time.Time
	elevatorEmergencyTimer := time.NewTimer(ElevatorEmergencyTimeInterval)
	Orders := make(map[int]driver.OrderEvent)

	*DoorState		= false
	*BetweenFloors 	= false
	*Direction 		= driver.Down
	*lastPassedFloor = driver.Elev_get_floor_sensor_signal()


	ElevatorStateToManager<-sendState() // Probably dont need the sendState function
	// Depends if structs are passed as values rather than reference

	if *lastPassedFloor == -1{
		log.Fatal("[FATAL]\tElevator initialized between floors")
	}



	for{
		select{

		case order:=<-NewOrder:
			Orders[order.Checksum] = order
			log.Println("Order delegated to this elevator")
			driver.Elev_set_button_lamp(order.Button,order.Floor,true)
			

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
				if !elevatorEmergencyTimer.Stop() {
					<-elevatorEmergencyTimer.C
				}
				elevatorEmergencyTimer.Reset(ElevatorEmergencyTimeInterval)
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
				if !elevatorEmergencyTimer.Stop() {
					<-elevatorEmergencyTimer.C
				}
				elevatorEmergencyTimer.Reset(ElevatorEmergencyTimeInterval)
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
			if !elevatorEmergencyTimer.Stop() {
				<-elevatorEmergencyTimer.C
			}
			elevatorEmergencyTimer.Reset(ElevatorEmergencyTimeInterval)


		case <-StopButton:
			ElevatorEmergency<-true
			driver.Elev_set_motor_direction(driver.MotorStop)
			driver.Elev_set_stop_lamp(true)
			//Reset the state, then send it.


		case <-elevatorEmergencyTimer.C:
			//System is idle
			if !*BetweenFloors && !*DoorState{
				elevatorEmergencyTimer.Reset(ElevatorEmergencyTimeInterval)

			}else{
					ElevatorEmergency<-true
			}

		}	

	}
}








//This function can probably be rewritten to half its length, if you just irierate over one segment twice, but change
//the direction for each iteration if no order in the current direction is found.
// Could also be rewritten to several small functions. 
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
				//log.Println("Order on floor confirmed")
				OrderComplete<-v
				delete(orders,k)
			}
			if *Direction == driver.Up{
				if v.Floor>*LastPassedFloor{
					orderOnNextFloors=true
					//log.Println("Order on the floor above")
				}
			}else{
				if v.Floor<*LastPassedFloor{
					orderOnNextFloors=true
					//log.Println("Order on the floor below")
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
					//log.Println("Order on floor confirmed")
				}
				if *Direction == driver.Up{
					if v.Floor>*LastPassedFloor{
						orderOnNextFloors=true
						//log.Println("Order on the floor above")
					}
				}else{
					if v.Floor<*LastPassedFloor{
						orderOnNextFloors=true
						//log.Println("Order on the floor below")
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


