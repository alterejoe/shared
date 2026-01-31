package create

import (
	"clerk-portal/shared/interfaces"
	"fmt"

	"github.com/casbin/casbin/v2"
	pgadapter "github.com/pckhoi/casbin-pgx-adapter/v2"
)

func GetCasbin(auth interfaces.DBAuth) (*casbin.Enforcer, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&search_path=%s",
		auth.GetUser(),
		auth.GetPassword(),
		auth.GetHost(),
		auth.GetPort(),
		auth.GetDBName(),
		auth.GetSSLMode(),
		auth.GetSchema(),
	)

	//
	// pgadapter.WithDatabase(auth.GetDBName()),
	adapter, err := pgadapter.NewAdapter(
		// conf, // <-- THIS must be *pgx.ConnConfig
		dsn,
		pgadapter.WithDatabase(auth.GetDBName()),
		pgadapter.WithSkipTableCreate(),
		pgadapter.WithTableName("casbin_rule"),
	)
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer("rbac.conf", adapter)
	if err != nil {
		return nil, err
	}

	return enforcer, nil
}
