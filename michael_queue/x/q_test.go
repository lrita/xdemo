package x

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestQueueNotOptimized(t *testing.T) {
	const max = 1000
	const c = 10
	q := NewQueue()

	for i := 0; i < c; i++ {
		go func() {
			time.Sleep(time.Millisecond)
			for i := 0; i < max; i++ {
				q.EnqueueNotOptimized(i)
			}
		}()
	}

	var (
		lock sync.Mutex
		wg   sync.WaitGroup
		j    int32
		jmap = make(map[int]int)
	)

	for i := 0; i < c; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				v := q.DequeueNotOptimized()
				if v == nil {
					if atomic.LoadInt32(&j) == max*c {
						return
					}
					continue
				}

				lock.Lock()
				atomic.AddInt32(&j, 1)
				jmap[v.(int)]++
				lock.Unlock()
			}
		}()
	}

	wg.Wait()

	for k, v := range jmap {
		if v != c {
			t.Fatalf("key %v got %v", k, v)
		}
	}
}

func TestQueue(t *testing.T) {
	const max = 1000
	const c = 10
	q := NewQueue()

	for i := 0; i < c; i++ {
		go func() {
			time.Sleep(time.Millisecond)
			for i := 0; i < max; i++ {
				q.Enqueue(i)
			}
		}()
	}

	var (
		lock sync.Mutex
		wg   sync.WaitGroup
		j    int32
		jmap = make(map[int]int)
	)

	for i := 0; i < c; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				v := q.Dequeue()
				if v == nil {
					if atomic.LoadInt32(&j) == max*c {
						return
					}
					continue
				}

				lock.Lock()
				atomic.AddInt32(&j, 1)
				jmap[v.(int)]++
				lock.Unlock()
			}
		}()
	}

	wg.Wait()

	for k, v := range jmap {
		if v != c {
			t.Fatalf("key %v got %v", k, v)
		}
	}
}

func BenchmarkQueueNotOptimized(b *testing.B) {
	i := uint32(0)
	q := NewQueue()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		x := atomic.AddUint32(&i, 1)
		if x%2 == 0 {
			for pb.Next() {
				q.EnqueueNotOptimized(pb)
			}
		} else {
			for pb.Next() {
				q.DequeueNotOptimized()
			}
		}
	})
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
				q.DequeueNotOptimized()
			}
		}
	})
}
