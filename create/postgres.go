package create

import (
	"github.com/alterejoe/shared/interfaces"
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetPGPool(auth interfaces.DBAuth) (*pgxpool.Pool, error) {
	user := auth.GetUser()
	password := auth.GetPassword()
	dbName := auth.GetDBName()
	sslmode := auth.GetSSLMode()
	host := auth.GetHost()
	port := auth.GetPort()
	schema := auth.GetSchema()

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
		// fmt.Println(err)
		// panic(err)
		return nil, err
	}
	return pool, nil
}
