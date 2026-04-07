package create

import (
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/driver"
	"github.com/alterejoe/shared/structs"
)

func CHConn(auth *structs.CHAuth) (driver.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{auth.Host + ":" + auth.Port},
		Auth: clickhouse.Auth{
			Database: auth.Database,
			Username: auth.User,
			Password: auth.Password,
		},
		MaxOpenConns:     10,
		MaxIdleConns:     5,
		ConnMaxLifetime:  10 * time.Minute,
		DialTimeout:      10 * time.Second,
		ReadTimeout:      60 * time.Second,
		ConnOpenStrategy: clickhouse.ConnOpenInOrder,
	})
	if err != nil {
		return nil, err
	}
	return conn, nil
}
