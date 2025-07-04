package controller

import (
	"fmt"
	"net/mail"

	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/bll/adapter"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	"onichankimochi.com/astro_cat_backend/src/server/utils"
)

type Contact struct {
	EnvSettings *schemas.EnvSettings
	Adapter     *adapter.AdapterCollection
	Logger      logging.Logger
}

func NewContactController(
	logger logging.Logger,
	adapter *adapter.AdapterCollection,
	envSettings *schemas.EnvSettings,
) *Contact {
	return &Contact{
		Logger:      logger,
		Adapter:     adapter,
		EnvSettings: envSettings,
	}
}

func (c *Contact) SendMessage(req *schemas.ContactRequest) *errors.Error {
	if req.Name == "" || req.Email == "" || req.Subject == "" || req.Message == "" {
		return &errors.ContactError.MissingFields
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		return &errors.ContactError.InvalidEmailFormat
	}
	subject := req.Subject
	body := fmt.Sprintf(
		`Soy %s, con email %s y teléfono %s.

	Mi consulta es: %s`,
		req.Name,
		req.Email,
		req.Phone,
		req.Message,
	)

	if err := utils.SendEmail(c.EnvSettings, c.EnvSettings.EmailFrom, subject, body); err != nil {
		return &errors.ContactError.FailedToSendEmail
	}

	return nil
}
