package subscriber

import (
	"sync"

	"github.com/google/uuid"
)

type Cache struct {
	mu sync.RWMutex

	data map[uuid.UUID]*Subscriber
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[uuid.UUID]*Subscriber),
	}
}

func (c *Cache) UpsertItem(key uuid.UUID, value *Subscriber) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
}

func (c *Cache) GetItem(key uuid.UUID) (*Subscriber, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, ok := c.data[key]
	return data, ok
}

func (c *Cache) RemoveItem(key uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}
