package main

import (
    "context"
    "log"
    "sync"
    "time"

    "google.golang.org/grpc"
    "metricspb"
)

// Connection pool to reuse gRPC connections.
type connPool struct {
    pool *sync.Pool
}

// NewConnPool initializes a connection pool with gRPC connections.
func NewConnPool(target string, size int) *connPool {
    p := &sync.Pool{
        New: func() interface{} {
            conn, err := grpc.Dial(target, grpc.WithInsecure())
            if err != nil {
                log.Fatalf("Failed to connect: %v", err)
            }
            return conn
        },
    }
    return &connPool{pool: p}
}

// Get retrieves a connection from the pool.
func (p *connPool) Get() *grpc.ClientConn {
    return p.pool.Get().(*grpc.ClientConn)
}

// Put returns a connection to the pool.
func (p *connPool) Put(conn *grpc.ClientConn) {
    p.pool.Put(conn)
}

func main() {
    pool := NewConnPool("localhost:50051", 5)
    conn := pool.Get()
    defer pool.Put(conn)

    client := metricspb.NewMetricsServiceClient(conn)
    stream, err := client.SendMetrics(context.Background())
    if err != nil {
        log.Fatalf("Failed to create stream: %v", err)
    }

    // Send 10 metrics to the server.
    for i := 0; i < 10; i++ {
        metric := &metricspb.Metric{
            Timestamp: time.Now().Unix(),
            Name:      "cpu_usage",
            Value:     float64(i) * 10.5,
        }
        if err := stream.Send(metric); err != nil {
            log.Fatalf("Failed to send metric: %v", err)
        }
        log.Printf("Sent Metric: %v", metric)
        time.Sleep(1 * time.Second) // Simulate interval.
    }

    res, err := stream.CloseAndRecv()
    if err != nil {
        log.Fatalf("Failed to receive response: %v", err)
    }
    log.Printf("Server Response: %v", res.Status)
}
