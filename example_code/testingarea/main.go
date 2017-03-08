package main

import (
	"fmt"
	//"time" 
	
)
/*
type ButtonType int

type Vertex struct {
	X ButtonType
	Y ButtonType
}

const(
	Up ButtonType = iota
	Down
	)

func another(y* ButtonType){
	*y = 100	
}

func main() {
	v := Vertex{1, 2}
	p := &v
	p.X = 1e9
	fmt.Println(v)
	
	y:= &v.Y
	*y=10
	*y = Up
	fmt.Println(v.Y)

	another(y)
	
	fmt.Println(*y)
	
	fmt.Println(v.Y)
	
	var b ButtonType
	b=5
	another(&b)
	fmt.Println(b)
	
	ne := make(<-chan time.Time)
	cha := make(<-chan time.Time)
	random2(&cha)
	Loop:
	for{
		select{
		case <-cha:
			fmt.Println("Cha")
			break Loop
		case <-ne:
			fmt.Println("Ting skjer")
		default:
		
		}
	}
}


func random2(cha * <-chan time.Time){
	*cha = time.After(1*time.Second)

}*/

var localIp string

func orderdelegate(
	orderAssignedToMap map[int]string,
	stateMap map[string]int,
	currentPeers []string,
	fitnessMap map[string]int) (bool, map[int]string) {


	for elevator,state := range stateMap{
		for _,peer:= range currentPeers{
			if elevator == peer{
				
				fitnessMap[elevator]=state
			}
		}

	}


	var maxValue int
	var OrderGivenToIp string
	for ip,value := range fitnessMap{
		if value > maxValue{
			maxValue=value
			OrderGivenToIp = ip
		}
	}
	orderAssignedToMap[maxValue]=OrderGivenToIp

	fmt.Println("Order given to: ",OrderGivenToIp, " Max value is :", maxValue)
	if OrderGivenToIp==localIp{
		return true,orderAssignedToMap
	}else{
		return false,orderAssignedToMap
	}
}

func main() {
	//var fitness int
	localIp = "129.21.23.431"
	otherIp := "254.21.23.4444444"
	fitnessMap := make(map[string]int)
	stateMap := make(map[string]int)
	var currentPeers []string
	orderAssignedToMap := make(map[int]string)

	stateMap[localIp]=25
	stateMap[otherIp]=20
	currentPeers=append(currentPeers,localIp,otherIp)
	

	if ok,orderAssignedToMap:=orderdelegate(orderAssignedToMap,stateMap,currentPeers,fitnessMap); ok{
		fmt.Println("Order taken")
		for k,_ := range orderAssignedToMap{
			fmt.Println("By :", k)
		}
	}else{
		fmt.Println("Order was given to another elevator")
	}


}




