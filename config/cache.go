package config

import "encoding/json"

type Cache struct {
	data map[string]map[string]string
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]map[string]string),
	}
}

func (c *Cache) Put(k1, k2, v string) {
	if c.data[k1] == nil {
		c.data[k1] = make(map[string]string)
	}
	c.data[k1][k2] = v
}

func (c *Cache) Get(k1, k2 string) string {
	return c.data[k1][k2]
}

func (c *Cache) Flush() {
	c.data = make(map[string]map[string]string)
}

func (c *Cache) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.data)
}

func (c *Cache) UnmarshalJSON(data []byte) error {
	aux := make(map[string]map[string]string)
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	c.data = aux
	return nil
}
