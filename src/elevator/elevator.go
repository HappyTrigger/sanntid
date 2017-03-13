package elevator

import (
	"log"
	"time"

	"../driver"
	"../utilities"
)

const (
	elevatorEmergencyTimeInterval = 6 * time.Second
	doorOpenTime                  = 3 * time.Second
)

var elevatorState utilities.State

func Run(
	NewOrder <-chan driver.OrderEvent,
	SensorEvent <-chan int,
	ElevatorEmergency chan<- bool,
	OrderComplete chan<- driver.OrderEvent,
	ElevatorStateToManager chan<- utilities.State,
	StopButton <-chan bool) {

	var doorClose <-chan time.Time
	elevatorEmergencyTimer := time.NewTimer(elevatorEmergencyTimeInterval)
	orders := make(map[int]driver.OrderEvent)

	elevatorState.DoorState = false
	elevatorState.Idle = true
	elevatorState.Direction = driver.Down
	elevatorState.LastRegisterdFloor = <-SensorEvent
	ElevatorStateToManager <- elevatorState

	if elevatorState.LastRegisterdFloor == -1 {
		log.Fatal("[FATAL]\tElevator initialized between floors")
	}

	for {
		select {

		case order := <-NewOrder:
			orders[order.Checksum] = order
			driver.Elev_set_button_lamp(order.Button, order.Floor, true)

			if elevatorState.Idle {

				elevatorControl(&doorClose, orders, OrderComplete)
				ElevatorStateToManager <- elevatorState

				if !elevatorEmergencyTimer.Stop() {
					<-elevatorEmergencyTimer.C
				}
				elevatorEmergencyTimer.Reset(elevatorEmergencyTimeInterval)
			}

		case elevatorState.LastRegisterdFloor = <-SensorEvent:

			driver.Elev_set_floor_indicator(elevatorState.LastRegisterdFloor)

			if elevatorState.LastRegisterdFloor != -1 {
				log.Println("On Floor ", elevatorState.LastRegisterdFloor+1)
				elevatorControl(&doorClose, orders, OrderComplete)

				ElevatorStateToManager <- elevatorState
				if !elevatorEmergencyTimer.Stop() {
					<-elevatorEmergencyTimer.C
				}
				elevatorEmergencyTimer.Reset(elevatorEmergencyTimeInterval)
			}

		case <-doorClose:
			driver.Elev_set_door_open_lamp(false)
			elevatorState.DoorState = false
			elevatorControl(&doorClose, orders, OrderComplete)

			ElevatorStateToManager <- elevatorState
			if !elevatorEmergencyTimer.Stop() {
				<-elevatorEmergencyTimer.C
			}
			elevatorEmergencyTimer.Reset(elevatorEmergencyTimeInterval)

		case <-StopButton:
			ElevatorEmergency <- true
			driver.Elev_set_motor_direction(driver.MotorStop)
			driver.Elev_set_stop_lamp(true)

		case <-elevatorEmergencyTimer.C:
			if elevatorState.Idle {
				elevatorEmergencyTimer.Reset(elevatorEmergencyTimeInterval)

			} else {
				ElevatorEmergency <- true
			}

		}

	}
}

func elevatorControl(
	doorClose *<-chan time.Time,
	Orders map[int]driver.OrderEvent,
	OrderComplete chan<- driver.OrderEvent) {

	orderOnFloor, orderOnNextFloors := findNextDestination(Orders, OrderComplete)

	if orderOnFloor {
		driver.Elev_set_motor_direction(driver.MotorStop)
		driver.Elev_set_door_open_lamp(true)
		*doorClose = time.After(doorOpenTime)
		elevatorState.DoorState = true
	} else {
		if orderOnNextFloors {
			if elevatorState.Direction == driver.Up {
				driver.Elev_set_motor_direction(driver.MotorUp)
			} else {
				driver.Elev_set_motor_direction(driver.MotorDown)
			}
			elevatorState.Idle = false
		} else {
			elevatorState.Idle = true

		}
	}
}

func findNextDestination(Orders map[int]driver.OrderEvent,
	OrderComplete chan<- driver.OrderEvent) (bool, bool) {

	orderOnNextFloors := false
	orderOnFloor := false

	for i := 0; i < 2; i++ {
		for checksum, order := range Orders {
			if order.Button == elevatorState.Direction || order.Button == driver.Internal {
				if elevatorState.LastRegisterdFloor == order.Floor {
					orderOnFloor = true
					driver.Elev_set_button_lamp(order.Button, order.Floor, false)

					OrderComplete <- order
					delete(Orders, checksum)
				}
				if elevatorState.Direction == driver.Up {
					if order.Floor > elevatorState.LastRegisterdFloor {
						orderOnNextFloors = true

					}
				} else {
					if order.Floor < elevatorState.LastRegisterdFloor {
						orderOnNextFloors = true
					}
				}
			} else {

				switch elevatorState.Direction {
				case driver.Up:
					if order.Floor > elevatorState.LastRegisterdFloor {
						orderOnNextFloors = true
					}

				case driver.Down:
					if order.Floor < elevatorState.LastRegisterdFloor {
						orderOnNextFloors = true

					}
				}
			}
		}

		if !orderOnFloor && !orderOnNextFloors {
			if elevatorState.Direction == driver.Up {
				elevatorState.Direction = driver.Down
			} else {
				elevatorState.Direction = driver.Up
			}
		}
	}
	return orderOnFloor, orderOnNextFloors
}
