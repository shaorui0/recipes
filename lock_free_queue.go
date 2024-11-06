package main

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

// Node represents a node in the queue.
type Node struct {
	value interface{}
	next  *Node
}

// LockFreeQueue represents a lock-free queue.
type LockFreeQueue struct {
	head *Node
	tail *Node
}

// NewLockFreeQueue creates a new lock-free queue.
func NewLockFreeQueue() *LockFreeQueue {
	node := &Node{}
	return &LockFreeQueue{
		head: node,
		tail: node,
	}
}

// Enqueue adds an element to the end of the queue.
func (q *LockFreeQueue) Enqueue(value interface{}) {
	node := &Node{value: value}
	for {
		tail := q.tail
		next := tail.next
		if tail == q.tail { // Check if tail is consistent
			if next == nil {
				// Try to link the new node at the end of the list
				if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&tail.next)), unsafe.Pointer(next), unsafe.Pointer(node)) {
					// Enqueue is done. Try to swing the tail to the inserted node
					atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail)), unsafe.Pointer(tail), unsafe.Pointer(node))
					return
				}
			} else {
				// Tail is falling behind. Try to advance it
				atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail)), unsafe.Pointer(tail), unsafe.Pointer(next))
			}
		}
	}
}

// Dequeue removes and returns an element from the front of the queue.
func (q *LockFreeQueue) Dequeue() (interface{}, bool) {
	for {
		head := q.head
		tail := q.tail
		next := head.next
		if head == q.head { // Check if head is consistent
			if head == tail {
				if next == nil {
					// Queue is empty
					return nil, false
				}
				// Tail is falling behind. Try to advance it
				atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail)), unsafe.Pointer(tail), unsafe.Pointer(next))
			} else {
				// Read value before CAS, otherwise another dequeue might free the next node
				value := next.value
				// Try to swing the head to the next node
				if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head)), unsafe.Pointer(head), unsafe.Pointer(next)) {
					return value, true
				}
			}
		}
	}
}

func main() {
	queue := NewLockFreeQueue()

	// Enqueue elements
	queue.Enqueue(1)
	queue.Enqueue(2)
	queue.Enqueue(3)

	// Dequeue elements
	if value, ok := queue.Dequeue(); ok {
		fmt.Println("Dequeued:", value)
	}
	if value, ok := queue.Dequeue(); ok {
		fmt.Println("Dequeued:", value)
	}
	if value, ok := queue.Dequeue(); ok {
		fmt.Println("Dequeued:", value)
	}
	if value, ok := queue.Dequeue(); ok {
		fmt.Println("Dequeued:", value)
	} else {
		fmt.Println("Queue is empty")
	}
}