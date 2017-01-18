package main

import (
	. "fmt"
	. "runtime"
	. "sync"
)

func adder(step int, ch chan int, wg *WaitGroup) {
	for x := 0; x < 100; x++ {
		i := <-ch
		i += step
		ch <- i
		Println(i)
	}
	wg.Done()
}

func main() {
	GOMAXPROCS(NumCPU()) //limits the number of threads that can execute user-level commands
	ch := make(chan int, 1)
	ch <- 0

//	m := 0 // := fordi den blir definert for første gang
//	m = 3
	var wg WaitGroup
	wg.Add(2)

	go adder(1, ch, &wg)
	go adder(-1, ch, &wg)

	wg.Wait()

	i := <-ch
	Println(i)
}


	go func(){
		for 1..10000
			addCh <- true
		done <- true
	}()

	go func(){
		for 1..10000
			subCh <- true
		done <- true
	}()
	<-done
	<-done
	fmt.Println("done!: " <-getCh)



func numberServer(addCh chan bool, subCh chan bool, getCh chan int) {}
	i := 0
	for {
		select {
			<- addCh:
				i++
			<- subCh:
				i--
			getCh <- i:
		}
	}
}