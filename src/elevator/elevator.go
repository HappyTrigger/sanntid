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
	ElevatorEmergency <-chan bool)) {


	state := State_OnFloor
	Direction := driver.Down

	lastPassedFloor := driver.Elev_get_floor_sensor_signal()
	Orders := make(map[int]driver.OrderEvent)


	if lastPassedFloor == -1{
		log.Fatal("[FATAL]\tElevator initialized between floors")
	}

	for{

		select{

			case order:=<-NewOrder:
				Orders[order.OrderId]=order
				driver.Elev_set_button_lamp(order.Button,order.Floor,true)
				log.Println("New order in map")





			case sensor:=<-SensorEvent:
				if sensor !=-1{
					lastPassedFloor = sensor
					driver.Elev_set_floor_indicator(lastPassedFloor)
					state = State_OnFloor
				}


			case stop:=<-StopButton:
				log.Println("Stop butten has been pressed:",stop)
				state = State_Failiure


			default:
				//Do nothing 

			}

		switch state{

			case State_OnFloor:

				orderOnFloor,orderOnNextFloors:=OrderOnTheFloor(Orders,
					&Direction,
					lastPassedFloor)
				if orderOnFloor !=-1 {
					driver.Elev_set_motor_direction(driver.MotorStop)
					driver.Elev_set_door_open_lamp(true)

				}
				if orderOnNextFloors {
					
					
				}

				time.Sleep(1*time.Millisecond)





			case State_Moving:

				driver.Elev_set_motor_direction(Direction)

			case State_Failiure:
				//

			default:
				//
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
	log.Println("Direction:",*Direction)
	log.Println("OrderOnFloor: ",orderOnFloor)
	log.Println("orderOnNextFloors: ",orderOnNextFloors)

	return orderOnFloor,orderOnNextFloors
}