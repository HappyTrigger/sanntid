package elevator

/*
The elevator-module is responsible for everything regarding the elevator control and exectution.
It controls lighting, motorfunctioality, and keeps track of the current elevator-state.
It also sends state-updates to the manager when a state-change has been registerd.

*/
import (
	"log"
	"time"

	//"../dummydriver"
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
	elevatorState.LastRegisterdFloor = driver.Elev_get_floor_sensor_signal()

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
				log.Println("On Floor ", elevatorState.LastRegisterdFloor)
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

func findNextDestination(Orders map[int]driver.OrderEvent,
	OrderComplete chan<- driver.OrderEvent) (bool, bool) {

	orderOnNextFloors := false
	orderOnFloor := false

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
	}
	return orderOnFloor, orderOnNextFloors
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
				log.Println("Driving up")
			} else {
				driver.Elev_set_motor_direction(driver.MotorDown)
				log.Println("Driving Down")
			}
			elevatorState.Idle = false

		} else {
			elevatorState.Idle = true

		}
	}
}
