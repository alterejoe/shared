package create

import (
	"github.com/alterejoe/structs"

	"gopkg.in/gomail.v2"
)

func CreateMailDialer(auth structs.MailAuth) *gomail.Dialer {

	user := auth.User
	pass := auth.Password
	host := auth.Host
	port := auth.Port

	return gomail.NewDialer(host, port, user, pass)
}
