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

// Object pool for Metric objects to reduce allocations.
var metricPool = sync.Pool{
    New: func() interface{} {
        return &metricspb.Metric{}
    },
}

// Worker pool for goroutine reuse.
type workerPool struct {
    jobs    chan *metricspb.Metric
    workers int
}

// NewWorkerPool initializes the worker pool.
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

// worker processes metrics from the job channel.
func (wp *workerPool) worker() {
    for metric := range wp.jobs { // waiting for job from channel
        processMetric(metric)
        metricPool.Put(metric) // Return the metric object to the pool.
    }
}

// processMetric simulates processing of a metric.
func processMetric(metric *metricspb.Metric) {
    fmt.Printf("Processing Metric: %s, Value: %.2f, Timestamp: %d\n",
        metric.Name, metric.Value, metric.Timestamp)
    time.Sleep(500 * time.Millisecond) // Simulate processing delay.
}

// MetricsServiceServer implementation with worker pool.
type MetricsServer struct {
    metricspb.UnimplementedMetricsServiceServer
    pool *workerPool
}

// SendMetrics receives metrics from the stream and enqueues them in the worker pool.
func (s *MetricsServer) SendMetrics(stream metricspb.MetricsService_SendMetricsServer) error {
    for {
        metric := metricPool.Get().(*metricspb.Metric)
        err := stream.RecvMsg(metric)
        if err != nil {
            if err.Error() == "EOF" {
                break
            }
            return err
        }
        s.pool.jobs <- metric // Send the metric to the worker pool.
    }
    return stream.SendAndClose(&metricspb.UploadResponse{Status: "Metrics Processed"})
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    server := &MetricsServer{pool: NewWorkerPool(5)} // Initialize with 5 workers.
    metricspb.RegisterMetricsServiceServer(grpcServer, server)

    log.Println("gRPC Server listening on :50051")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
