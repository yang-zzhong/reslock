package reslock

import (
	"sync"
	"testing"
)

func TestLocker_Lock(t *testing.T) {
	locker := New()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			locker.Lock("hello")
			wg.Done()
		}()
	}
	wg.Wait()
	if locker.keys["hello"] != 10 {
		t.Fatalf("key count error")
	}
}

func TestLocker_Unlock(t *testing.T) {
	locker := New()
	locker.keys["hello"] = 10
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			locker.Unlock("hello")
			wg.Done()
		}()
	}
	wg.Wait()
	if _, ok := locker.keys["hello"]; ok {
		t.Fatalf("unlock error")
	}
}

func Benchmark_Lock(b *testing.B) {
	locker := New()
	for i := 0; i < b.N; i++ {
		locker.Lock("hello")
	}
}
