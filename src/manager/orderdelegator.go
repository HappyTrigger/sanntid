package manager

import (
	"math"

	"../driver"
	"../utilities"
)

func increaseFitnessPerOrder(peer string,
	orderAssignedToElevator map[int]string) float64 {
	var fitness float64

	for _, elevatorId := range orderAssignedToElevator {
		if elevatorId == peer {
			fitness += float64(doorOpenTime)
		}
	}

	return fitness
}

func orderDelegator(StateMap map[string]utilities.State,
	OrderEvent driver.OrderEvent, currentPeers []string, orderAssignedToMap map[int]string) bool {

	fitnessMap := make(map[string]float64)

	for elevator, state := range StateMap {
		for _, peer := range currentPeers {
			if elevator == peer {
				if state.Idle {
					fitnessMap[elevator] = math.Abs(float64(state.LastRegisterdFloor-OrderEvent.Floor)) * driver.TravelTimeBetweenFloors
				} else {
					floorDifference := float64(OrderEvent.Floor - state.LastRegisterdFloor)
					switch OrderEvent.Button {
					case driver.Up:
						if state.Direction == driver.Up {
							if floorDifference >= 0 { //order Up above the elevator, and elevator moving up
								fitnessMap[elevator] += float64(floorDifference)
							} else {
								//order Up below, and elevator moving up
								fitnessMap[elevator] += float64(math.Abs(float64(floorDifference)) + float64((driver.N_FLOORS-state.LastRegisterdFloor)*2))
							}

						} else {
							//Order Up bellow, and eleavtor moving down
							if floorDifference <= 0 {
								fitnessMap[elevator] += float64(state.LastRegisterdFloor*2) - math.Abs(float64(floorDifference))
							} else {
								//Order Up above and eleavtor moving down
								fitnessMap[elevator] += float64(math.Abs(float64(floorDifference)) + float64(state.LastRegisterdFloor*2))
							}
						}

					case driver.Down:
						if state.Direction == driver.Down {

							if floorDifference <= 0 { // Order downwards bellow, and elevator moving down
								fitnessMap[elevator] += math.Abs(float64(floorDifference))
							} else {
								//Order downwards above, and elevator moving down
								fitnessMap[elevator] += float64(math.Abs(float64(floorDifference)) + float64(state.LastRegisterdFloor*2))
							}
						} else {
							//Order downwards above and elevator moving up
							if floorDifference >= 0 {
								fitnessMap[elevator] += float64((driver.N_FLOORS-state.LastRegisterdFloor)*2) - float64(floorDifference)
							} else {
								//Order downwards bellow and elevator moving up
								fitnessMap[elevator] += float64(math.Abs(float64(floorDifference)) + float64((driver.N_FLOORS-state.LastRegisterdFloor)*2))
							}
						}
					}
					fitnessMap[elevator] = fitnessMap[elevator] * driver.TravelTimeBetweenFloors
					fitnessMap[elevator] += increaseFitnessPerOrder(peer, orderAssignedToMap)
				}
			}
		}
	}
	var minFitness float64
	minFitness = 100
	var currentId string
	for elevatorId, fitness := range fitnessMap {
		if fitness == minFitness {
			if elevatorId > currentId {
				currentId = elevatorId
			}
		}
		if fitness < minFitness {
			minFitness = fitness
			currentId = elevatorId
		}
	}

	orderAssignedToMap[OrderEvent.Checksum] = currentId
	if currentId == localId {
		return true
	} else {
		return false
	}
}
