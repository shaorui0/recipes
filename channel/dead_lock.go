package main

import (
    "fmt"
    "sync"
    "time"
)

var mu sync.Mutex

func deadlock() {
    mu.Lock()
    fmt.Println("Main goroutine acquired the lock")
    defer mu.Unlock()

    go func() {
        fmt.Println("New goroutine trying to acquire the lock")
        mu.Lock() // Will block forever.
        fmt.Println("New goroutine acquired the lock")
        defer mu.Unlock()
    }()

    // Sleep to give the new goroutine time to run
    time.Sleep(2 * time.Second)
    fmt.Println("Main goroutine releasing the lock")
}

func main() {
    deadlock()
    // Sleep to keep the main function running
    time.Sleep(5 * time.Second)
}