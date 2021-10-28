package cache

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	state "session/game/gameState"
	mapSt "session/game/mapState"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

var mycache *cache.Cache
var ctx context.Context

func Init() {
	ctx = context.TODO()
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
	WordlId uint32
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

func StoreChnk(key ChunkKey, obj mapSt.Chunk) error {
	stringKey := key.toString()
	if err := mycache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   stringKey,
		Value: obj,
		TTL:   (time.Second * 3),
	}); err != nil {
		return err
	}

	return nil
}

func GetChnk(key ChunkKey) (mapSt.Chunk, error) {
	stringKey := key.toString()
	var wanted mapSt.Chunk
	err := mycache.Get(ctx, stringKey, &wanted)
	return wanted, err
}

func StoreSt(state *state.GameState) error {
	key := strconv.Itoa(int(state.Id))
	if err := mycache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: state,
		TTL:   (time.Hour * 100),
	}); err != nil {
		return err
	}

	return nil
}

func GetSt(key uint64) (state.GameState, error) {
	stringKey := strconv.Itoa(int(key))
	var wanted state.GameState
	err := mycache.Get(ctx, stringKey, &wanted)
	return wanted, err
}
