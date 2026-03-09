package store

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
	"strconv"
)

type item struct {
	value string
	expiresAt time.Time
}

type LususStore struct {
	data map[string]item
	mu sync.RWMutex
	aofFile *os.File
}

func NewLususStore(aofFilename string) *LususStore {
	ls := &LususStore{
		data: make(map[string]item),
	}

	// Load persisted data is present
	if _, err := os.Stat(aofFilename) ; err == nil {
		_ = ls.loadAOF(aofFilename)
	}else{
		fmt.Println("No persisted data found")
	}

	// Open AOF file for persistence
	if f, err := os.OpenFile(aofFilename, os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0644 ); err == nil {
		ls.aofFile = f
	}else{
		fmt.Fprintf(os.Stderr, "Warning: could not open AOF file: %v\n", err)
	}

	return ls
}

func (ls *LususStore) loadAOF(aofFilename string) error {
	f, err := os.Open(aofFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan(){
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 4)
		cmd := strings.ToUpper(parts[0])
		key:= parts[1]
		switch cmd {
		case "SET":
			val := UnescapeNewlines(parts[2])
			if len(parts) == 4 {
				if ttl, err := strconv.Atoi(parts[3]); err == nil {
					ls.data[key] = item{
						value: val,
						expiresAt: time.Now().Add(time.Duration(ttl)*time.Second),
					}
				}
			}
			if len(parts) == 3 {
				ls.data[key] = item{
					value: val,
				}
			}
		case "DEL":
			if len(parts) >= 2 {
				delete(ls.data, key)
			}
		}
	}
	return scanner.Err()
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
	ls.data[key] = item{
		value: value,
		expiresAt: expiresAt,
	}
	if ls.aofFile != nil {
		if len(ttl) > 0 {
			fmt.Fprintf(ls.aofFile, "SET %s %s %d\r\n", key, EscapeNewlines(value), ttl[0])
		}else{
			fmt.Fprintf(ls.aofFile, "SET %s %s\r\n", key, EscapeNewlines(value))
		}
	}
	return nil
}

func (ls *LususStore) Get(key string) (string, bool) {
	ls.mu.RLock()
	it, ok := ls.data[key]
	ls.mu.RUnlock()
	if !ok {
		return "", false
	} 
	
	if !it.expiresAt.IsZero() && time.Now().After(it.expiresAt) {
		ls.mu.Lock()
		delete(ls.data, key)
		ls.mu.Unlock()
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
	if ls.aofFile != nil {
		fmt.Fprintf(ls.aofFile, "DEL %s\r\n", key)
	}
	return nil
}
