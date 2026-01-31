package create

import (
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
)

func GetRedisSessionManager(pool *redis.Pool) *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Store = redisstore.New(pool)
	sessionManager.Lifetime = 12 * time.Hour
	// sessionManager.IdleTimeout = 5 * time.Second

	return sessionManager
}
