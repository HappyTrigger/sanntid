package elevator






type State int

const(
	State_idle = iota
	State_moving

)

type DoorState bool

const(
	DoorOpen DoorState = true
	DoorClosed = false
)

/*
const(
	Dir_up = iota 
	Dir_down
	Internal
	
)

*/
