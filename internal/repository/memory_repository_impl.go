package repository

import (
	"fmt"
	"sync"
)

type MemoryRepositoryImpl struct {
	data  map[string]int64
	mutex sync.RWMutex
}

func CreateMemoryRepository() *MemoryRepositoryImpl {
	return &MemoryRepositoryImpl{
		data: make(map[string]int64),
	}
}

func (m *MemoryRepositoryImpl) Set(key string, value int64) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.data[key]; exists {
		return fmt.Errorf("переменная %s уже определена", key)
	}

	m.data[key] = value
	return nil
}

func (m *MemoryRepositoryImpl) Get(key string) (int64, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	val, exists := m.data[key]
	return val, exists
}

func (m *MemoryRepositoryImpl) Exists(key string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	_, exists := m.data[key]
	return exists
}
