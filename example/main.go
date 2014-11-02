package main

import (
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/cachestore"
	//"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

var (
	cacher   *xorm.LRUCacher
	CacheDir string   = "."
	cfg      []string = []string{"leveldb"}
)

func main() {

	//cfg := strings.SplitN(cacherName, ":", 2)

	engine, err := xorm.NewEngine("mysql", "root:root@/coscms?charset=utf8")
	if err != nil {
		log.Fatalf("The database connection failed: %v\n", err)
	}
	switch strings.ToLower(cfg[0]) {
	case "memory":
		ccStore := xorm.NewMemoryStore()
		//this.CacheStore = ccStore
		cacher = xorm.NewLRUCacher(ccStore, 1000) //NewLRUCacher(store core.CacheStore, maxElementSize int)
	case "leveldb":
		storagePath := CacheDir + "/leveldb/dbcache"
		if len(cfg) == 2 {
			storagePath = cfg[1]
		}
		ccStore := cachestore.NewLevelDBStore(storagePath)
		//this.CacheStore = ccStore
		//ccStore.Debug = this.ShowDebug
		cacher = xorm.NewLRUCacher(ccStore, 999999999)
	case "memcache":
		conn := "127.0.0.1:11211"
		if len(cfg) == 2 {
			conn = cfg[1]
		}
		ccStore := cachestore.NewMemCache(strings.Split(conn, ";"))
		//this.CacheStore = ccStore
		//ccStore.Debug = this.ShowDebug
		cacher = xorm.NewLRUCacher(ccStore, 999999999)
	}
	if cacher != nil {
		cacher.Expired = 86400 * time.Second
		engine.SetDefaultCacher(cacher)
	}

	//engine.Where(querystring, ...)
}
