package reslock

import "sync"

var (
	StdLocker Locker
)

type Locker struct {
	keys map[string]int
	lock sync.RWMutex
}

func init() {
	StdLocker = Locker{keys: make(map[string]int)}
}

func New() *Locker {
	return &Locker{keys: make(map[string]int)}
}

func (l *Locker) Lock(key string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.keys[key] += 1
}

func (l *Locker) Locked(key string) bool {
	l.lock.RLock()
	defer l.lock.RUnlock()
	if _, ok := l.keys[key]; ok {
		return true
	}
	return false
}

func (l *Locker) Unlock(key string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if _, ok := l.keys[key]; !ok {
		return
	}
	l.keys[key] -= 1
	if l.keys[key] <= 0 {
		delete(l.keys, key)
	}
}
