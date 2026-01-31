package create

import (
	"github.com/alterejoe/shared/interfaces"

	"gopkg.in/gomail.v2"
)

func GetMailDialer(auth interfaces.MailAuth) *gomail.Dialer {

	user := auth.GetUser()
	pass := auth.GetPassword()
	host := auth.GetHost()
	port := auth.GetPort()

	return gomail.NewDialer(host, port, user, pass)
}
