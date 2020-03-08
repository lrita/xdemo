// concurrent queue using mutex by Maged M. Michael and Michael L. Scott
package y

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type Node struct {
	Next  *Node
	Value interface{}
}

type Queue struct {
	H_lock sync.Mutex
	T_lock sync.Mutex
	Head   *Node
	Tail   *Node
}

func NewQueue() *Queue {
	n := &Node{}
	return &Queue{
		Head: n,
		Tail: n,
	}
}

func (q *Queue) Enqueue(v interface{}) {
	n := &Node{Next: nil, Value: v}
	q.T_lock.Lock()
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&q.Tail.Next)), unsafe.Pointer(n)) // release-store
	q.Tail = n
	q.T_lock.Unlock()
}

func (q *Queue) Dequeue() interface{} {
	q.H_lock.Lock()
	node := (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.Head.Next)))) // acquire-load
	if node == nil {
		q.H_lock.Unlock()
		return nil
	}
	v := node.Value
	q.Head = node
	q.H_lock.Unlock()
	return v
}
