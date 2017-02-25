package main

import(
	"fmt"
	"time"
	)



func main() {
	c := make(<-chan time.Time)
	c2 := make(<-chan time.Time)
	fmt.Println("Starting")

	i := 0

	state := false

	if !state {
		fmt.Println("Staring")
	}

	somefunc(&c2)
	somefunc2(&c)
	Loop:
	for{
		select{
		case time2:=<-c:
			fmt.Println("Timeout",time2)
			break Loop

		case <-c2:
			fmt.Println("C2 triggerd")


		}
	}
	fmt.Println("Int: ", i)

}

func somefunc(c *<-chan time.Time) {
	*c = time.After(1*time.Second)
}

func somefunc2(c *<-chan time.Time) {
	*c = time.After(2*time.Second)
}