package cache

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

var (
	// ErrNotFound is returned when the cache
	// is consulted for an item which does not
	// exist in it.
	ErrNotFound = errors.New("not found")

	// DefaultBytesLimit is the initial capacity of
	// the memcache in terms of total bytes stored.
	// It can be set by using the MaxBytesLimit option
	// with a different value during initialization.
	DefaultBytesLimit = uint(64 * 1024 * 1024) // 64MB

	// DefaultItemLimit is the initial capacity of
	// the memcache in terms of total items stored.
	// It can be set by using the MaxItemLimit option
	// with a different value during initialization.
	DefaultItemLimit = uint(64 * 1024 * 4) // ~256K items
)

// Cache is the interface for basic cache impementations
type Cache interface {
	Fetch(key string) ([]byte, error)
	Set(key string, value []byte, exp time.Duration) error
	Delete(key string) error
}

type memcache struct {
	sync.RWMutex
	list         *list.List
	items        map[string]*list.Element
	maxItems     uint
	maxBytes     uint
	currentBytes uint
}

type cachedItem struct {
	key      string
	payload  []byte
	expireAt time.Time
}

// NewMemCache returns an in-memory implementation
// of the Cache interface.
func NewMemCache(opts ...cacheOption) Cache {
	m := &memcache{
		list:     list.New(),
		items:    map[string]*list.Element{},
		maxBytes: DefaultBytesLimit,
		maxItems: DefaultItemLimit,
	}

	for _, o := range opts {
		o(m)
	}

	return m
}

type cacheOption func(*memcache)

// MaxItemLimit sets the maximum number
// of items the cache can hold. When this
// limit is exceeded, the least accessed
// item gets deleted from the cache.
func MaxItemLimit(l uint) cacheOption {
	return func(m *memcache) {
		m.maxItems = l
	}
}

// MaxBytesLimit sets the maximum total size
// of items the cache can hold. When this
// limit is exceeded, the least accessed
// item gets deleted from the cache.
func MaxBytesLimit(b uint) cacheOption {
	return func(m *memcache) {
		m.maxBytes = b
	}
}

// Fetch retrieves an item from the cache.
// Frequently accessed items are less likely to be evicted.
func (m *memcache) Fetch(key string) ([]byte, error) {
	m.Lock()
	defer m.Unlock()
	element, ok := m.items[key]

	if !ok {
		return nil, ErrNotFound
	}

	cachedItem := element.Value.(*cachedItem)
	if cachedItem.expireAt.Before(time.Now().UTC()) {
		m.list.Remove(element)
		delete(m.items, cachedItem.key)
		return nil, ErrNotFound
	}

	m.list.MoveToFront(element)
	return cachedItem.payload, nil
}

// Set writes an item in the cache.
func (m *memcache) Set(key string, value []byte, exp time.Duration) error {
	m.Lock()
	defer m.Unlock()

	element, ok := m.items[key]
	if ok {
		item := element.Value.(*cachedItem)
		m.currentBytes -= uint(len(item.payload))
		item.payload = value
		item.expireAt = time.Now().UTC().Add(exp)
		m.currentBytes += uint(len(value))

		m.list.MoveToFront(element)
		return nil
	}

	// add new item to the list and the map
	m.currentBytes += uint(len(value))
	m.items[key] = m.list.PushFront(&cachedItem{key, value, time.Now().UTC().Add(exp)})

	// remove least accessed item if total stored items exceed limit
	for uint(m.list.Len()) > m.maxItems || m.currentBytes > m.maxBytes {
		e := m.list.Back()
		if e != nil {
			delete(m.items, e.Value.(*cachedItem).key)
			m.list.Remove(m.list.Back())
			m.currentBytes -= uint(len(e.Value.(*cachedItem).payload))
		}
	}

	return nil
}

// Delete removes an item from the cache.
func (m *memcache) Delete(key string) error {
	m.Lock()
	defer m.Unlock()

	element, ok := m.items[key]
	if ok {
		delete(m.items, key)
	}

	if element != nil {
		m.list.Remove(element)
	}

	return nil
}
