package driver

import(
	//"log"
	//"os"
	"time"
)


const (
	N_FLOORS    = 4
	N_BUTTONS  = 3
	InvalidFloor = -1
	MOTOR_SPEED   = 2800
	PollInterval = 1 * time.Millisecond
)



func pollFloorSensor(sensorEventChan chan<- int) {
	state := -1

	for {
		sensorSignal := Elev_get_floor_sensor_signal()
        //log.Println("POlling")
		if state != sensorSignal {
			state = sensorSignal
			sensorEventChan <- state
		}
		time.Sleep(PollInterval)
	}
}

func pollButtons(order chan<- OrderEvent) {
	for{
		
		time.Sleep(3*time.Second)
		order <- OrderEvent{2, ButtonType(Up),0}
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
		time.Sleep(PollInterval)
	}
}

func Init(OrderEvent chan<- OrderEvent, sensorEventChan chan<- int, stopButtonChan chan<- bool) {
	go pollFloorSensor(sensorEventChan)
	go pollButtons(OrderEvent)
	go pollStopButton(stopButtonChan)
}







func Elev_set_motor_direction(Direction MotorDirection) {
    
}


func Elev_set_button_lamp(button ButtonType,floor int,value bool) {

}


func Elev_set_floor_indicator(floor int) {
   
}


func Elev_set_door_open_lamp(doorOpen bool) {

}


func Elev_set_stop_lamp(lightOn bool) {

}



func elev_get_button_signal(button ButtonType,floor int) bool {

    return false
}


func Elev_get_floor_sensor_signal() int {
    return 1
}


func  elev_get_stop_signal() bool {
    return false
}


func elev_get_obstruction_signal() int{
    return 1
}