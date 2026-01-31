package create

import (
	"context"
	"fmt"
	"net/url"

	"github.com/alterejoe/structs"

	"github.com/jackc/pgx/v5/pgxpool"
)

func PGPool(auth *structs.DBAuth) *pgxpool.Pool {
	user := auth.User
	password := auth.Password
	dbName := auth.DBName
	sslmode := auth.SSLMode
	host := auth.Host
	port := auth.Port
	schema := auth.Schema

	fmt.Println(user, password, dbName, sslmode, host, port)
	u := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(user, password),
		Host:     fmt.Sprintf("%s:%s", host, port),
		Path:     dbName,
		RawQuery: fmt.Sprintf("sslmode=%s&search_path=%s", sslmode, schema),
	}

	q := u.Query()
	q.Set("sslmode", sslmode)
	u.RawQuery = q.Encode()

	connString := u.String()

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		panic(err)
	}
	return pool
}
