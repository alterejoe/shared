package interfaces

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type User interface {
	ID() uuid.UUID
	Name() string
	Email() string
}

type Sanitizer any

type Logger interface {
	Debug(string, ...any)
	Info(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
}

type SessionManager interface {
	LoadAndSave(http.Handler) http.Handler
	Get(r context.Context, key string) any
	Put(r context.Context, key string, value any)
	Remove(r context.Context, key string)
	Destroy(ctx context.Context) error
	Token(ctx context.Context) string
	SetFlashMessage(r context.Context, t, msg string)
	GetFlashMessage(r context.Context) (string, string)
}

type Response interface {
	ServerSuccess(w http.ResponseWriter, r *http.Request)
	ServerCreated(w http.ResponseWriter, r *http.Request)
	ServerError(w http.ResponseWriter, r *http.Request, err error)
	ClientError(w http.ResponseWriter, r *http.Request, err error)
}
