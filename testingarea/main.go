package main

import (
	"fmt"
	"time" 
	
)
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

}