package config

import "encoding/json"

// Cache stores values indexed by a cache name and a cache key.
type Cache struct {
	data map[string]map[string]string
}

// NewCache creates a new empty Cache.
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]map[string]string),
	}
}

// Put stores a value index by a cache name and a cache key.
func (c *Cache) Put(cacheName, key, v string) {
	if c.data[cacheName] == nil {
		c.data[cacheName] = make(map[string]string)
	}
	c.data[cacheName][key] = v
}

// Get returns the value pointed to by cacheName and key.
func (c *Cache) Get(cacheName, key string) string {
	return c.data[cacheName][key]
}

// Flush removes all values from the cache.
func (c *Cache) Flush() {
	c.data = make(map[string]map[string]string)
}

// MarshalJSON makes Cache implement json.Marshaler.
func (c *Cache) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.data)
}

// UnmarshalJSON makes Cache implement json.Unmarshaler
func (c *Cache) UnmarshalJSON(data []byte) error {
	aux := make(map[string]map[string]string)
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	c.data = aux
	return nil
}
