package main

//Pull from https://github.com/johhat/sanntidsheis if needed

import (	
	"./networking"
	"./utilities"
	"./manager"
	"./elevator"
	"fmt"
)

func main() {
	
	sendMsg := make(chan utilities.Message)
	recMsg := make(chan utilities.Message)



	
	go manager.Run(sendMsg, recMsg)
	go networking.Run(sendMsg,recMsg)
	go elevator.Run()

	for{
	}
	fmt.Println("Exiting program")

}