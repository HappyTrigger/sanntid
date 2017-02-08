package main

//Pull from https://github.com/johhat/sanntidsheis if needed

import (	
	"./networking"
	"./utilities"
	"./manager"
	"./elevator"
	"log"
)

func main() {
	
	sendMsg := make(chan utilities.Message,50)
	recMsg := make(chan utilities.Message,50)



	
	go manager.Run(sendMsg, recMsg)
	go networking.Run(sendMsg,recMsg)
	go elevator.Run()

	for{
	}
	log.Println("Exiting program")

}