package create

import (
	"context"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionManager struct {
	*scs.SessionManager
}

func (d *SessionManager) SetFlashMessage(r context.Context, t, msg string) {
	d.Put(r, "flashmsg-type", t)
	d.Put(r, "flashmsg", msg)
}

func (d *SessionManager) GetFlashMessage(r context.Context) (string, string) {
	t := d.Get(r, "flashmsg-type")
	msg := d.Get(r, "flashmsg")

	if t == nil || msg == nil {
		return "", ""
	}
	var typeStr string
	typeStr = t.(string)
	var msgStr string
	msgStr = msg.(string)

	d.Put(r, "flashmsg-type", "")
	d.Put(r, "flashmsg", "")

	return typeStr, msgStr
}

func GetPGSessionManager(pool *pgxpool.Pool) *SessionManager {
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(pool)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.IdleTimeout = 60 * time.Minute
	// sessionManager.IdleTimeout = 10 * time.Second

	return &SessionManager{
		SessionManager: sessionManager,
	}
}
