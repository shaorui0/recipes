package main

import (
    "context"
    "fmt"
    "time"
)

// startTaskWithContext 启动一个带有上下文的任务
func startTaskWithContext(ctx context.Context) {
    go func() {
        for {
            select {
            case <-ctx.Done():
                fmt.Println("Goroutine exiting")
                return
            default:
                // 执行一些工作
                fmt.Println("Working...")
                time.Sleep(1 * time.Second)
            }
        }
    }()
}

func main() {
    // 创建一个带有取消功能的上下文
    ctx, cancel := context.WithCancel(context.Background())

    // 启动带有上下文的任务
    startTaskWithContext(ctx)

    // 让任务运行一段时间
    time.Sleep(5 * time.Second)

    // 取消上下文，通知 goroutine 退出
    cancel()

    // 等待一段时间以确保 goroutine 退出
    time.Sleep(1 * time.Second)

    fmt.Println("Main function exiting")
}