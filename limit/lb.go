package limit

import (
	"sync/atomic"
	"time"
)

func NewLeakyBucket(maxCount int64, count int64, d time.Duration) Limiter  {
	return (&leakyBucket{
		true,
		maxCount,
		make(chan struct{}),
		int64(d) / count,
		time.Now().UnixNano(),
		0,
	}).start()
}

// rete limiter based leakBucket
type leakyBucket struct {
	enabled     bool
	Maximal     int64
	exitCh      chan struct{}
	rate        int64
	refreshTime int64 // in nanoseconds
	count       int64
}


func (s *leakyBucket) Enabled() bool     { return s.enabled }
func (s *leakyBucket) SetEnabled(b bool) { s.enabled = b }
func (s *leakyBucket) Count() int64      { return atomic.LoadInt64(&s.count) }
func (s *leakyBucket) Available() int64  { return atomic.LoadInt64(&s.count) }
func (s *leakyBucket) Capacity() int64   { return s.Maximal }
func (s *leakyBucket) Close() {
	close(s.exitCh)
}

func (s *leakyBucket) start() *leakyBucket {
	if s.rate < 1000 {
		return nil
	}
	return s
}

func (s *leakyBucket) max(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}

func (s *leakyBucket) take(count int) (requestAt time.Time, ok bool) {
	requestAt = time.Now()

	rtm := atomic.LoadInt64(&s.refreshTime)
	cnt := atomic.LoadInt64(&s.count)
	// 桶中最少是0，也就是漏完，不能为负
	//fmt.Println((requestAt.UnixNano()-rtm)/s.rate)
	//fmt.Println(cnt-(requestAt.UnixNano()-rtm)/s.rate)
	atomic.StoreInt64(&s.count, s.max(0, cnt-(requestAt.UnixNano()-rtm)/s.rate))
	cnt = atomic.LoadInt64(&s.count)
	// fmt.Println(cnt)
	if cnt + int64(count) > s.Maximal {
		ok = false
		return
	}
	atomic.AddInt64(&s.count, int64(count))
	atomic.StoreInt64(&s.refreshTime, requestAt.UnixNano())
	ok = true
	return
}

func (s *leakyBucket) Take(count int) (ok bool) {
	_, ok = s.take(count)
	return
}

func (s *leakyBucket) TakeBlocked(count int) (requestAt time.Time) {
	var ok bool
	requestAt, ok = s.take(count)
	for !ok {
		time.Sleep(time.Duration(s.rate - (1000 - 1)))
		_, ok = s.take(count)
	}
	time.Sleep(time.Duration(s.rate-int64(time.Now().Sub(requestAt))) - time.Millisecond)
	return
}

func (s *leakyBucket) Acquire()  bool {
	return s.Take(1)
}

func (s *leakyBucket) AcquireBlocked() (requestAt time.Time) {
	return s.TakeBlocked(1)
}