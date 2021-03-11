// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cache

import (
	"fmt"
	"time"

	"github.com/sitename/sitename/modules/graceful"
	"github.com/sitename/sitename/modules/nosql"

	"gitea.com/go-chi/cache"
	"github.com/go-redis/redis/v8"
	"github.com/unknwon/com"
)

// RedisCacher represents a redis cache adapter implementation.
type RedisCacher struct {
	c          redis.UniversalClient
	prefix     string
	hsetName   string
	occupyMode bool
}

// Put puts value into cache with key and expire time.
// If expired is 0, it lives forever.
func (c *RedisCacher) Put(key string, val interface{}, expire int64) error {
	key = c.prefix + key
	if expire == 0 {
		if err := c.c.Set(graceful.GetManager().HammerContext(), key, com.ToStr(val), 0).Err(); err != nil {
			return err
		}
	} else {
		dur, err := time.ParseDuration(com.ToStr(expire) + "s")
		if err != nil {
			return err
		}
		if err = c.c.Set(graceful.GetManager().HammerContext(), key, com.ToStr(val), dur).Err(); err != nil {
			return err
		}
	}

	if c.occupyMode {
		return nil
	}
	return c.c.HSet(graceful.GetManager().HammerContext(), c.hsetName, key, "0").Err()
}

// Get gets cached value by given key.
func (c *RedisCacher) Get(key string) interface{} {
	val, err := c.c.Get(graceful.GetManager().HammerContext(), c.prefix+key).Result()
	if err != nil {
		return nil
	}
	return val
}

// Delete deletes cached value by given key.
func (c *RedisCacher) Delete(key string) error {
	key = c.prefix + key
	if err := c.c.Del(graceful.GetManager().HammerContext(), key).Err(); err != nil {
		return err
	}

	if c.occupyMode {
		return nil
	}
	return c.c.HDel(graceful.GetManager().HammerContext(), c.hsetName, key).Err()
}

// Incr increases cached int-type value by given key as a counter.
func (c *RedisCacher) Incr(key string) error {
	if !c.IsExist(key) {
		return fmt.Errorf("key '%s' not exist", key)
	}
	return c.c.Incr(graceful.GetManager().HammerContext(), c.prefix+key).Err()
}

// Decr decreases cached int-type value by given key as a counter.
func (c *RedisCacher) Decr(key string) error {
	if !c.IsExist(key) {
		return fmt.Errorf("key '%s' not exist", key)
	}
	return c.c.Decr(graceful.GetManager().HammerContext(), c.prefix+key).Err()
}

// IsExist returns true if cached value exists.
func (c *RedisCacher) IsExist(key string) bool {
	if c.c.Exists(graceful.GetManager().HammerContext(), c.prefix+key).Val() == 1 {
		return true
	}

	if !c.occupyMode {
		c.c.HDel(graceful.GetManager().HammerContext(), c.hsetName, c.prefix+key)
	}
	return false
}

// Flush deletes all cached data.
func (c *RedisCacher) Flush() error {
	if c.occupyMode {
		return c.c.FlushDB(graceful.GetManager().HammerContext()).Err()
	}

	keys, err := c.c.HKeys(graceful.GetManager().HammerContext(), c.hsetName).Result()
	if err != nil {
		return err
	}
	if err = c.c.Del(graceful.GetManager().HammerContext(), keys...).Err(); err != nil {
		return err
	}
	return c.c.Del(graceful.GetManager().HammerContext(), c.hsetName).Err()
}

// StartAndGC starts GC routine based on config string settings.
// AdapterConfig: network=tcp,addr=:6379,password=macaron,db=0,pool_size=100,idle_timeout=180,hset_name=MacaronCache,prefix=cache:
func (c *RedisCacher) StartAndGC(opts cache.Options) error {
	c.hsetName = "MacaronCache"
	c.occupyMode = opts.OccupyMode

	uri := nosql.ToRedisURI(opts.AdapterConfig)

	c.c = nosql.GetManager().GetRedisClient(uri.String())

	for k, v := range uri.Query() {
		switch k {
		case "hset_name":
			c.hsetName = v[0]
		case "prefix":
			c.prefix = v[0]
		}
	}

	return c.c.Ping(graceful.GetManager().HammerContext()).Err()
}

func init() {
	cache.Register("redis", &RedisCacher{})
}
