package create

import (
	"context"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/alterejoe/shared/structs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func PGSessionManager(pool *pgxpool.Pool, cookiename string) *SessionManager {
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(pool)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.IdleTimeout = 60 * time.Minute
	sessionManager.Cookie.Name = cookiename
	// sessionManager.IdleTimeout = 10 * time.Second

	return &SessionManager{
		SessionManager: sessionManager,
	}
}

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

func (d *SessionManager) SetElectionID(ctx context.Context, electionID uuid.UUID) {
	d.Put(ctx, "electionID", electionID.String())
}

func (d *SessionManager) GetElectionID(ctx context.Context) uuid.UUID {
	electionID := d.Get(ctx, "electionID")
	electionIDstring := electionID.(string)
	electionid, err := uuid.Parse(electionIDstring)
	if err != nil {
		panic(err)
	}
	return electionid
}

func (d *SessionManager) GetGroupName(ctx context.Context) string {
	return d.Get(ctx, "group_name").(string)
}

func (d *SessionManager) GetClientGroupID(ctx context.Context) uuid.UUID {
	sg := d.Get(ctx, "authenticatedClientGroup")
	sgstring := sg.(string)
	clientgroup, err := uuid.Parse(sgstring)
	if err != nil {
		panic(err)
	}
	return clientgroup
}
func (d *SessionManager) GetAdminGroupID(ctx context.Context) uuid.UUID {
	// return d.Get(ctx, "authenticatedAdminGroup").(uuid.UUID)
	sg := d.Get(ctx, "authenticatedAdminGroup")
	sgstring := sg.(string)
	admingroup, err := uuid.Parse(sgstring)
	if err != nil {
		panic(err)
	}
	return admingroup
}
