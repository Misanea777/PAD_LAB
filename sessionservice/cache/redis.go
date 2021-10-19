package cache

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

var mycache *cache.Cache
var ctx context.Context

func Init() {
	ctx = context.TODO()
	// ring := redis.NewRing(&redis.RingOptions{
	// 	Addrs: map[string]string{
	// 		"server1": ":6379",
	// 		"server2": ":6380",
	// 	},
	// })

	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:8081"
	}
	var RedisClient = redis.NewClient(&redis.Options{
		Addr: dsn,
	})
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		panic(err.Error())
	}

	mycache = cache.New(&cache.Options{
		Redis:      RedisClient,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

}

func Store(key string, obj interface{}) error {

	if err := mycache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: obj,
		TTL:   time.Hour,
	}); err != nil {
		return err
	}

	return nil
}

func Get(key string) error {
	var wanted interface{}
	if err := mycache.Get(ctx, key, &wanted); err == nil {
		return err
	}

	return nil
}
