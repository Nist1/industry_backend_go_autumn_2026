package main

import (
	"container/list"
	"sync"
)

type LRU[K comparable, V any] interface {
	Get(key K) (value V, ok bool)
	Set(key K, value V)
}

type entry[K comparable, V any] struct {
	key   K
	value V
}

type LRUCache[K comparable, V any] struct {
	capacity int
	mu       sync.Mutex
	ll       list.List
	items    map[K]*list.Element
}

func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity: capacity,
		items:    make(map[K]*list.Element, max(capacity, 0)),
	}
}

func (c *LRUCache[K, V]) Get(key K) (value V, ok bool) {
	c.mu.Lock()
	// defer гарантирует разблокировку при любом выходе из функции
	defer c.mu.Unlock()

	elem, ok := c.items[key]
	if !ok {
		return value, false
	}
	c.ll.MoveToFront(elem)

	return elem.Value.(*entry[K, V]).value, true
}

func (c *LRUCache[K, V]) Set(key K, value V) {
	if c.capacity <= 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		elem.Value.(*entry[K, V]).value = value
		return
	}
	if c.ll.Len() >= c.capacity {
		oldest := c.ll.Back()
		c.ll.Remove(oldest)
		delete(c.items, oldest.Value.(*entry[K, V]).key)
	}
	c.items[key] = c.ll.PushFront(&entry[K, V]{key: key, value: value})
}
