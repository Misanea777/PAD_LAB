package db

import (
	"os"
	"time"

	"github.com/PAD_LAB/validators"
	"github.com/go-redis/redis/v7"
)

var RedisClient *redis.Client

func InitRedis() {
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "goRedis:6379"
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         dsn,
		PoolSize:     20,
		MinIdleConns: 10,
		Password:     "password",
	})
	_, err := RedisClient.Ping().Result()
	if err != nil {
		panic(err.Error())
	}
}

func SaveToken(ID string, td *validators.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := RedisClient.Set(td.AccessUuid, ID, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := RedisClient.Set(td.RefreshUuid, ID, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func FetchAuth(authD *validators.AccessDetails) (string, error) {
	userid, err := RedisClient.Get(authD.AccessUuid).Result()
	if err != nil {
		return "", err
	}

	return userid, nil
}
