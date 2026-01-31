package interfaces

import (
	// auth0_api "clerk-portal/auth0_admin_api/db"
	// admin "clerk-portal/runtime_admin/db"
	// client "clerk-portal/runtime_client/db"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RLSMode defines the Row Level Security mode for transactions
type RLSMode int

const (
	RLSUser   RLSMode = iota // Use authenticated user's RLS context
	RLSSystem                // Use system-level RLS context
	RLSNone                  // No RLS (disabled)
)

// CommonContext provides shared functionality across all contexts
type CommonContext interface {
	Logger() Logger
	Session() SessionManager
	Sanitizer() Sanitizer
	Response
}

// DBQueryContext provides database query operations with RLS support
// T should be your project's generated sqlc Queries type (e.g., *db.Queries)
type DBQueryContext[T any] interface {
	GetDBPool() *pgxpool.Pool
	SetRLS(ctx context.Context, tx pgx.Tx) error
	SetSystemRLS(ctx context.Context, tx pgx.Tx) error
	WithTx(ctx context.Context, mode RLSMode, fn func(T) error) error
}

// UserContext combines common functionality with database access
type UserContext[T any] interface {
	CommonContext
	DBQueryContext[T]
}

// AuthContext adds authentication capabilities to UserContext
type AuthContext[T any] interface {
	UserContext[T]
	OidcAuthenticator() OIDCAuthenticator
}

// AdminAuthContext adds JWT issuing capability for admin contexts
type AdminAuthContext[T any] interface {
	AuthContext[T]
	JwtIssuer() JWTIssuer
}

// ApiAuthContext provides both JWT issuing and validation
type ApiAuthContext[T any] interface {
	UserContext[T]
	JwtIssuer() JWTIssuer
	JwtValidator() JWTValidator
}
