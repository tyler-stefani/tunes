package cache

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
)

const (
	ENTRY_DURATION = 10 * time.Minute
)

type BC struct {
	cache *bigcache.BigCache
}

func NewCache() *BC {
	cache, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(ENTRY_DURATION))
	return &BC{cache}
}

func (c *BC) Get(key string) string {
	json, _ := c.cache.Get(key)
	return string(json)
}

func (c *BC) Put(key string, json string) {
	c.cache.Set(key, []byte(json))
}
