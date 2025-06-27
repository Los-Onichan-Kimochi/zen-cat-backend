package api

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

// @Summary Recovery Password
// @Description Envía un código de recuperación (PIN) al correo electrónico del usuario
// @Tags ForgotPassword
// @Accept json
// @Produce json
// @Param request body schemas.ForgotPasswordRequest true "Email del usuario"
// @Success 200 {object} schemas.ForgotPasswordResponse "Código enviado exitosamente"
// @Failure 400 {object} errors.Error "Bad Request - Error al enviar el código"
// @Failure 404 {object} errors.Error "Not Found - Usuario no encontrado"
// @Failure 422 {object} errors.Error "Unprocessable Entity - Formato incorrecto"
// @Failure 500 {object} errors.Error "Internal Server Error"
// @Router /forgot-password/ [post]
func (a *Api) ForgotPassword(c echo.Context) error {
	var request schemas.ForgotPasswordRequest

	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	// Validate email format
	if request.Email == "" {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidUserEmail, c)
	}

	// Basic email format validation
	if !strings.Contains(request.Email, "@") || !strings.Contains(request.Email, ".") {
		return errors.HandleError(errors.BadRequestError.InvalidUpdatedByValue, c) // Using a generic bad request error for invalid email format
	}

	response, err := a.BllController.ForgotPassword.GenerateResetPin(request.Email)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}
