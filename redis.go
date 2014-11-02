package cachestore

import (
	"github.com/go-xorm/cachestore/redigo/redis"
	"encoding/json"
	"errors"
	"log"
)

var (
	// the collection name of redis for cache adapter.
	DefaultKey string = "comcmsRedis"
)

// Redis cache adapter.
type RedisCache struct {
	c        redis.Conn
	conninfo string
	key      string
	Debug    bool
}

// create new redis cache with default collection name.
func NewRedisCache(cf map[string]string) *RedisCache {
	rc := &RedisCache{key: DefaultKey}
	if _, ok := cf["key"]; !ok {
		cf["key"] = DefaultKey
	}
	rc.key = cf["key"]
	rc.conninfo = cf["conn"]
	var err error
	rc.c, err = rc.connectInit()
	if err != nil {
		rc.c = nil
	}
	return rc
}

// Get cache from redis.
func (rc *RedisCache) Get(key string) (interface{}, error) {
	if rc.c == nil {
		var err error
		rc.c, err = rc.connectInit()
		if err != nil {
			return nil, err
		}
	}
	val, err := rc.c.Do("HGET", rc.key, key)
	if err != nil {
		if rc.Debug {
			log.Println("[Redis]GetErr: ", err, "Key:", key)
		}
		return nil, err
	}
	var v interface{}
	err = Decode(val.([]byte), &v)
	if err != nil {
		if rc.Debug {
			log.Println("[Redis]DecodeErr: ", err, "Key:", key)
		}
		return nil, err
	}
	if rc.Debug {
		log.Println("[Redis]Get: ", key)
	}
	return v, err
}

// put cache to redis.
// timeout is ignored.
func (rc *RedisCache) Put(key string, value interface{}) error {
	if rc.c == nil {
		var err error
		rc.c, err = rc.connectInit()
		if err != nil {
			return err
		}
	}
	val, err := Encode(value)
	if err != nil {
		if rc.Debug {
			log.Println("[Redis]EncodeErr: ", err, "Key:", key)
		}
		return err
	}
	_, err = rc.c.Do("HSET", rc.key, key, val)
	if err != nil {
		if rc.Debug {
			log.Println("[Redis]PutErr: ", err, "Key:", key)
		}
		return err
	}
	if rc.Debug {
		log.Println("[Redis]Put: ", key)
	}
	return err
}

// delete cache in redis.
func (rc *RedisCache) Del(key string) error {
	if rc.c == nil {
		var err error
		rc.c, err = rc.connectInit()
		if err != nil {
			return err
		}
	}
	_, err := rc.c.Do("HDEL", rc.key, key)
	if err != nil {
		if rc.Debug {
			log.Println("[Redis]DelErr: ", err, "Key:", key)
		}
		return err
	}
	if rc.Debug {
		log.Println("[Redis]Del: ", key)
	}
	return err
}

// check cache exist in redis.
func (rc *RedisCache) IsExist(key string) bool {
	if rc.c == nil {
		var err error
		rc.c, err = rc.connectInit()
		if err != nil {
			return false
		}
	}
	v, err := redis.Bool(rc.c.Do("HEXISTS", rc.key, key))
	if err != nil {
		return false
	}
	return v
}

// increase counter in redis.
func (rc *RedisCache) Incr(key string, delta uint64) error {
	if rc.c == nil {
		var err error
		rc.c, err = rc.connectInit()
		if err != nil {
			return err
		}
	}
	_, err := redis.Bool(rc.c.Do("HINCRBY", rc.key, key, delta))
	if err != nil {
		return err
	}
	return nil
}

// decrease counter in redis.
func (rc *RedisCache) Decr(key string, delta uint64) error {
	if rc.c == nil {
		var err error
		rc.c, err = rc.connectInit()
		if err != nil {
			return err
		}
	}
	_, err := redis.Bool(rc.c.Do("HINCRBY", rc.key, key, -delta))
	if err != nil {
		return err
	}
	return nil
}

// clean all cache in redis. delete this redis collection.
func (rc *RedisCache) ClearAll() error {
	if rc.c == nil {
		var err error
		rc.c, err = rc.connectInit()
		if err != nil {
			return err
		}
	}
	_, err := rc.c.Do("DEL", rc.key)
	return err
}

// start redis cache adapter.
// config is like {"key":"collection key","conn":"connection info"}
// the cache item in redis are stored forever,
// so no gc operation.
func (rc *RedisCache) Connect(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["key"]; !ok {
		cf["key"] = DefaultKey
	}
	if _, ok := cf["conn"]; !ok {
		return errors.New("config has no conn key")
	}
	rc.key = cf["key"]
	rc.conninfo = cf["conn"]
	var err error
	rc.c, err = rc.connectInit()
	if err != nil {
		return err
	}
	if rc.c == nil {
		return errors.New("dial tcp conn error")
	}
	return nil
}

// connect to redis.
func (rc *RedisCache) connectInit() (redis.Conn, error) {
	c, err := redis.Dial("tcp", rc.conninfo)
	if err != nil {
		return nil, err
	}
	return c, nil
}
