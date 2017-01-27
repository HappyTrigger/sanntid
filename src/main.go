package main

import (	
	"./networking"
	"fmt"
)

func main() {

	go networking.Run()
	for{
		
	}
	fmt.Println("Exiting program")

}