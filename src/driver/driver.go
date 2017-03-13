package driver

import (
	"log"
	"os"
	"time"

	"./io"
)

const (
	N_FLOORS                  = 4
	N_BUTTONS                 = 3
	InvalidFloor              = -1
	MOTOR_SPEED               = 2800
	pollInterval              = 1 * time.Millisecond
	TRAVELTIME_BETWEEN_FLOORS = 2.7 //Seconds
)

var lamp_channel_matrix = [N_FLOORS][N_BUTTONS]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var button_channel_matrix = [N_FLOORS][N_BUTTONS]int{
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}

func pollFloorSensor(sensorEventChan chan<- int) {
	state := -1

	for {
		sensorSignal := elev_get_floor_sensor_signal()
		if state != sensorSignal {
			state = sensorSignal
			sensorEventChan <- state
		}
		time.Sleep(pollInterval)
	}
}

func pollButtons(order chan<- OrderEvent) {

	var isPressed [N_BUTTONS][N_FLOORS]bool

	for {
		for f := 0; f < N_FLOORS; f++ {
			for button := 0; button < N_BUTTONS; button++ {
				if isPressed[ButtonType(button)][f] != elev_get_button_signal(ButtonType(button), f) {
					isPressed[ButtonType(button)][f] = !isPressed[ButtonType(button)][f]
					if isPressed[ButtonType(button)][f] {
						order <- OrderEvent{f, ButtonType(button), 0}
					}
				}
			}
		}
		time.Sleep(pollInterval)
	}
}

func pollStopButton(stopButtonChan chan<- bool) {
	isPressed := elev_get_stop_signal()

	for {
		if isPressed != elev_get_stop_signal() {
			isPressed = !isPressed

			if isPressed {
				stopButtonChan <- true
			}
		}
		time.Sleep(pollInterval)
	}
}

func Init(OrderEvent chan<- OrderEvent, sensorEventChan chan<- int, stopButtonChan chan<- bool) {
	go pollFloorSensor(sensorEventChan)
	go pollButtons(OrderEvent)
	go pollStopButton(stopButtonChan)
}

func init() {
	err := io.Init()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	for f := 0; f < N_FLOORS; f++ {
		for b := 0; b < N_BUTTONS; b++ {
			Elev_set_button_lamp(ButtonType(b), f, false)
		}
	}

	Elev_set_stop_lamp(false)
	Elev_set_door_open_lamp(false)
	Elev_set_motor_direction(MotorDown)

	timeout := time.After(10 * time.Second)

	for elev_get_floor_sensor_signal() == InvalidFloor {
		select {
		case <-timeout:
			log.Fatal("Timeout in driver. Did not get to valid floor in time.")
			os.Exit(1)
		default:
		}
	}

	Elev_set_motor_direction(MotorStop)
	Elev_set_floor_indicator(elev_get_floor_sensor_signal())

}

func Elev_set_motor_direction(Direction MotorDirection) {

	switch Direction {

	case MotorUp:
		io.ClearBit(MOTORDIR)
		io.WriteAnalog(MOTOR, MOTOR_SPEED)
	case MotorDown:
		io.SetBit(MOTORDIR)
		io.WriteAnalog(MOTOR, MOTOR_SPEED)
	case MotorStop:
		io.WriteAnalog(MOTOR, 0)
	}
}

func Elev_set_button_lamp(button ButtonType, floor int, value bool) {

	if floor >= 0 && floor < N_FLOORS {
		if button >= 0 && button < N_BUTTONS {
			switch button {
			case Up, Down, Internal:
				if value {
					io.SetBit(lamp_channel_matrix[floor][button])
				} else {
					io.ClearBit(lamp_channel_matrix[floor][button])
				}
			}
		}
	}
}

func Elev_set_floor_indicator(floor int) {
	if floor >= 0 && floor < N_FLOORS {
		if floor&0x02 != 0 {
			io.SetBit(LIGHT_FLOOR_IND1)
		} else {
			io.ClearBit(LIGHT_FLOOR_IND1)
		}

		if floor&0x01 != 0 {
			io.SetBit(LIGHT_FLOOR_IND2)
		} else {
			io.ClearBit(LIGHT_FLOOR_IND2)
		}

		// Binary encoding. One light must always be on.
		if floor&0x02 != 0 {
			io.SetBit(LIGHT_FLOOR_IND1)
		} else {
			io.ClearBit(LIGHT_FLOOR_IND1)
		}

		if floor&0x01 != 0 {
			io.SetBit(LIGHT_FLOOR_IND2)
		} else {
			io.ClearBit(LIGHT_FLOOR_IND2)
		}
	}
}

func Elev_set_door_open_lamp(doorOpen bool) {
	if doorOpen {
		io.SetBit(LIGHT_DOOR_OPEN)
	} else {
		io.ClearBit(LIGHT_DOOR_OPEN)
	}
}

func Elev_set_stop_lamp(lightOn bool) {
	if lightOn {
		io.SetBit(LIGHT_STOP)
	} else {
		io.ClearBit(LIGHT_STOP)
	}
}

func elev_get_button_signal(button ButtonType, floor int) bool {

	if floor >= 0 && floor < N_FLOORS {
		if button >= 0 && button < N_BUTTONS {
			return io.ReadBit(button_channel_matrix[floor][button]) != 0
		}
	} else {
		log.Println("Tried to get button or floor out of bounds")
	}
	return false
}

func elev_get_floor_sensor_signal() int {
	switch {
	case io.ReadBit(SENSOR_FLOOR1) != 0:
		return 0
	case io.ReadBit(SENSOR_FLOOR2) != 0:
		return 1
	case io.ReadBit(SENSOR_FLOOR3) != 0:
		return 2
	case io.ReadBit(SENSOR_FLOOR4) != 0:
		return 3

	default:
		return -1
	}
}

func elev_get_stop_signal() bool {
	return io.ReadBit(STOP) == 1
}

func elev_get_obstruction_signal() int {
	return io.ReadBit(OBSTRUCTION)
}
