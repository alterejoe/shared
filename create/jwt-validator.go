package create

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// ServiceJWTValidator validates self-signed JWTs using a shared HMAC-SHA256 secret.
// Same usage pattern as JWTValidator but no Auth0/JWKS dependency.
type ServiceJWTValidator struct {
	secret   []byte
	audience string
	issuer   string
}

// NewServiceJWTValidator creates a validator from the shared secret.
// Reads the same secret the issuer signs with.
func NewServiceJWTValidator(secret, audience, issuer string) *ServiceJWTValidator {
	return &ServiceJWTValidator{
		secret:   []byte(secret),
		audience: audience,
		issuer:   issuer,
	}
}

// ValidateToken parses and validates a JWT string.
// Matches the same method signature as JWTValidator.ValidateToken().
// Returns the parsed claims or an error if invalid/expired.
func (v *ServiceJWTValidator) ValidateToken(ctx context.Context, tokenStr string) (interface{}, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&ServiceClaims{},
		func(t *jwt.Token) (interface{}, error) {
			// Enforce HMAC signing — reject RS256 or any other method
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return v.secret, nil
		},
		jwt.WithIssuer(v.issuer),
		jwt.WithAudience(v.audience),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*ServiceClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
