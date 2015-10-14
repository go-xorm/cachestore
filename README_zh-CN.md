简介
======
cachestore 是一个用于 `xorm` 的缓存包，目前支持的缓存引擎有：goleveldb、memcache和redis，后续我们还会根据需要添加更多其它的存储引擎。
同时，也欢迎大家和我们一起来完善它，尽量对更多的缓存引擎实现支持。

使用
======

## 引子
`xorm` 默认提供了基于内存的缓存支持，一般我们会这样使用它：

	package main

	import (
		_ "github.com/go-sql-driver/mysql"
		"github.com/go-xorm/xorm"
	)

	func main() {
		engine, err := xorm.NewEngine("mysql", "username:password@/database?charset=utf8")
		if err != nil {
			panic(err)
		}
		ccStore := xorm.NewMemoryStore()
		cacher := xorm.NewLRUCacher(ccStore, 1000)
		//也可以指定缓存有效时间，如下
		//cacher.Expired = 86400 * time.Second
		engine.SetDefaultCacher(cacher)
		
		
		//下面开始执行数据库查询
		//engine.Where(querystring, ...)
	
	}

## 实战

### 1、使用goleveldb作为缓存

	package main

	import (
		_ "github.com/go-sql-driver/mysql"
		"github.com/go-xorm/cachestore"
		"github.com/go-xorm/xorm"
	)

	func main() {
		engine, err := xorm.NewEngine("mysql", "username:password@/database?charset=utf8")
		if err != nil {
			panic(err)
		}
		storagePath := "data/leveldb/dbcache"
		ccStore := cachestore.NewLevelDBStore(storagePath)
		ccStore := xorm.NewMemoryStore()
		cacher := xorm.NewLRUCacher(ccStore, 99999999)
		engine.SetDefaultCacher(cacher)
		
		//下面开始执行数据库查询
		//engine.Where(querystring, ...)
	
	}

### 2、使用memcache作为缓存

	package main

	import (
		_ "github.com/go-sql-driver/mysql"
		"github.com/go-xorm/cachestore"
		"github.com/go-xorm/xorm"
	)

	func main() {
		engine, err := xorm.NewEngine("mysql", "username:password@/database?charset=utf8")
		if err != nil {
			panic(err)
		}
		configs := []string{
			"127.0.0.1:11211",
			"192.168.1.2:11211",
		}
		ccStore := cachestore.NewMemCache(configs)
		cacher := xorm.NewLRUCacher(ccStore, 99999999)
		engine.SetDefaultCacher(cacher)
		
		//下面开始执行数据库查询
		//engine.Where(querystring, ...)
	
	}


### 3、使用redis作为缓存

	package main

	import (
		_ "github.com/go-sql-driver/mysql"
		"github.com/go-xorm/cachestore"
		"github.com/go-xorm/xorm"
	)

	func main() {
		engine, err := xorm.NewEngine("mysql", "username:password@/database?charset=utf8")
		if err != nil {
			panic(err)
		}
		configs := map[string]string{
			"conn":"localhost:6379",
			"key":"default", // the collection name of redis for cache adapter.
		}
		ccStore := cachestore.NewRedisCache(configs)
		cacher := xorm.NewLRUCacher(ccStore, 99999999)
		engine.SetDefaultCacher(cacher)
		
		//下面开始执行数据库查询
		//engine.Where(querystring, ...)
	
	}


### 4、也可以将以上缓存引擎结合起来，按需切换
	实例代码详见本包中的文件`example/main.go`。
	
#EOF