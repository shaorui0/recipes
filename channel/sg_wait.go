package main

import (
    "fmt"
    "sync"
    "time"
)

var wg sync.WaitGroup // global

// task 是一个简单的任务函数
func task(id int) {
    defer wg.Done() // 确保在函数结束时调用 Done
    fmt.Printf("Task %d is starting\n", id)
    time.Sleep(2 * time.Second) // 模拟一些工作
    fmt.Printf("Task %d is done\n", id)
}

func main() {
    numTasks := 3 // 要启动的任务数量

    for i := 1; i <= numTasks; i++ {
        wg.Add(1) // 增加 WaitGroup 计数
        go task(i) // 启动一个新的 goroutine
    }

    wg.Wait() // 等待所有 goroutine 完成
    fmt.Println("All tasks are done")
}