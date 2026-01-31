package interfaces

import (
	auth0_api "clerk-portal/auth0_admin_api/db"
	admin "clerk-portal/runtime_admin/db"
	client "clerk-portal/runtime_client/db"
	"context"

	"github.com/jackc/pgx/v5"
)

type CommonContext interface {
	Logger() Logger
	Session() SessionManager
	Sanitizer() Sanitizer
	Response
}

type DBQueryContext[T any] interface {
	Close()
	BeginTx(r context.Context, withRLS, system bool) (pgx.Tx, error)
	TxWithCleanup(r context.Context) (pgx.Tx, func(success bool), error)
	TxNoRLSWithCleanup(r context.Context) (pgx.Tx, func(success bool), error)
	TxSystemRLSWithCleanup(r context.Context) (pgx.Tx, func(success bool), error)
	Tx(context.Context) (pgx.Tx, error)
	TxNoRLS(context.Context) (pgx.Tx, error)
	TxSystemRLS(context.Context) (pgx.Tx, error)
	SetRLS(r context.Context, tx pgx.Tx) error

	NewQueries(pgx.Tx) T
	QueryTx(r context.Context, params StandardSQLCQuery[T], tx pgx.Tx) (any, error)
}

type ClientDBContext = DBQueryContext[*client.Queries]
type AdminDBContext = DBQueryContext[*admin.Queries]
type AdminAuth0ApiDBContext = DBQueryContext[*auth0_api.Queries]

type ClientUserContext interface {
	CommonContext
	ClientDBContext
}

type AdminUserContext interface {
	CommonContext
	AdminDBContext
}

type ClientAuthContext interface {
	CommonContext
	ClientDBContext
	OidcAuthenticator() OIDCAuthenticator
}

type AdminAuthContext interface {
	CommonContext
	AdminDBContext
	OidcAuthenticator() OIDCAuthenticator
	JwtIssuer() JWTIssuer
}

type Auth0ApiAuthContext interface {
	CommonContext
	AdminAuth0ApiDBContext
	JwtIssuer() JWTIssuer
	JwtValidator() JWTValidator
}
