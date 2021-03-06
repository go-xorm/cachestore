package cachestore

import (
	"log"

	"github.com/syndtr/goleveldb/leveldb"
	//"reflect"
)

// LevelDBStore implements CacheStore provide local machine
type LevelDBStore struct {
	store *leveldb.DB
	Debug bool
	v     interface{}
}

func NewLevelDBStore(dbfile string) *LevelDBStore {
	db := &LevelDBStore{}
	if h, err := leveldb.OpenFile(dbfile, nil); err != nil {
		panic(err)
	} else {
		db.store = h
	}
	return db
}

func (s *LevelDBStore) Put(key string, value interface{}) error {
	val, err := Encode(value)
	if err != nil {
		if s.Debug {
			log.Println("[LevelDB]EncodeErr: ", err, "Key:", key)
		}
		return err
	}
	err = s.store.Put([]byte(key), val, nil)
	if err != nil {
		if s.Debug {
			log.Println("[LevelDB]PutErr: ", err, "Key:", key)
		}
		return err
	}
	if s.Debug {
		log.Println("[LevelDB]Put: ", key)
	}
	return err
}

func (s *LevelDBStore) Get(key string) (interface{}, error) {
	data, err := s.store.Get([]byte(key), nil)
	if err != nil {
		if s.Debug {
			log.Println("[LevelDB]GetErr: ", err, "Key:", key)
		}
		return nil, err
	}

	err = Decode(data, &s.v)
	if err != nil {
		if s.Debug {
			log.Println("[LevelDB]DecodeErr: ", err, "Key:", key)
		}
		return nil, err
	}
	if s.Debug {
		log.Println("[LevelDB]Get: ", key, s.v)
	}
	return s.v, err
}

func (s *LevelDBStore) Del(key string) error {
	err := s.store.Delete([]byte(key), nil)
	if err != nil {
		if s.Debug {
			log.Println("[LevelDB]DelErr: ", err, "Key:", key)
		}
		return err
	}
	if s.Debug {
		log.Println("[LevelDB]Del: ", key)
	}
	return err
}

func (s *LevelDBStore) Close() {
	s.store.Close()
}
