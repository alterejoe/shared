package create

import (
	"github.com/alterejoe/shared/structs"

	"github.com/gomodule/redigo/redis"
)

func CreateRedisPool(auth structs.RedisAuth) *redis.Pool {
	pool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
	return pool
}
