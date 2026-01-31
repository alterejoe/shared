package env

import (
	"github.com/coreos/go-oidc/v3/oidc"
)

type OIDCEnvConfig struct {
	Domain       string
	ClientID     string
	ClientSecret string
	CallbackURL  string
	Scopes       []string
}

// Build config using a prefix like "AUTH0_ADMIN_UI" or "AUTH0_CLIENT_UI"
func NewOIDCEnvConfig(prefix string) (*OIDCEnvConfig, error) {
	// required := []string{"DOMAIN", "CLIENT_ID", "CLIENT_SECRET", "CALLBACK_URL"}
	// for _, key := range required {
	// 	if getEnv(prefix, key) == "" {
	// 		return nil, fmt.Errorf("missing env var %s_%s", prefix, key)
	// 	}
	// }
	if err := RequiredEnvs([]string{"DOMAIN", "CLIENT_ID", "CLIENT_SECRET", "CALLBACK_URL"}, prefix); err != nil {
		return nil, err
	}

	return &OIDCEnvConfig{
		Domain:       getEnv(prefix, "DOMAIN"),
		ClientID:     getEnv(prefix, "CLIENT_ID"),
		ClientSecret: getEnv(prefix, "CLIENT_SECRET"),
		CallbackURL:  getEnv(prefix, "CALLBACK_URL"),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}, nil
}

// Satisfy OIDCClientConfig interface
func (c *OIDCEnvConfig) GetDomain() string       { return c.Domain }
func (c *OIDCEnvConfig) GetClientID() string     { return c.ClientID }
func (c *OIDCEnvConfig) GetClientSecret() string { return c.ClientSecret }
func (c *OIDCEnvConfig) GetCallbackURL() string  { return c.CallbackURL }
func (c *OIDCEnvConfig) GetScopes() []string     { return c.Scopes }
