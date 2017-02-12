package elevator






type State int

const(
	State_Init State = iota
	State_OnFloor
	State_Moving
	State_Failiure
)

/*
const(
	Dir_up = iota 
	Dir_down
	Internal
	
)

*/
