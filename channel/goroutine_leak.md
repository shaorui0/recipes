## **2. Common Causes of Goroutine Leaks**

Here are the most frequent reasons why goroutines **fail to exit gracefully**, causing leaks:

### **2.1 Blocking on Channel Operations**

- If a goroutine **waits indefinitely** on a channel that no one is writing to, it will never terminate.

### Example:

```go
func leakExample() {
    ch := make(chan int) // No writer on the channel.
    go func() {
        <-ch // Blocks forever.
    }()
}
```

- **Fix**: Ensure the channel is **closed or written to** in all paths, or use **timeouts** with `select`.

---

### **2.2 Forgotten Goroutines in Background Tasks**

- Goroutines spawned without any way to stop them may continue running indefinitely.

### Example:

```go
func startLeakyTask() {
    go func() {
        for {
            // Do some work.
            time.Sleep(1 * time.Second) // Keeps running forever.
        }
    }()
}
```

- **Fix**: Use **context cancellation** to control when the goroutine should exit.

---

### **2.3 Deadlocks (Blocked on Mutexes or Channels)**

- If two or more goroutines **block waiting on each other** (e.g., through channels or mutexes), they may never exit.

### Example:

```go
var mu sync.Mutex

func deadlock() {
    mu.Lock()
    defer mu.Unlock()
    go func() {
        mu.Lock() // Will block forever.
        defer mu.Unlock()
    }()
}
```

- **Fix**: Ensure that **locks and channels** are always released properly to avoid deadlocks.

---

### **2.4 Missing Exit Conditions in `select` Statements**

- Goroutines using `select` without a proper **exit condition** may block forever.

### Example:

```go
func leakInSelect() {
    ch := make(chan int)
    go func() {
        select {
        case val := <-ch:
            fmt.Println(val)
        // No way to exit the select block.
        }
    }()
}
```

- **Fix**: Add **timeouts or cancellation signals** to ensure the goroutine exits when appropriate.

---

### **2.5 Waiting on Network or I/O Operations Indefinitely**

- If a goroutine is **waiting for a network call or I/O operation** that never completes, it can leak.

### Example:

```go
func leakyNetworkCall() {
    go func() {
        _, err := http.Get("http://example.com") // Network may hang.
        if err != nil {
            fmt.Println("Request failed")
        }
    }()
}
```

- **Fix**: Use **timeouts** for network or I/O operations.