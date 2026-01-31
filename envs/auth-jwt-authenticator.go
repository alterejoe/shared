package env

import "fmt"

type JwtValidatorEnv struct {
	Domain, Audience string
	Scope            []string
}

func NewJwtValidatorENV(prefix string) (*JwtValidatorEnv, error) {
	required := []string{"DOMAIN", "AUDIENCE", "SCOPE"}
	for _, key := range required {
		if getEnv(prefix, key) == "" {
			return nil, fmt.Errorf("missing env var %s_%s", prefix, key)
		}
	}
	return &JwtValidatorEnv{
		Domain:   getEnv(prefix, "DOMAIN"),
		Audience: getEnv(prefix, "AUDIENCE"),
		Scope:    cleanScopes(getEnv(prefix, "SCOPE")),
	}, nil
}

func (v *JwtValidatorEnv) GetDomain() string   { return v.Domain }
func (v *JwtValidatorEnv) GetAudience() string { return v.Audience }
func (v *JwtValidatorEnv) GetScope() []string  { return v.Scope }
