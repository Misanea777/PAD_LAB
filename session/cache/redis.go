package cache

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	"session/myTypes"

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
		dsn = "sessionCache:6379"
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

type ChunkKey struct {
	WordlId uint64
	PosX    int64
	PosY    int64
}

func (ck ChunkKey) toString() string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(int(ck.WordlId)))
	sb.WriteString("-")
	sb.WriteString(strconv.Itoa(int(ck.PosX)))
	sb.WriteString("-")
	sb.WriteString(strconv.Itoa(int(ck.PosY)))
	return sb.String()
}

func Store(key ChunkKey, obj myTypes.Chunk) error {
	stringKey := key.toString()
	if err := mycache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   stringKey,
		Value: obj,
		TTL:   time.Second * 3,
	}); err != nil {
		return err
	}

	return nil
}

func Get(key ChunkKey) (myTypes.Chunk, error) {
	stringKey := key.toString()
	var wanted myTypes.Chunk
	err := mycache.Get(ctx, stringKey, &wanted)
	return wanted, err
}
