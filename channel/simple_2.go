package main

import (
    "fmt"
	"time"
	"sync"
)

func main() {
	// init wg
	// wg := make(sync.WaitGroup)
	var wg sync.WaitGroup

	// two goroutine
	wg.Add(2)

	// init channel
	ch := make(chan bool) // unbuffered channel
	// sender 
	go func(){
		defer wg.Done()
        for i := 0; i < 5; i++ {
			ch <- true
			time.Sleep(time.Second)
		}
		close(ch) // Close the channel when done sending
	}()
	
	
	// receiver
	go func(){
		defer wg.Done()
        for i := 0; i < 5; i++ {
			<- ch
			fmt.Println("received one message.")
		}
	}()

	// exit
	// time.Sleep(8 * time.Second)
	wg.Wait()
}