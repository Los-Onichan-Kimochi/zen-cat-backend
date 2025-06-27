package utils

import (
	"gopkg.in/gomail.v2"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

func SendEmail(env *schemas.EnvSettings, to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", env.EmailFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(env.EmailHost, env.EmailPort, env.EmailUser, env.EmailPassword)
	return d.DialAndSend(m)
}
