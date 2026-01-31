package structs

import "fmt"

type DBAuth struct {
	User     string
	Password string
	DBName   string
	Host     string
	Port     string
	SSLMode  string
	Schema   string
}

func (a DBAuth) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		a.User,
		a.Password,
		a.Host,
		a.Port,
		a.DBName,
		a.SSLMode,
	)
}

type RedisAuth struct {
	Host string
	Port string
}
type MailAuth struct {
	User     string
	Password string
	Host     string
	Port     int
}
type S3Auth struct {
	Bucket    string
	Endpoint  string
	Region    string
	AccessKey string
	SecretKey string
}
