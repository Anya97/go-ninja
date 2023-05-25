package cache

import "errors"

type Cache struct {
	cacheMap map[string]interface{}
}

func New() *Cache {
	return &Cache{
		cacheMap: make(map[string]interface{}),
	}
}

func (c *Cache) Set(key string, value interface{}) {
	c.cacheMap[key] = value
	return
}

func (c *Cache) Get(key string) (interface{}, error) {
	valueByKey, exists := c.cacheMap[key]
	if exists {
		return valueByKey, nil
	}
	return nil, errors.New("No such key. Please, use valid key.")
}

func (c *Cache) Delete(key string) {
	delete(c.cacheMap, key)
	return
}
