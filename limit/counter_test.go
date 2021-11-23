package limit

import (
	"testing"
	"time"
)

func TestCounterLimiter(b *testing.T) {
	var counter int
	l := NewCounter(100, time.Second) // one req per 10ms
	defer l.Close()
	for i := 0; i < 120; i++ {
		ok := l.Take(1)
		if !ok {
			b.Logf("#%d Take() returns not ok, remained ticks: %vns, counter: %v", i, l.(interface{ Ticks() int64 }).Ticks()-time.Now().UnixNano(), l.(interface{ Count() int }).Count())
			time.Sleep(100 * time.Millisecond)
		} else {
			//time.Sleep(5 * time.Millisecond)
			counter++
		}
	}
	b.Logf("%v requests allowed.", counter)
}
