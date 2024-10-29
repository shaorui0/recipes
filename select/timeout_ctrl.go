package main

import (
    "fmt"
    "time"
)

func main() {
    ch := make(chan int)

    go func() {
        time.Sleep(2 * time.Second)
        ch <- 1
    }()

    select {
    case res := <-ch:
        fmt.Println("Received:", res)
    case <-time.After(1 * time.Second):
        fmt.Println("Timeout")
    }
}