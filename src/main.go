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
	
	sendMsg := make(chan utilities.Message,100)
	recMsg := make(chan utilities.Message)
	connectionStatus := make(chan utilities.ConnectionStatus)
	//stateIsNew := make(chan bool)




	go manager.Run(sendMsg, recMsg, connectionStatus)
	go elevator.Run()
	go networking.Run(sendMsg,recMsg,connectionStatus)	
	for{
	}
	log.Println("Exiting program")

}