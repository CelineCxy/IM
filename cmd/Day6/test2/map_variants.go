package main

import "sync"

// ConcurrentMap 统一三种 Map 的操作。
type ConcurrentMap interface {
	Set(key, value string)
	Get(key string) (string, bool)
	Delete(key string)
}

// MutexMap 使用普通互斥锁。
// 无论读取还是写入，同一时间只允许一个 goroutine 操作。
type MutexMap struct {
	mu   sync.Mutex
	data map[string]string
}

func NewMutexMap() *MutexMap {
	return &MutexMap{
		data: make(map[string]string),
	}
}

func (m *MutexMap) Set(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = value
}

func (m *MutexMap) Get(key string) (string, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	value, ok := m.data[key]
	return value, ok
}

func (m *MutexMap) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, key)
}

// RWMutexMap 就是 main.go 中已经实现的 SafeMap。
// 使用类型别名，避免重复编写相同代码。
type RWMutexMap = SafeMap

func NewRWMutexMap() *RWMutexMap {
	return &RWMutexMap{
		data: make(map[string]string),
	}
}

// SyncMap 使用标准库提供的 sync.Map。
// sync.Map 的零值可以直接使用，不需要初始化内部 map。
type SyncMap struct {
	data sync.Map
}

func NewSyncMap() *SyncMap {
	return &SyncMap{}
}

func (m *SyncMap) Set(key, value string) {
	m.data.Store(key, value)
}

func (m *SyncMap) Get(key string) (string, bool) {
	value, ok := m.data.Load(key)
	if !ok {
		return "", false
	}

	stringValue, ok := value.(string)
	return stringValue, ok
}

func (m *SyncMap) Delete(key string) {
	m.data.Delete(key)
}
