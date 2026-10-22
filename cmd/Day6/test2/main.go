package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type SafeMap struct {
	mu   sync.RWMutex
	data map[string]string
}

func (m *SafeMap) Set(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = value
}

func (m *SafeMap) Get(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, ok := m.data[key]
	return value, ok
}

func (m *SafeMap) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

func main() {
	m := SafeMap{
		mu:   sync.RWMutex{},
		data: make(map[string]string),
	}
	m.Set("name", "xiaoming")
	m.Set("age", "18")
	value, ok := m.Get("name")
	if !ok {
		return
	}
	fmt.Println(value)
	fmt.Println(m.Get("age"))
	var count atomic.Int64
	count.Store(10)
	if !count.CompareAndSwap(8, 20) {
		fmt.Println("更新失败")
		return
	}
	fmt.Printf("更新成功，当前值为：%d", count.Load())
}
