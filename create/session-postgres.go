package create

import (
	"context"
	"encoding/gob"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/alterejoe/shared/structs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserPreferences struct {
	FleetingAlign string `json:"fleeting_align"`
}

func DefaultPreferences() UserPreferences {
	return UserPreferences{
		FleetingAlign: "Top",
	}
}

func PGSessionManager(pool *pgxpool.Pool, cookiename string) *SessionManager {
	gob.Register(UserPreferences{})

	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(pool)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.IdleTimeout = 60 * time.Minute
	sessionManager.Cookie.Name = cookiename

	return &SessionManager{
		SessionManager: sessionManager,
	}
}

type SessionManager struct {
	*scs.SessionManager
}

// Flash messages

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

	typeStr := t.(string)
	msgStr := msg.(string)

	d.Put(r, "flashmsg-type", "")
	d.Put(r, "flashmsg", "")

	return typeStr, msgStr
}

// Auth

func (d *SessionManager) GetAuthUserID(ctx context.Context) pgtype.UUID {
	uid := d.Get(ctx, "authenticatedUserID").(string)
	useruuid := uuid.MustParse(uid)

	return pgtype.UUID{
		Bytes: useruuid,
		Valid: true,
	}
}

func (d *SessionManager) SetAuthUser(ctx context.Context, user structs.User) {
	d.Put(ctx, "authenticatedUserID", user.ID)
	d.Put(ctx, "user_name", user.Name)
	d.Put(ctx, "user_email", user.Email)
}

func (d *SessionManager) DeleteAuthUser(ctx context.Context) {
	d.Remove(ctx, "authenticatedUserID")
	d.Remove(ctx, "user_name")
	d.Remove(ctx, "user_email")
}

// Preferences

func (d *SessionManager) GetPreferences(ctx context.Context) UserPreferences {
	prefs := d.Get(ctx, "user_preferences")
	if prefs == nil {
		return DefaultPreferences()
	}

	p, ok := prefs.(UserPreferences)
	if !ok {
		return DefaultPreferences()
	}

	return p
}

func (d *SessionManager) SetPreference(ctx context.Context, key string, value string) {
	prefs := d.GetPreferences(ctx)

	switch key {
	case "fleeting_align":
		prefs.FleetingAlign = value
	}

	d.Put(ctx, "user_preferences", prefs)
}

func (d *SessionManager) SetPreferences(ctx context.Context, prefs UserPreferences) {
	d.Put(ctx, "user_preferences", prefs)
}
