package main

import "fmt"

// Small, focused interface
type Storer interface {
	Store(key, value string)
}

type Fetcher interface {
	Fetch(key string) (string, bool)
}

// Composed interface — only use this when you need both
type Cache interface {
	Storer
	Fetcher
}

// Concrete type
type MemCache struct {
	data map[string]string
}

func NewMemCache() *MemCache {
	return &MemCache{data: make(map[string]string)}
}

func (m *MemCache) Store(key, value string) {
	m.data[key] = value
}

func (m *MemCache) Fetch(key string) (string, bool) {
	v, ok := m.data[key]
	return  v, ok
}

// This function only needs to store — it takes the narrow interface
func Populate(s Storer) {
	s.Store("lang", "Go")
	s.Store("day", "1")
}

// This function only needs to fetch
func Display(f Fetcher) {
	if v, ok := f.Fetch("lang"); ok {
		fmt.Println("lang:", v)
	}
}

func main() {
	c := NewMemCache()
	Populate(c)
	Display(c)
}