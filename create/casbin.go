package create

import (
	"fmt"

	"github.com/alterejoe/shared/structs"

	"github.com/casbin/casbin/v2"
	pgadapter "github.com/pckhoi/casbin-pgx-adapter/v2"
)

func Casbin(auth *structs.DBAuth) *casbin.Enforcer {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&search_path=%s",
		auth.User,
		auth.Password,
		auth.Host,
		auth.Port,
		auth.DBName,
		auth.SSLMode,
		auth.Schema,
	)

	//
	// pgadapter.WithDatabase(auth.GetDBName()),
	adapter, err := pgadapter.NewAdapter(
		// conf, // <-- THIS must be *pgx.ConnConfig
		dsn,
		pgadapter.WithDatabase(auth.DBName),
		pgadapter.WithSkipTableCreate(),
		pgadapter.WithTableName("casbin_rule"),
	)
	if err != nil {
		panic(err)
	}

	enforcer, err := casbin.NewEnforcer("rbac.conf", adapter)
	if err != nil {
		panic(err)
	}

	return enforcer
}
