package env

import "os"

type DBENV struct {
}

func (auth *DBENV) GetUser() string {
	return os.Getenv("DB_USER")
}

func (auth *DBENV) GetPassword() string {
	return os.Getenv("DB_PASS")
}

func (auth *DBENV) GetDBName() string {
	return os.Getenv("DB_NAME")
}

func (auth *DBENV) GetHost() string {
	return os.Getenv("DB_HOST")
}

func (auth *DBENV) GetPort() string {
	return os.Getenv("DB_PORT")
}

func (auth *DBENV) GetSSLMode() string {
	return os.Getenv("DB_SSLMODE")
}

func (auth *DBENV) GetSchema() string {
	return os.Getenv("DB_SCHEMA")
}
