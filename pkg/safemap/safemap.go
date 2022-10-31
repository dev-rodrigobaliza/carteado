package safemap

import (
	"errors"
	"sync"
)

const (
	MAX_INT = int(^uint(0) >> 1)
)

var (
	errNotFoundIndex = errors.New("index not found")
	errNotFoundKey   = errors.New("key not found")
	errNotFoundValue = errors.New("value not found")
)

type SafeMap[K comparable, V any] struct {
	data      map[K]V
	index     map[K]int
	nextIndex int
	mu        sync.RWMutex
}

func New[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		data:      make(map[K]V),
		index:     make(map[K]int),
		nextIndex: 0,
	}
}

func (s *SafeMap[K, V]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[K]V)
	s.index = make(map[K]int)
	s.nextIndex = 0
}

func (s *SafeMap[K, V]) DeleteKey(key K) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.index[key]
	if !ok {
		return errNotFoundKey
	}

	delete(s.data, key)
	delete(s.index, key)

	return nil
}

func (s *SafeMap[K, V]) DeleteIndex(index int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for key, idx := range s.index {
		if idx == index {
			delete(s.data, key)
			delete(s.index, key)

			return nil
		}
	}

	return errNotFoundIndex
}

func (s *SafeMap[K, V]) GetAllKeys() []K {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]K, 0)

	for key := range s.index {
		keys = append(keys, key)
	}

	return keys
}

func (s *SafeMap[K, V]) GetAllIndexes() []int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	indexes := make([]int, 0)

	for _, index := range s.index {
		indexes = append(indexes, index)
	}

	return indexes
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

func (s *SafeMap[K, V]) GetFirstValue(remove bool) (V, error) {
	var key K
	var value V
	var err error
	index := -1

	s.mu.RLock()

	for k, idx := range s.index {
		if index == -1 || idx < index {
			key = k
			index = idx
		}
	}

	s.mu.RUnlock()

	if index == -1 {
		err = errNotFoundValue
		return value, err
	}

	value = s.data[key]

	if remove {
		s.DeleteKey(key)
	}

	return value, nil
}

func (s *SafeMap[K, V]) GetFirstKey() (K, error) {
	var key K
	var err error
	index := -1

	s.mu.RLock()
	defer s.mu.RUnlock()

	for k, idx := range s.index {
		if index == -1 || idx < index {
			key = k
			index = idx
		}
	}

	if index == -1 {
		err = errNotFoundValue
		return key, err
	}

	return key, nil
}

func (s *SafeMap[K, V]) GetIndexedValue(index int, remove bool) (V, error) {
	var key K
	var value V
	ok := false

	s.mu.RLock()

	for k, idx := range s.index {
		if idx == index {
			key = k
			value = s.data[k]

			ok = true
		}
	}

	s.mu.RUnlock()

	if ok {
		if remove {
			s.DeleteKey(key)
		}

		return value, nil
	}

	return value, errNotFoundValue
}

func (s *SafeMap[K, V]) GetLastValue(remove bool) (V, error) {
	var key K
	var value V
	index := MAX_INT

	s.mu.RLock()

	for k, idx := range s.index {
		if index == MAX_INT || idx > index {
			key = k
			index = idx
		}
	}

	s.mu.RUnlock()

	if index == MAX_INT {
		return value, errNotFoundValue
	}

	value = s.data[key]

	if remove {
		s.DeleteKey(key)
	}

	return value, nil
}

func (s *SafeMap[K, V]) GetOneValue(key K, remove bool) (V, error) {
	var value V
	var err error

	s.mu.RLock()

	value, ok := s.data[key]

	s.mu.RUnlock()

	if !ok {
		err = errNotFoundKey
		return value, err
	}

	if remove {
		s.DeleteKey(key)
	}

	return value, nil
}

func (s *SafeMap[K, V]) HasIndex(index int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, idx := range s.index {
		if idx == index {
			return true
		}
	}

	return false
}

func (s *SafeMap[K, V]) HasKey(key K) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data[key]

	return ok
}

func (s *SafeMap[K, V]) Insert(key K, value V) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	index := s.nextIndex
	s.nextIndex++

	s.data[key] = value
	s.index[key] = index

	return index
}

func (s *SafeMap[K, V]) IsEmpty() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data) == 0
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
		return errNotFoundKey
	}

	s.data[key] = value

	return nil
}
