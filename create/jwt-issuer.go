package create

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ServiceJWTIssuer issues self-signed JWTs using a shared HMAC-SHA256 secret.
// Same usage pattern as JWTIssuer but no Auth0 dependency.
type ServiceJWTIssuer struct {
	secret   []byte
	audience string
	issuer   string
	ttl      time.Duration

	token  string
	expiry time.Time
	mu     sync.Mutex
}

// ServiceClaims are the claims embedded in the self-signed JWT.
type ServiceClaims struct {
	jwt.RegisteredClaims
}

// NewServiceJWTIssuer creates an issuer from environment variables.
// Reads: <prefix>_JWT_SECRET, <prefix>_JWT_AUDIENCE, <prefix>_JWT_ISSUER
// TTL defaults to 1 hour if not set.
func NewServiceJWTIssuer(secret, audience, issuer string, ttl time.Duration) *ServiceJWTIssuer {
	if ttl == 0 {
		ttl = time.Hour
	}
	return &ServiceJWTIssuer{
		secret:   []byte(secret),
		audience: audience,
		issuer:   issuer,
		ttl:      ttl,
	}
}

// GetJwtToken returns a cached JWT, refreshing it if it's within 30 seconds of expiry.
// Matches the same method signature as JWTIssuer.GetJwtToken().
func (s *ServiceJWTIssuer) GetJwtToken() (string, error) {
	// Fast path — check without lock
	if s.token != "" && time.Until(s.expiry) > 30*time.Second {
		return s.token, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Re-check after acquiring lock
	if s.token != "" && time.Until(s.expiry) > 30*time.Second {
		return s.token, nil
	}

	if len(s.secret) == 0 {
		return "", fmt.Errorf("service JWT secret is empty")
	}

	now := time.Now()
	expiry := now.Add(s.ttl)

	claims := ServiceClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Audience:  jwt.ClaimStrings{s.audience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiry),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	s.token = signed
	s.expiry = expiry

	return signed, nil
}
