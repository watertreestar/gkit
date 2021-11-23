package leakybucket

import (
	"fmt"
	"testing"
	"time"
)

func TestLeakyBucket(t *testing.T)  {
	// New LeakyBucket that leaks at the rate of 0.5/sec and a total capacity of 10.
	b := NewLeakyBucket(0.5, 10)
	b.Add(5)
	b.Add(5)
	// Bucket is now full!
	for i := 0; i < 20; i++{
		n := b.Add(1)
		time.Sleep(500*time.Millisecond)
		// n == 0
		fmt.Println(n)
	}

}
