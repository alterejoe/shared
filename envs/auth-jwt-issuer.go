package env

type JwtIssuerEnv struct {
	Domain       string
	ClientID     string
	ClientSecret string
	Audience     string
	Scope        string
}

func NewJwtIssuerENV(prefix string) (*JwtIssuerEnv, error) {
	if err := RequiredEnvs([]string{"DOMAIN", "SCOPE", "AUDIENCE", "CLIENT_ID", "CLIENT_SECRET"}, prefix); err != nil {
		return nil, err
	}

	return &JwtIssuerEnv{
		Domain:       getEnv(prefix, "DOMAIN"),
		ClientID:     getEnv(prefix, "CLIENT_ID"),
		ClientSecret: getEnv(prefix, "CLIENT_SECRET"),
		Audience:     getEnv(prefix, "AUDIENCE"),
		Scope:        getEnv(prefix, "SCOPE"),
	}, nil

	// return &JwtIssuerEnv{
	// 	Domain:       getEnv(targetjwtprefix, "DOMAIN"),
	// 	ClientID:     getEnv(targetjwtprefix, "CLIENT_ID"),
	// 	ClientSecret: getEnv(targetjwtprefix, "CLIENT_SECRET"),
	// 	Audience:     getEnv(targetjwtprefix, "AUDIENCE"),
	// 	Scope:        getEnv(targetjwtprefix, "SCOPE"),
	// }, nil
}

func (v *JwtIssuerEnv) GetDomain() string       { return v.Domain }
func (v *JwtIssuerEnv) GetClientID() string     { return v.ClientID }
func (v *JwtIssuerEnv) GetClientSecret() string { return v.ClientSecret }
func (v *JwtIssuerEnv) GetAudience() string     { return v.Audience }
func (v *JwtIssuerEnv) GetScope() string        { return v.Scope }
