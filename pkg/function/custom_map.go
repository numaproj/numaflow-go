package function

import (
	"hash"
	"hash/fnv"
)

type CustomMap struct {
	hashMap map[uint64]chan Datum
	hash    hash.Hash64
}

func NewCustomMap() *CustomMap {
	return &CustomMap{
		hashMap: make(map[uint64]chan Datum),
		hash:    fnv.New64(),
	}
}

func (c *CustomMap) getIfPresent(key []string) (chan Datum, bool) {
	hashVal := c.generateHash(key)
	val, ok := c.hashMap[hashVal]
	return val, ok
}

func (c *CustomMap) putIfNotPresent(key []string, val chan Datum) {
	hashVal := c.generateHash(key)
	_, ok := c.hashMap[hashVal]
	if !ok {
		c.hashMap[hashVal] = val
	}
}

func (c *CustomMap) generateHash(key []string) uint64 {
	c.hash.Reset()
	for _, k := range key {
		_, _ = c.hash.Write([]byte(k))
	}
	return c.hash.Sum64()
}

func (c *CustomMap) closeChannels() {
	for _, c := range c.hashMap {
		close(c)
	}
}
