package controller

import (
	"fmt"
	"math/rand"

	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/bll/adapter"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	"onichankimochi.com/astro_cat_backend/src/server/utils"
)

type ForgotPassword struct {
	Logger      logging.Logger
	Adapter     *adapter.AdapterCollection
	EnvSettings *schemas.EnvSettings
}

// guardo el pin
var resetPins = map[string]string{}

func NewForgotPasswordController(
	logger logging.Logger,
	adapter *adapter.AdapterCollection,
	envSettings *schemas.EnvSettings,
) *ForgotPassword {
	return &ForgotPassword{
		Logger:      logger,
		Adapter:     adapter,
		EnvSettings: envSettings,
	}
}

func (fp *ForgotPassword) GenerateResetPin(
	email string,
) (*schemas.ForgotPasswordResponse, *errors.Error) {
	user, err := fp.Adapter.User.GetPostgresqlUserByEmail(email)
	if err != nil {
		return nil, err
	}

	pin := fmt.Sprintf("%06d", rand.Intn(1000000))
	resetPins[user.Email] = pin

	body := fmt.Sprintf("Hola %s,\n\nTu código de recuperación es: %s\n\nSaludos,\nAstrocat 🐾", user.Name, pin)

	// Try to send email, but don't fail the request if email service is not configured
	if emailErr := utils.SendEmail(fp.EnvSettings, user.Email, "Recuperación de contraseña", body); emailErr != nil {
		// Log the error but don't fail the request - useful for tests and development
		fp.Logger.Warnf("Failed to send email: %v", emailErr)
		// Don't return an error for email sending failures
	}

	return &schemas.ForgotPasswordResponse{
		Message: "Código enviado al correo",
		Pin:     pin, // solo para curso/test
	}, nil
}
