package main


import (
	"fmt"
	"time"
	"strconv"


)

type RawMessage struct {
	Data []byte
	Ip   string
}



func worker(done chan<-bool, c chan<-int, raw chan<-RawMessage) {
    fmt.Print("working...")
    time.Sleep(time.Second)
    fmt.Println("done")
    
    for i:=0; i<10; i++{
    	if i==9{
    		c<-i
    	}
    }
    msg := strconv.Itoa(52)
	buf := []byte(msg)
    raw <- RawMessage{Data: buf, Ip: "192.168.1.5"}

    // Send a value to notify that we're done.
    done <- true
}

func main() {

    // Start a worker goroutine, giving it the channel to
    // notify on.
    done := make(chan bool, 1)
    c := make(chan int)
    raw := make(chan RawMessage)
    go worker(done,c,raw)
    fmt.Println("The number is", <-c)
    //msg:=<-raw
    fmt.Println("The Ip is", <-raw)
    <- done
}








