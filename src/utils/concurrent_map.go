package utils

import (
	"sync"
)

type ConcurrentMap struct {
	data  map[string]string
	mutex *sync.Mutex
}

func (m *ConcurrentMap) Get(key string) string {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.data[key]
}

func (m *ConcurrentMap) Put(key string, value string) {
	m.mutex.Lock()
	m.data[key] = value
	m.mutex.Unlock()
}

func (m *ConcurrentMap) Delete(key string) {
	m.mutex.Lock()
	delete(m.data, key)
	m.mutex.Unlock()
}

func (m *ConcurrentMap) GetData() map[string]string {
	data := make(map[string]string)
	m.mutex.Lock()
	for k,v := range m.data {
		data[k] = v
	}
	m.mutex.Unlock()
	return data
}

func (m *ConcurrentMap) SetData(data map[string]string) {
	m.mutex.Lock()
	m.data = make(map[string]string)
	for k, v := range data {
		m.data[k] = v
	}
	m.mutex.Unlock()
}

func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{
		data:  make(map[string]string),
		mutex: new(sync.Mutex),
	}
}
