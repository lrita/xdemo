// concurrent lock free queue by Maged M. Michael and Michael L. Scott
package x

import (
	"sync/atomic"
	"unsafe"
)

type Node struct {
	Next  *Node
	Value interface{}
}

type Queue struct {
	Head *Node
	Tail *Node
}

func NewQueue() *Queue {
	n := &Node{}
	return &Queue{
		Head: n,
		Tail: n,
	}
}

func (q *Queue) EnqueueNotOptimized(v interface{}) {
	ok := false
	node := &Node{Next: nil, Value: v}
	for !ok {
		tail := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.Tail)))) // acquire-load
		ok = atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&tail.Next)), nil, unsafe.Pointer(node))
		atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.Tail)),
			unsafe.Pointer(tail), atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&tail.Next))))
	}
}

func (q *Queue) DequeueNotOptimized() (v interface{}) {
	ok := false
	for !ok {
		head := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.Head))))    // acquire-load
		tail := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.Tail))))    // acquire-load
		next := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&head.Next)))) // acquire-load
		if next == nil {
			return nil
		}
		ok = atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.Head)),
			unsafe.Pointer(head), unsafe.Pointer(next))
		if ok {
			v = next.Value
		}
		if head == tail {
			atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.Tail)),
				unsafe.Pointer(tail), atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&tail.Next))))
		}
	}
	return v
}

func (q *Queue) Enqueue(v interface{}) {
	ok := false
	node := &Node{Next: nil, Value: v}
	for !ok {
		tail := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.Tail))))    // acquire-load
		next := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&tail.Next)))) // acquire-load
		if tail == (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.Tail)))) {
			if next == nil {
				ok = atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&tail.Next)), nil, unsafe.Pointer(node))
			}
			atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.Tail)),
				unsafe.Pointer(tail), atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&tail.Next))))
		}
	}
}

func (q *Queue) Dequeue() (v interface{}) {
	ok := false
	for !ok {
		head := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.Head))))    // acquire-load
		tail := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.Tail))))    // acquire-load
		next := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&head.Next)))) // acquire-load
		if head == (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.Head)))) {
			if head == tail {
				if next == nil {
					return nil
				}
				atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.Tail)),
					unsafe.Pointer(tail), (unsafe.Pointer(next)))
			} else {
				ok = atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.Head)),
					unsafe.Pointer(head), unsafe.Pointer(next))
				if ok {
					v = next.Value
				}
			}
		}
	}
	return v
}
