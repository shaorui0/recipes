package main

import (
    "fmt"
	"time"
)

func main() {

	// init channel
	ch := make(chan bool) // unbuffered channel
	// sender 
	go func(){
        for i := 0; i < 5; i++ {
			ch <- true
			time.Sleep(time.Second)
		}
		close(ch) // Close the channel when done sending
	}()
	
	
	// receiver
	go func(){
        for i := 0; i < 5; i++ {
			<- ch
			fmt.Println("received one message.")
		}
	}()

	// exit
	time.Sleep(8 * time.Second)
}