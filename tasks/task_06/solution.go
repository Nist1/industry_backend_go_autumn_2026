package main

import "container/list"

type entry[K comparable, V any] struct {
	key   K
	value V
}

type LRUCache[K comparable, V any] struct {
	capacity int
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
	if elem, ok := c.items[key]; ok {
		// перезапись меняет только значение
		// значение меняется через указатель
		elem.Value.(*entry[K, V]).value = value
		return
	}
	if c.ll.Len() >= c.capacity {
		// сначала вытесняем самый давний
		// и удаляем ключ из map
		oldest := c.ll.Back()
		c.ll.Remove(oldest)
		delete(c.items, oldest.Value.(*entry[K, V]).key)
	}

	c.items[key] = c.ll.PushFront(&entry[K, V]{key: key, value: value})
}

type LRU[K comparable, V any] interface {
	Get(K) (V, bool)
	Set(K, V)
}
