package main

import (
    "fmt"
)

func main() {
    ch := make(chan int)

    go func() {
        for val := range ch {
            fmt.Println(val)
        }
        fmt.Println("Channel closed")
    }()

    // 发送一些数据
    ch <- 1
    ch <- 2
    ch <- 3

    // 关闭通道
    close(ch)

    // 等待 goroutine 处理完毕
    // 在实际应用中，通常会使用 sync.WaitGroup 或其他同步机制
    // 这里简单地使用一个 sleep 来等待
    time.Sleep(2 * time.Second)
}
