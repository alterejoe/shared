package interfaces

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// token access for user login flow
type OIDCClientCreds interface {
	GetDomain() string
	GetClientID() string
	GetClientSecret() string
	GetCallbackURL() string
	GetScopes() []string
}

// user login flow
type OIDCAuthenticator interface { // rename from UserLoginFlow to clarify purpose
	VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error)
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	GenerateRandomState() (string, error)
}

// type TokenGrantClient interface {
// 	GetDomain() string
// 	GetClientID() string
// 	GetClientSecret() string
// 	GetTokenAudience() string
// 	GetScopes() []string
// }
//
// // For verifying incoming tokens in middleware (JWT validation)
// type TokenValidator interface {
// 	GetDomain() string
// 	GetMiddlewareAudience() string
// }

// type JWTClientConfig interface {
// 	// GetDomain() string
// 	// GetClientID() string
// 	// GetClientSecret() string
// 	// GetTokenAudience() string
// 	// GetScopes() []string
// 	// GetJwtToken() (string, error) // basically the same as excahnge in Oidc
// }

// type JWTClientCreds interface {
// 	GetDomain() string
// 	GetClientID() string
// 	GetClientSecret() string
// 	GetTokenAudience() string
// 	GetScopes() []string
// }
// type JWTAuthenticator interface {
// 	GetDomain() string
// 	GetAudience() string
// 	GetScopes() []string
// }

type JWTIssuer interface {
	GetJwtToken() (string, error)
}

type JWTValidator interface {
	ValidateToken(ctx context.Context, token string) (any, error)
}
