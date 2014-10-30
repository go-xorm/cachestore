package cachestore

import (
	"fmt"
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
		fmt.Println("[LevelDB]EncodeErr: ", err)
		return err
	}
	err = s.store.Put([]byte(key), val, nil)
	if err != nil {
		fmt.Println("[LevelDB]PutErr: ", err)
	}
	if s.Debug {
		fmt.Println("[LevelDB]Put: ", key)
	}
	return err
}

func (s *LevelDBStore) Get(key string) (interface{}, error) {
	data, err := s.store.Get([]byte(key), nil)
	if err != nil {
		fmt.Println("[LevelDB]GetErr: ", err)
		return nil, err
	}

	err = Decode(data, &s.v)
	if err != nil {
		fmt.Println("[LevelDB]DecodeErr: ", err)
	}
	if s.Debug {
		fmt.Println("[LevelDB]Get: ", key, s.v)
	}
	return s.v, err
}

func (s *LevelDBStore) Del(key string) error {
	err := s.store.Delete([]byte(key), nil)
	if err != nil {
		fmt.Println("[LevelDB]DelErr: ", err)
		return err
	}
	if s.Debug {
		fmt.Println("[LevelDB]Del: ", key)
	}
	return err
}

func (s *LevelDBStore) Close() {
	s.store.Close()
}
