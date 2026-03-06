package lusus 

import (
	"sync"
	"time"
)

type Item struct {
	value string
	ttl time.Time
}


type LususStore struct {
	data map[string]Item
	mu sync.RWMutex
}

func NewLususStore() *LususStore {
	return &LususStore{
		data : make(map[string]Item),
	}
}

func (ls *LususStore) Set (key, value string) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	ls.data[key] = Item {
		value: value,
	}
}

func (ls *LususStore) Get (key string) (value string, isValid bool) {
	ls.mu.RLock()
	defer ls.mu.RUnlock()
	val, ok := ls.data[key]

	return val.value, ok
}

// Now lets go ahead a create a well defined system prompt that encapsulates all this discussion and states the plan to implement this system in a expo app. I will be using sqlite as my local DB so for the first phase I can install the necessary packages and initialize the db and create a working instance with our preferred schema ( make sure transit only fields are not mentioned in this schema ) In this phase I want to just setup the DB with the schema and also ensure that my DB supports storing vector embeddings