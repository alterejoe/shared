package create

import (
	"clerk-portal/shared/interfaces"

	"github.com/gomodule/redigo/redis"
)

func GetRedisPool(auth interfaces.RedisAuth) *redis.Pool {
	pool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
	return pool
}
