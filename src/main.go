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
	
	fromManager := make(chan utilities.Message)
	toManager := make(chan utilities.Message)
	connectionStatus := make(chan utilities.ConnectionStatus)
	//stateIsNew := make(chan bool)




	go manager.Run(fromManager, toManager, connectionStatus)
	go elevator.Run()
	go networking.Run(fromManager,toManager,connectionStatus)	
	for{
	}
	log.Println("Exiting program")

}