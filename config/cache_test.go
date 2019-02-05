package config

import (
	"testing"
)

func TestGetOnEmptyCache(t *testing.T) {
	c := NewMapCache()
	v := c.Get("does not", "exist")
	if v != "" {
		t.Errorf("Expected empty string and got %s", v)
	}
}

func TestPutAndGet(t *testing.T) {
	c := NewMapCache()
	c.Put("c1", "k1", "v1")
	v := c.Get("c1", "k1")
	if v != "v1" {
		t.Errorf("Expected '%s' to be 'v1'", v)
	}
}
func TestFlush(t *testing.T) {
	c := NewMapCache()
	c.Put("c1", "k1", "v1")
	c.Flush()
	if c.Get("c1", "k1") != "" {
		t.Error("Expected cache to be empty")
	}
}
func TestMarshalEmptyCache(t *testing.T) {
	c := NewMapCache()
	bytes, err := c.MarshalJSON()

	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}
	if string(bytes) != "{}" {
		t.Errorf("Expected empty JSON object but got %s", bytes)
	}
}

func TestMarshalPopulatedCache(t *testing.T) {
	c := NewMapCache()
	c.Put("c1", "k1", "v1")
	c.Put("c1", "k2", "v2")
	c.Put("c2", "k1", "v1")

	bytes, err := c.MarshalJSON()
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}
	if string(bytes) != `{"c1":{"k1":"v1","k2":"v2"},"c2":{"k1":"v1"}}` {
		t.Errorf("Got unexpected JSON object: %s", bytes)
	}
}
func TestUnmarshalPopulatedCache(t *testing.T) {
	c := NewMapCache()

	err := c.UnmarshalJSON([]byte(`{"c1":{"k1":"v1","k2":"v2"},"c2":{"k1":"v1"}}`))
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	v := c.Get("c1", "k1")
	if v != "v1" {
		t.Errorf("Expected '%s' to be 'v1'", v)
	}
	v = c.Get("c1", "k2")
	if v != "v2" {
		t.Errorf("Expected '%s' to be 'v2'", v)
	}
	v = c.Get("c2", "k1")
	if v != "v1" {
		t.Errorf("Expected '%s' to be 'v1'", v)
	}
}

func TestUnmarshalBrokenJSON(t *testing.T) {
	c := NewMapCache()
	err := c.UnmarshalJSON([]byte(`broken`))
	if err == nil {
		t.Errorf("Expected non-nil error")
	}
}
