# pool

0. object pool
   1. 限制频繁创建新对象
        Limit the number of concurrent client connections.
        Use worker pools to process metrics asynchronously.
        Manage backpressure by blocking new connections if the pool is full.
1. conenction pool
   1. client connection pool
      1. 限制频繁创建新连接
      2. 一开始创建好，使用现成的connection
   2. server connection pool
      1. `Acquire a connection slot from the pool (blocks if no slots available).`, 限制接收的连接个数
      2. 创建好容量，利用slot 进行占据
      3. slot -> channel

---
## **1. Understanding the Purpose of Each Pool**

### **Connection Pool**

- **What it does**: Manages the number of **active connections**.
- **Purpose**:
    - Ensures the server doesn’t get overwhelmed by too many concurrent client connections.
    - Helps **limit resource consumption** (e.g., file descriptors, memory).
- **When to use**:
    - When you expect **high client concurrency** (many clients connecting simultaneously).

### **Goroutine Pool (Worker Pool)**

- **What it does**: Limits the **number of concurrent tasks** (i.e., goroutines) for processing data.
- **Purpose**:
    - Avoids **goroutine explosion**, where too many goroutines are created and overuse CPU/memory.
    - Helps keep the system **stable under heavy load** by reusing workers.
- **When to use**:
    - When each connection streams **many messages**, and **each message needs to be processed** by a worker thread (goroutine).

---

## **2. Which Pool Do You Need?**

### **Scenario Analysis: Client Streaming gRPC**

1. **High number of concurrent client connections**:
    - If **many clients** will open **simultaneous streaming connections**, you need a **connection pool** to **limit active connections**.
    - Without a connection pool, the server could get overwhelmed by too many open connections, causing resource exhaustion.
2. **High-frequency message streams from each client**:
    - If each **client connection streams many messages** (e.g., metrics every second), **goroutine reuse** (worker pool) is important to **avoid spawning a new goroutine for every message**.
    - This prevents **excessive context switching** and **memory pressure** from thousands of active goroutines.

---

### **Decision: When to Use One or Both**

| **Case** | **Use Connection Pool?** | **Use Goroutine Pool?** |
| --- | --- | --- |
| High **number of concurrent clients** | ✅ Yes | Optional (if message processing is lightweight) |
| High **frequency of streamed messages** per client | ❌ No (if clients are few) | ✅ Yes, to avoid goroutine explosion |
| Both **many clients** and **frequent messages** | ✅ Yes | ✅ Yes |

---

### **3. Recommended Approach for Your Case**

Since your scenario involves **client streaming of metrics**, where clients send frequent messages (e.g., **one message per second**) and potentially many clients may connect:

- **You need both** a **connection pool** and **goroutine pool**:
    - **Connection Pool**: Limits the number of simultaneous client connections.
    - **Goroutine Pool (Worker Pool)**: Ensures efficient reuse of goroutines for **processing each message** from the stream.

---

## **4. Implementation Strategy**

You can **combine both the connection pool and goroutine pool** as shown below:

- **Connection Pool**: Limits concurrent client connections.
- **Worker Pool**: Handles the processing of individual metrics efficiently.

## **2. Recommended Pool Sizes Based on Typical 5G Scenarios**

Given the scale of **5G metrics collection**, here are **ballpark numbers** for **pool sizes**:

| **Metric** | **Typical Number** | **Explanation** |
| --- | --- | --- |
| **Concurrent Connections (Connection Pool)** | **100–1,000** | Each gNodeB, IoT device, or base station may open a stream to send metrics. |
| **Goroutine Pool Size (Worker Pool)** | **2× to 5× the CPU core count** (e.g., 32–64) | Depends on the processing load. More workers help if the workload is CPU-bound. |
| **Message Queue Buffer Size** | **10,000+** messages | Buffer to handle bursts and avoid blocking if workers are busy. |
| **Message Frequency** | **1–10 messages per second per client** | Example: Metrics reporting at high frequency from each base station. |

---

## **5. Recommended Pool Sizes Based on Typical 5G Scenarios**

Given the scale of **5G metrics collection**, here are **ballpark numbers** for **pool sizes**:

| **Metric** | **Typical Number** | **Explanation** |
| --- | --- | --- |
| **Concurrent Connections (Connection Pool)** | **100–1,000** | Each gNodeB, IoT device, or base station may open a stream to send metrics. |
| **Goroutine Pool Size (Worker Pool)** | **2× to 5× the CPU core count** (e.g., 32–64) | Depends on the processing load. More workers help if the workload is CPU-bound. |
| **Message Queue Buffer Size** | **10,000+** messages | Buffer to handle bursts and avoid blocking if workers are busy. |
| **Message Frequency** | **1–10 messages per second per client** | Example: Metrics reporting at high frequency from each base station. |