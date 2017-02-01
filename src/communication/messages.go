package Communication


type NewOrder struct{
    Floor int 
	Direction int 
}

type ElevatorState struct{
	TravelingDirection int 
	InternalOrders[]int
	CurrentFloor int 
}

