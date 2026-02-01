package create

import (
	"context"
	"net/url"
	"time"

	envs "github.com/alterejoe/envs"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

type JWTValidator struct {
	validator *validator.Validator
}

func (v *JWTValidator) ValidateToken(ctx context.Context, token string) (interface{}, error) {
	return v.validator.ValidateToken(ctx, token)
}

type CustomClaims struct {
	Scope string `json:"scope"`
}

func (c *CustomClaims) Validate(ctx context.Context) error {
	// optional: enforce required claims
	// if c.Scope == "" { return errors.New("missing scope") }

	return nil
}
func JwtValidator(targetjwtprefix string) *JWTValidator {
	e, err := envs.NewJwtValidatorENV(targetjwtprefix)
	if err != nil {
		panic(err)
	}

	issuer := "https://" + e.GetDomain() + "/"

	u, err := url.Parse(issuer)
	if err != nil {
		panic(err)
	}

	provider := jwks.NewCachingProvider(u, 5*time.Minute)

	v, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuer,
		[]string{e.GetAudience()},
		validator.WithCustomClaims(func() validator.CustomClaims { return &CustomClaims{} }),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		panic(err)
	}

	return &JWTValidator{validator: v}
}
