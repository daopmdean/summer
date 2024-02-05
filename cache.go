package summer

import (
	"sync"
	"time"
)

type LCache struct {
	m          map[string]*cacheItem
	l          sync.RWMutex
	maxItem    int
	refreshTTL bool
}

type cacheItem struct {
	value        interface{}
	lastAccessed int64
}

// NewLCache ...
// maxItem: maximum item
// ttl: time to live (second)
// short hand for NewLCacheRefreshMode(maxItem, ttl, true)
func NewLCache(maxItem int, ttl int) (m *LCache) {
	m = NewLCacheRefreshMode(maxItem, ttl, true)
	return
}

// NewLCacheRefreshMode ...
// maxItem: maximum item
// ttl: time to live (second)
// refreshTTL : if refreshTTL == true, when someone access item, item ttl will be refreshed
func NewLCacheRefreshMode(maxItem int, ttl int, refreshTTL bool) (m *LCache) {
	m = &LCache{m: make(map[string]*cacheItem, maxItem), maxItem: maxItem, refreshTTL: refreshTTL}
	go func() {
		for now := range time.Tick(time.Second) {
			m.l.Lock()
			for k, v := range m.m {
				if now.Unix()-v.lastAccessed > int64(ttl) {
					delete(m.m, k)
				}
			}
			m.l.Unlock()
		}
	}()
	return
}

// Get cache length
func (c *LCache) Len() int {
	c.l.RLock()
	defer c.l.RUnlock()

	return len(c.m)
}

// Put data to cache
func (c *LCache) Put(k string, v interface{}) {
	c.l.Lock()

	c.m[k] = &cacheItem{
		value:        v,
		lastAccessed: time.Now().Unix(),
	}

	c.l.Unlock()
}

// Get data from cache
func (c *LCache) Get(k string) (v interface{}, ok bool) {
	if !c.refreshTTL {
		c.l.RLock()
		if ci, okv := c.m[k]; okv {
			ok = okv
			v = ci.value
		}
		c.l.RUnlock()
		return
	}

	c.l.Lock()
	if ci, okv := c.m[k]; okv {
		ok = okv
		v = ci.value
		ci.lastAccessed = time.Now().Unix()
	}
	c.l.Unlock()
	return
}

// Check if key is exist
func (c *LCache) ContainsKey(k string) (ok bool) {
	c.l.Lock()
	if ci, okv := c.m[k]; okv {
		ok = okv
		if c.refreshTTL {
			ci.lastAccessed = time.Now().Unix()
		}
	}
	c.l.Unlock()
	return
}

// Remove data in cache
func (c *LCache) Remove(k string) {
	c.l.Lock()
	delete(c.m, k)
	c.l.Unlock()
}

// Remove all data in cache
func (c *LCache) Cleanup() {
	c.l.Lock()
	c.m = make(map[string]*cacheItem, c.maxItem)
	c.l.Unlock()
}
