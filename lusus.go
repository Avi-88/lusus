package lusus 

import (
	"time"
	"fmt"
	"sync"
)

type item struct {
	value string
	expiresAt time.Time
}

type LususStore struct {
	data map[string]item
	mu sync.RWMutex
}

func NewLususStore() *LususStore {
	return &LususStore{
		data: make(map[string]item),
	}
}

func (ls *LususStore) Set(key, value string, ttl ...int) error {
	var expiresAt time.Time
	if len(ttl) > 0 {
		if ttl[0] > 0 {
			expiresAt = time.Now().Add(time.Duration(ttl[0])*time.Second)
		} else {
			return fmt.Errorf("ttl must be non-negative, got %d", ttl[0])
		}
	}
	ls.mu.Lock()
	defer ls.mu.Unlock()
	ls.data[key] = item {
		value: value,
		expiresAt: expiresAt,
	}
	return nil
}

func (ls *LususStore) Get(key string) (string, bool) {
	ls.mu.RLock()
	defer ls.mu.RUnlock()
	it, ok := ls.data[key]
	if !ok {
		return "", false
	} else if time.Now().After(it.expiresAt) {
		ls.mu.RUnlock()
		ls.mu.Lock()
		delete(ls.data, key)
		ls.mu.Unlock()
		ls.mu.RLock()
		return "", false
	}
	return it.value, true
}

func (ls *LususStore) Delete(key string) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	_ , exists := ls.data[key]
	if !exists {
		return fmt.Errorf("key not found")
	}
	delete(ls.data, key)
	return nil
}
