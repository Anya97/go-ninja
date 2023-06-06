package cache

import "C"
import (
	"errors"
	"runtime"
	"sync"
	"time"
)

type Cleaner struct {
	Interval time.Duration
	stop     chan bool
}

type Data struct {
	Value    interface{}
	ExpireAt int64
}

type Cache struct {
	cacheMap map[string]Data
	ttl      time.Duration
	mutex    *sync.RWMutex
	cleaner  *Cleaner
}

func New(ttl time.Duration, cleanUpInterval time.Duration) *Cache {
	cache := &Cache{
		cacheMap: make(map[string]Data),
		ttl:      ttl,
		mutex:    &sync.RWMutex{},
	}

	if cleanUpInterval > 0 {
		clean(cleanUpInterval, cache)
		runtime.SetFinalizer(cache, stopCleaning)
	}

	return cache
}

func clean(cleanUpInterval time.Duration, cache *Cache) {
	cleaner := &Cleaner{
		Interval: cleanUpInterval,
		stop:     make(chan bool),
	}

	cache.cleaner = cleaner
	go cleaner.Cleaning(cache)
}

func stopCleaning(cache *Cache) {
	cache.cleaner.stop <- true
}

func (c *Cleaner) Cleaning(cache *Cache) {
	ticker := time.NewTicker(c.Interval)

	for {
		select {
		case <-ticker.C:
			cache.purge()
		case <-c.stop:
			ticker.Stop()
		}
	}
}

func (c *Cache) purge() {
	now := time.Now().UnixNano()
	for key, data := range c.cacheMap {
		if data.ExpireAt < now {
			delete(c.cacheMap, key)
		}
	}
}

func (c *Cache) Set(key string, value interface{}) {
	c.mutex.Lock()
	c.cacheMap[key] = Data{Value: value, ExpireAt: time.Now().Add(c.ttl).UnixNano()}
	c.mutex.Unlock()
	return
}

func (c *Cache) Get(key string) (interface{}, error) {
	c.mutex.RLock()
	valueByKey, exists := c.cacheMap[key]
	c.mutex.RUnlock()
	if exists {
		return valueByKey.Value, nil
	}
	return nil, errors.New("No such key. Please, use valid key.")
}

func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	delete(c.cacheMap, key)
	c.mutex.Unlock()
	return
}
