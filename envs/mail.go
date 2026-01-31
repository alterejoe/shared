package env

import (
	"os"
	"strconv"
)

type MailENV struct {
}

func (auth *MailENV) GetUser() string {
	return os.Getenv("MAIL_USER")
}

func (auth *MailENV) GetPassword() string {
	return os.Getenv("MAIL_PASS")
}

func (auth *MailENV) GetHost() string {
	return os.Getenv("MAIL_HOST")
}

func (auth *MailENV) GetPort() int {
	val, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		panic(err)
	}
	return val
}
