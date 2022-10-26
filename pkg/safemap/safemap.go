package safemap

import (
	"fmt"
	"sync"
)

type SafeMap[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

func New[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		data: make(map[K]V),
	}
}

func (s *SafeMap[K, V]) Delete(key K) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[key]
	if !ok {
		return fmt.Errorf("key %v not found", key)
	}

	delete(s.data, key)

	return nil
}

func (s *SafeMap[K, V]) GetAllKeys() []K {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]K, 0)

	for key := range s.data {
		keys = append(keys, key)
	}

	return keys
}

func (s *SafeMap[K, V]) GetAllValues() []V {
	s.mu.RLock()
	defer s.mu.RUnlock()

	values := make([]V, 0)

	for _, value := range s.data {
		values = append(values, value)
	}

	return values
}

func (s *SafeMap[K, V]) GetOneValue(key K) (V, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.data[key]
	if !ok {
		return value, fmt.Errorf("key %v not found", key)
	}

	return value, nil
}

func (s *SafeMap[K, V]) HasKey(key K) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data[key]

	return ok
}

func (s *SafeMap[K, V]) Insert(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
}

func (s *SafeMap[K, V]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data)
}

func (s *SafeMap[K, V]) Updade(key K, value V) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[key]
	if !ok {
		return fmt.Errorf("key %v not found", key)
	}

	s.data[key] = value

	return nil
}