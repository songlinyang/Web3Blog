package myredis

import (
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func InitRedis() *redis.Client {
	db, e := strconv.Atoi(os.Getenv("REDIS_DB"))
	if e != nil {
		panic(e)
	}
	zap.S().Debug("Redis DB:", db, "Redis Addr:", os.Getenv("REDIS_ADDR"))
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "", // no password set
		DB:       db, // use default DB
	})
	zap.S().Debug("Redis server:", rdb)
	return rdb
}
