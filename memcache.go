package cachestore

import (
	"coscms/app/base/lib/cachestore/memcache"
	"encoding/json"
	"errors"
	"fmt"
)

// Memcache adapter.
type MemcacheCache struct {
	c        *memcache.Connection
	conninfo string
	LifeTime uint64
	Debug    bool
}

// create new memcache adapter.
func NewMemCache(conn string) *MemcacheCache {
	rc := &MemcacheCache{LifeTime: 86400}
	rc.conninfo = conn
	rc.c = rc.connectInit()
	return rc
}

// get value from memcache.
func (rc *MemcacheCache) Get(key string) (interface{}, error) {
	if rc.c == nil {
		rc.c = rc.connectInit()
		if rc.c == nil {
			return nil, nil
		}
	}
	val, err := rc.c.Get(Md5(key))
	if err != nil || len(val) < 1 {
		if err != nil {
			fmt.Println("[Memcache]GetErr: ", err)
			rc.c.Close()
			rc.c = nil
		}
		return nil, err
	}

	var v interface{}
	err = Decode(val[0].Value, &v)
	if err != nil {
		fmt.Println("[Memcache]DecodeErr: ", err)
	}

	if rc.Debug {
		fmt.Println("[Memcache]Get: ", key)
	}
	return v, err
}

// put value to memcache. only support string.
func (rc *MemcacheCache) Put(key string, value interface{}) error {
	if rc.c == nil {
		rc.c = rc.connectInit()
		if rc.c == nil {
			return nil
		}
	}

	val, err := Encode(value)
	if err != nil {
		fmt.Println("[Memcache]EncodeErr: ", err)
		return err
	}

	stored, err := rc.c.Set(Md5(key), 0, rc.LifeTime, val)
	if err != nil || stored == false {
		if err != nil {
			fmt.Println("[Memcache]PutErr: ", err)
			rc.c.Close()
			rc.c = nil
		}
		return errors.New("stored fail")
	}
	if rc.Debug {
		fmt.Println("[Memcache]Put: ", key)
	}
	return err
}

// delete value in memcache.
func (rc *MemcacheCache) Del(key string) error {
	if rc.c == nil {
		rc.c = rc.connectInit()
		if rc.c == nil {
			return nil
		}
	}
	_, err := rc.c.Delete(Md5(key))
	if err != nil {
		fmt.Println("[Memcache]DelErr: ", err)
		rc.c.Close()
		rc.c = nil
	}
	if rc.Debug {
		fmt.Println("[Memcache]Del: ", key)
	}
	return err
}

// [Not Support]
// increase counter.
func (rc *MemcacheCache) Incr(key string) error {
	return errors.New("not support in memcache")
}

// [Not Support]
// decrease counter.
func (rc *MemcacheCache) Decr(key string) error {
	return errors.New("not support in memcache")
}

// check value exists in memcache.
func (rc *MemcacheCache) IsExist(key string) bool {
	if rc.c == nil {
		rc.c = rc.connectInit()
	}
	v, err := rc.c.Get(key)
	if err != nil {
		rc.c.Close()
		rc.c = nil
		return false
	}
	if len(v) == 0 {
		return false
	}
	return true
}

// clear all cached in memcache.
func (rc *MemcacheCache) ClearAll() error {
	if rc.c == nil {
		rc.c = rc.connectInit()
		if rc.c == nil {
			return nil
		}
	}
	err := rc.c.FlushAll()
	return err
}

// start memcache adapter.
// config string is like {"conn":"connection info"}.
// if connecting error, return.
func (rc *MemcacheCache) Connect(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["conn"]; !ok {
		return errors.New("config has no conn key")
	}
	rc.conninfo = cf["conn"]
	rc.c = rc.connectInit()
	if rc.c == nil {
		return errors.New("dial tcp conn error")
	}
	return nil
}

// connect to memcache and keep the connection.
func (rc *MemcacheCache) connectInit() *memcache.Connection {
	c, err := memcache.Connect(rc.conninfo)
	if err != nil {
		fmt.Println("[Memcahe]Connect failure:", err)
		return nil
	}
	if rc.Debug {
		fmt.Println("[Memcahe]Connect success:", rc.conninfo)
	}
	return c
}
