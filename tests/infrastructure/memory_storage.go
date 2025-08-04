package infrastructure

import (
	"fmt"
	"sync"
)

// InMemoryStorage provides an in-memory implementation of the Storage interface for testing
type InMemoryStorage struct {
	data  map[string]interface{}
	mutex sync.RWMutex
}

// NewInMemoryStorage creates a new in-memory storage instance
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]interface{}),
	}
}

// Store saves a value with the given key
func (s *InMemoryStorage) Store(key string, value interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	s.data[key] = value
	return nil
}

// Get retrieves a value by key
func (s *InMemoryStorage) Get(key string) (interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	value, exists := s.data[key]
	if !exists {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	
	return value, nil
}

// Exists checks if a key exists in storage
func (s *InMemoryStorage) Exists(key string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	_, exists := s.data[key]
	return exists
}

// ListAll retrieves all stored values
func (s *InMemoryStorage) ListAll() ([]interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	values := make([]interface{}, 0, len(s.data))
	for _, value := range s.data {
		values = append(values, value)
	}
	
	return values, nil
}

// Delete removes a value by key
func (s *InMemoryStorage) Delete(key string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if _, exists := s.data[key]; !exists {
		return fmt.Errorf("key not found: %s", key)
	}
	
	delete(s.data, key)
	return nil
}