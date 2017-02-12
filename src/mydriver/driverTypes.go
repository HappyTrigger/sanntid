package driver


type ButtonType int  

const(
	Up ButtonType = iota
	Down
	Internal
)


type OrderEvent struct{
	Floor int
	Button ButtonType 
	OrderId int
}



type MotorDirection int

const (
	MotorUp MotorDirection = iota
	MotorStop
	MotorDown
)


