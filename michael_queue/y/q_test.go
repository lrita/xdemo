package y

import (
	"sync/atomic"
	"testing"
)

func TestQueue(t *testing.T) {
	const max = 100
	q := NewQueue()
	go func() {
		for i := 0; i < max; i++ {
			q.Enqueue(i)
		}
	}()

	j := 0
	for {
		v := q.Dequeue()
		if v == nil {
			continue
		}
		if v.(int) != j {
			t.Fatalf("got(%v) != want(%v)", v, j)
			return
		} else {
			j++
			if j == max {
				break
			}
		}
	}
}

func BenchmarkQueue(b *testing.B) {
	i := uint32(0)
	q := NewQueue()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		x := atomic.AddUint32(&i, 1)
		if x%2 == 0 {
			for pb.Next() {
				q.Enqueue(pb)
			}
		} else {
			for pb.Next() {
				q.Dequeue()
			}
		}
	})
}
