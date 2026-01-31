package env

import "os"

type RedisENV struct {
}

func (auth *RedisENV) GetHost() string {
	return os.Getenv("REDIS_HOST")
}

func (auth *RedisENV) GetPort() string {
	return os.Getenv("REDIS_PORT")
}
