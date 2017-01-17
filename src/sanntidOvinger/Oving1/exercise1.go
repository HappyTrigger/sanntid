package main

import(
	"fmt"
	"time"


)

var global_time = 0

func thread_1(){
	for j:=0; j<1000;j++{
		global_time = global_time+1
	}
}
func thread_2(){
	for k:=0; k<1000;k++{
		global_time = global_time-1
	}
}


func main(){
	go thread_1()
	go thread_2()
	time.Sleep(time.Second*3)
	

	fmt.Println(global_time)
}