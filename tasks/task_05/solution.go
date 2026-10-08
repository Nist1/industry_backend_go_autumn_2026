package main

type Cache[K comparable, V any] struct {
	capacity int
	items    map[K]V
}

func NewCache[K comparable, V any](capacity int) *Cache[K, V] {
	return &Cache[K, V]{
		capacity: capacity,
		items:    make(map[K]V, max(capacity, 0)),
	}
}

func (c *Cache[K, V]) Get(k K) (v V, ok bool) {
	// Проверка, чтобы отличить нулевое значение от несуществующего ключа
	v, ok = c.items[k]
	return v, ok
}

func (c *Cache[K, V]) Set(k K, v V) bool {
	if c.capacity <= 0 {
		return false
	}
	// Отказываем только новому ключу при заполненном кэше.
	// Существующий ключ обновляем всегда
	if _, exists := c.items[k]; !exists && len(c.items) >= c.capacity {
		return false
	}
	c.items[k] = v
	return true
}
