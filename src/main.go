package main

//Pull from https://github.com/johhat/sanntidsheis if needed

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