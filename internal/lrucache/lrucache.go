package lrucache

import (
	"container/list"
	"errors"
	"fmt"
	"sync"
)

type LRUCache[K comparable] interface {
	// Add adds a new key value pair and moves them to front. it returns an error if key already exists
	Add(K, interface{}) (interface{}, error)
	// Read returns error if key does not exist. ow, it moves the cache node
	Read(key K) (interface{}, error)
	// ReadSafe returns error if key does not exist. it never moves a cache node
	ReadSafe(key K) (any, error)
	// Remove returns error if key does not exist. ow, it removes the node from cache and returns it value
	Remove(key K) (any, error)
	// ReadMostRecent returns error if cache is empty. it never moves a cache node
	ReadMostRecent() (any, error)
	// Update returns error if key does not exist. ow, it updates value and moves the cache node
	Update(key K, v interface{}) error
	// UpdateSafe returns error if key does not exist. ow, it updates value without moving the cache node
	UpdateSafe(key K, value interface{}) error
}

type lruCache[K comparable] struct {
	queue      *list.List
	dictionary map[K]*list.Element
	size       int
	mu         sync.Mutex
}

func NewLRUCache[K comparable](bufferSize int) LRUCache[K] {
	return &lruCache[K]{
		queue:      list.New(),
		dictionary: make(map[K]*list.Element, bufferSize),
		size:       bufferSize,
	}
}

// Add adds a new key value pair and moves them to front. it returns an error if key already exists
func (l *lruCache[K]) Add(key K, value interface{}) (any, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.dictionary[key]; ok {
		message := fmt.Sprintf("cannot add an already existing element to lru cache. element %v already exists", key)
		return nil, errors.New(message)
	}
	cNode := CacheNode[K, any]{key, value}
	element := l.queue.PushFront(cNode)
	l.dictionary[key] = element
	var removedValue interface{} = nil
	if l.queue.Len() > l.size {
		removedValue = l.removeLeastUsed()
	}
	return removedValue, nil
}

// Read returns error if key does not exist. ow, it moves the cache node
func (l *lruCache[K]) Read(key K) (any, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	element, ok := l.dictionary[key]
	if !ok {
		return nil, errors.New("failed reading lruCache element. key does not exist")
	}
	l.queue.MoveToFront(element)
	return element.Value.(CacheNode[K, any]).Value, nil
}

func (l *lruCache[K]) removeLeastUsed() any {
	element := l.queue.Back()
	node := l.queue.Remove(element)
	delete(l.dictionary, node.(CacheNode[K, any]).key)
	return element.Value.(CacheNode[K, any]).Value
}

// ReadMostRecent returns error if cache is empty. it never moves a cache node
func (l *lruCache[K]) ReadMostRecent() (any, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.dictionary) == 0 {
		return nil, errors.New("empty cache")
	}
	return l.queue.Front().Value.(CacheNode[K, any]).Value, nil
}

// ReadSafe returns error if key does not exist. it never moves a cache node
func (l *lruCache[K]) ReadSafe(key K) (any, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	element, ok := l.dictionary[key]
	if !ok {
		return nil, errors.New("failed reading lruCache element. key does not exist")
	}
	return element.Value.(CacheNode[K, any]).Value, nil
}

// Remove returns error if key does not exist. ow, it removes the node from cache
func (l *lruCache[K]) Remove(key K) (any, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	element, ok := l.dictionary[key]
	if !ok {
		return nil, errors.New("failed removing lruCache element. key does not exist")
	}
	l.queue.Remove(element)
	delete(l.dictionary, key)
	return element.Value.(CacheNode[K, any]).Value, nil
}

// Update returns error if key does not exist. ow, it updates value and moves the cache node
func (l *lruCache[K]) Update(key K, value interface{}) error {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	if _, err := l.Read(key); err != nil {
		return err
	}
	l.dictionary[key].Value = CacheNode[K, any]{key, value}
	return nil
}

// UpdateSafe returns error if key does not exist. ow, it updates value without moving the cache node
func (l *lruCache[K]) UpdateSafe(key K, value interface{}) error {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	if _, err := l.ReadSafe(key); err != nil {
		return err
	}
	l.dictionary[key].Value = CacheNode[K, any]{key, value}
	return nil
}

type CacheNode[K comparable, V any] struct {
	key   K
	Value V
}
