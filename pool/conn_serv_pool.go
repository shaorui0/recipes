package main

import (
    "fmt"
    "log"
    "net"
    "sync"
    "time"

    "google.golang.org/grpc"
    "metricspb"
)

// Connection pool to manage active gRPC streams (limited pool size).
type connectionPool struct {
    pool     chan struct{} // Channel to limit active connections.
    poolSize int
}

// NewConnectionPool initializes the connection pool with a given size.
func NewConnectionPool(size int) *connectionPool {
    return &connectionPool{
        pool:     make(chan struct{}, size),
        poolSize: size,
    }
}

// Acquire reserves a slot in the pool for a new connection.
func (p *connectionPool) Acquire() {
    p.pool <- struct{}{} // Block if pool is full.
}

// Release frees up a slot in the pool.
func (p *connectionPool) Release() {
    <-p.pool
}

// Worker pool to reuse goroutines for processing metrics.
type workerPool struct {
    jobs    chan *metricspb.Metric
    workers int
}

// NewWorkerPool initializes the worker pool with a given number of workers.
func NewWorkerPool(workers int) *workerPool {
    wp := &workerPool{
        jobs:    make(chan *metricspb.Metric, 100),
        workers: workers,
    }
    for i := 0; i < workers; i++ {
        go wp.worker()
    }
    return wp
}

// worker processes metrics from the job queue.
func (wp *workerPool) worker() {
    for metric := range wp.jobs {
        processMetric(metric)
    }
}

// Metrics server implementation with a connection and worker pool.
type MetricsServer struct {
    metricspb.UnimplementedMetricsServiceServer
    connPool *connectionPool
    workerPool *workerPool
}

// SendMetrics handles client streaming by sending metrics to the worker pool.
func (s *MetricsServer) SendMetrics(stream metricspb.MetricsService_SendMetricsServer) error {
    // Acquire a connection slot from the pool (blocks if no slots available).
    s.connPool.Acquire()
    defer s.connPool.Release() // Release the connection slot when done.

    for {
        metric, err := stream.Recv()
        if err != nil {
            if err.Error() == "EOF" {
                break
            }
            return err
        }
        // Send the metric to the worker pool for processing.
        s.workerPool.jobs <- metric
    }

    return stream.SendAndClose(&metricspb.UploadResponse{
        Status: "Metrics Processed",
    })
}

// processMetric simulates the processing of a metric.
func processMetric(metric *metricspb.Metric) {
    fmt.Printf("Processing Metric: %s, Value: %.2f, Timestamp: %d\n",
        metric.Name, metric.Value, metric.Timestamp)
    time.Sleep(500 * time.Millisecond) // Simulate processing delay.
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()

    // Initialize connection pool and worker pool.
    connPool := NewConnectionPool(10) // Limit to 10 concurrent connections.
    workerPool := NewWorkerPool(5)    // 5 concurrent workers for processing.

    server := &MetricsServer{
        connPool:  connPool,
        workerPool: workerPool,
    }

    metricspb.RegisterMetricsServiceServer(grpcServer, server)

    log.Println("gRPC Server listening on :50051")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
