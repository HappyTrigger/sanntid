package elevator

import (
	".././utilities"
	"log"
	"time"
	".././mydriver"
)





func Init() {
	// get his position
	// "let the manager write the InternalOrders in the case where rebooting and
	//another manager gives back the memory of interalorders"
}



func Run(NewState chan<-utilities.State,
	NewOrder <-chan driver.OrderEvent,
	SensorEvent <-chan int,
	StopButton <-chan bool,
	DoorOpen <- chan bool,
	DoorClosed <-chan bool,
	ElevatorEmergency <-chan bool,
	orderComplete chan<-utilities.NewOrder) {



	Direction := driver.Down

	lastPassedFloor := driver.Elev_get_floor_sensor_signal()
	Orders := make(map[int]driver.OrderEvent)


	if lastPassedFloor == -1{
		log.Fatal("[FATAL]\tElevator initialized between floors")
	}

	var doorClose <-chan time.Time
	var State State
	State = State_idle
	doorClose = time.After(3*time.Second)

	for{

		select{

		case order:=<-NewOrder:
			Orders[order.OrderId] = order
			driver.Elev_set_button_lamp(order.Button,order.Floor,true)
			log.Println("New order in elevator")


			switch State {

				case State_idle:
				orderOnFloor,orderOnNextFloors:=OrderOnTheFloor(Orders,
				&Direction,
				lastPassedFloor)
				if orderOnFloor !=-1 {
					driver.Elev_set_motor_direction(driver.MotorStop)
					driver.Elev_set_door_open_lamp(true)
					doorClose = time.After(3*time.Second)
				}else{
					if orderOnNextFloors{
						if Direction==driver.Down{
							driver.Elev_set_motor_direction(driver.MotorDown)
						}else{
							driver.Elev_set_motor_direction(driver.MotorUp)
						}
						
					}
				}

								



				case State_moving:
					//Do nothing
				}

			


		case lastPassedFloor = <-SensorEvent:
			driver.Elev_set_floor_indicator(lastPassedFloor)
			if lastPassedFloor !=-1 {

				orderOnFloor,orderOnNextFloors:=OrderOnTheFloor(Orders,
				&Direction,
				lastPassedFloor)
			
				if orderOnFloor !=-1 {
					driver.Elev_set_motor_direction(driver.MotorStop)
					driver.Elev_set_door_open_lamp(true)
					doorClose = time.After(3*time.Second)
				}else{
					if orderOnNextFloors{
						if Direction==driver.Down{
							driver.Elev_set_motor_direction(driver.MotorDown)
						}else{
							driver.Elev_set_motor_direction(driver.MotorUp)
						}
						
					}
				}
			}



		case stop:=<-StopButton:
			log.Println("Stop butten has been pressed:",stop)



		case <-doorClose:


			driver.Elev_set_door_open_lamp(false)

			orderOnFloor, orderOnNextFloors := 
				OrderOnTheFloor(Orders, &Direction, lastPassedFloor)
			
			if orderOnFloor !=-1 {
				driver.Elev_set_motor_direction(driver.MotorStop)
				driver.Elev_set_door_open_lamp(true)
				doorClose = time.After(3*time.Second)
			} else {
				if orderOnNextFloors{
					if Direction == driver.Up{
						driver.Elev_set_motor_direction(driver.MotorUp)
					}else{
						driver.Elev_set_motor_direction(driver.MotorDown)
					}
					State=State_moving
				}else{

					State=State_idle
				}
			}
		}

	}
}

func OrderOnTheFloor(orders map[int]driver.OrderEvent,
	Direction *driver.ButtonType,currentFloor int)(int,bool){


	//log.Println("Order Checking func started")
	No_order := -1
	var orderOnNextFloors bool
	var orderOnFloor int 

	orderOnNextFloors = false
	orderOnFloor = -1

	for k,v:= range orders{
		//log.Println("Checking")
		if v.Button == *Direction || v.Button == driver.Internal{
			if currentFloor == v.Floor{
				orderOnFloor=k //Send orderNumber back
				driver.Elev_set_button_lamp(v.Button,v.Floor,false)
				log.Println("Order on floor")
				delete(orders,k)
			}
			
			if *Direction == driver.Up{
				if v.Floor>currentFloor{
					orderOnNextFloors=true
				}
			}else{
				if v.Floor<currentFloor{
					orderOnNextFloors=true
				}
			}
		}else{

			switch *Direction{
				case driver.Up:
					if v.Floor > currentFloor{
							orderOnNextFloors=true
					}

				case driver.Down:
					if v.Floor < currentFloor{
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
				if currentFloor == v.Floor{
					orderOnFloor=k //Send orderNumber back
					driver.Elev_set_button_lamp(v.Button,v.Floor,false)
					delete(orders,k)
					log.Println("Order on floor")
				}
				if *Direction == driver.Up{
					if v.Floor>currentFloor{
						orderOnNextFloors=true
						log.Println("Order Upwoards")
					}
				}else{
					if v.Floor<currentFloor{
						orderOnNextFloors=true
						log.Println("Order Downwards")
					}
				}
			}else{

			switch *Direction{
				case driver.Up:
					if v.Floor > currentFloor{
							orderOnNextFloors=true
					}

				case driver.Down:
					if v.Floor < currentFloor{
						orderOnNextFloors=true
							
					}
				}
			}
		}
	}
	return orderOnFloor,orderOnNextFloors
}