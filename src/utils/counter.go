package utils

import (
	"sync"
)

type Counter struct {
	name  string
	value int
}

var (
	counterMap = make(map[string]*Counter)
	mutex      = new(sync.RWMutex)
)

func NewCounter(name string) *Counter {
	mutex.Lock()
	defer mutex.Unlock()
	counter := counterMap[name]
	if counter == nil {
		counter = &Counter{
			name:  name,
			value: 0,
		}
		counterMap[name] = counter
	}

	return counter
}

func (c *Counter) Get() int {
	mutex.RLock()
	defer mutex.RUnlock()
	return c.value
}

func (c *Counter) Incr(delta int) {
	mutex.Lock()
	defer mutex.Unlock()
	c.value += delta
}

func DumpCounter() map[string]int {
	mutex.RLock()
	defer mutex.RUnlock()
	dump := make(map[string]int)
	for k, c := range counterMap {
		dump[k] = c.value
	}
	return dump
}
