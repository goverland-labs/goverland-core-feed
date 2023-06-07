package subscriber

import "sync"

type Cache struct {
	mu sync.RWMutex

	data map[string]*Subscriber
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]*Subscriber),
	}
}

func (c *Cache) UpsertItem(key string, value *Subscriber) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
}

func (c *Cache) GetItem(key string) (*Subscriber, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, ok := c.data[key]
	return data, ok
}

func (c *Cache) RemoveItem(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}
