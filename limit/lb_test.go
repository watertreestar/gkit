package limit

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestLeakyBucketLimiter(t *testing.T) {
	var counter int
	l := NewLeakyBucket(10, 1, 2*time.Second) // one req five second
	defer l.Close()
	//time.Sleep(5 * time.Second)
	prev := time.Now()
	for i := 0; i < 20; i++ {
		now := l.TakeBlocked(1)
		counter++
		fmt.Println(i, now.Sub(prev))
		prev = now
	}
	t.Logf("%v requests allowed.", counter)
}

func TestLeakyBucketLimiterNonBlocked(b *testing.T) {
	l := NewLeakyBucket(10,1, 2*time.Second) // one req per 10ms
	defer l.Close()
	l.Take(10)
	for i := 0; i < 20; i++ {
		ok := l.Take(1)
		time.Sleep(500 * time.Millisecond)
		fmt.Println(ok)
	}
}

func TestTB(b *testing.T) {
	var wg sync.WaitGroup
	var counter int64

	l := NewLeakyBucket(100, 100, time.Second)
	defer l.Close()

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runner(b, l, &counter)
		}()
	}
	wg.Wait()
	b.Logf("%v requests allowed.", counter)
}

func runner(b *testing.T, l Limiter, counter *int64) {
	for i := 0; i < 100; i++ {
		ok := l.Take(1)
		if !ok {
			b.Logf("#%d Take() returns not ok, available: %v", i, l.Available())
			time.Sleep(100 * time.Millisecond)
		} else {
			//b.Logf("OK: #%d Take(), counter: %v", i, l.count)
			atomic.AddInt64(counter, 1)
			time.Sleep(time.Duration(rand.Intn(15-5)+5) * time.Millisecond)
		}
	}
}

