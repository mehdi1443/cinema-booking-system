package database

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	"github.com/redis/go-redis/v9"
)

type Cache interface {
	GetRDB() *redis.Client
	GetRedisConstant() *bootstrap.RedisConst
}

type RedisDatabase struct {
	rdb       *redis.Client
	rdb_const *bootstrap.RedisConst
}

var (
	rdbOnce     sync.Once
	rdbInstance *RedisDatabase
)

func NewRedisDatabase(redisConfig *bootstrap.Redis, redisConst *bootstrap.RedisConst) *RedisDatabase {
	rdbOnce.Do(func() {
		rdbNumber, _ := strconv.Atoi(redisConfig.RDBNumber)
		address := fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port)
		rdb := redis.NewClient(&redis.Options{
			Addr:     address,
			Password: redisConfig.Password,
			DB:       rdbNumber,
		})
		_, err := rdb.Ping(context.Background()).Result()
		if err != nil {
			log.Fatalf("Critical error: failed connecting to Redis at %s: %v", address, err)
		}
		rdbInstance = &RedisDatabase{rdb: rdb, rdb_const: redisConst}
	})

	return rdbInstance
}

func (rdb *RedisDatabase) GetRDB() *redis.Client {
	return rdbInstance.rdb
}

func (rdb *RedisDatabase) GetRedisConstant() *bootstrap.RedisConst {
	return rdb.rdb_const
}
