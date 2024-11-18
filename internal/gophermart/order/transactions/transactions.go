package transactions

import (
	"sync"
	"time"
)

const (
	SLEEP = 10 // сколько спит перед проверкой разблокирована ли транзакция
)

type Mutex struct {
	keys map[string] struct{}
	mu sync.Mutex
}

func New() (*Mutex, error) {
	var t Mutex
	t.keys = make(map[string]struct{})
	return &t, nil
}

func (t *Mutex) Lock(key string) {
	for {
		time.Sleep(SLEEP * time.Millisecond)
		t.mu.Lock()
		_, ok := t.keys[key]
		if ok {
			t.mu.Unlock()
			continue
		} else {
			t.keys[key] = struct{}{}
			t.mu.Unlock()
			break
		}
	}
}

func (t *Mutex) Unlock(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.keys, key)
}