package api

import (
	"net/http"
	"net/mail"

	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

// @Summary Enviar mensaje de contacto
// @Description Enviar un mensaje desde el formulario de contacto público
// @Tags Contact
// @Accept json
// @Produce json
// @Param request body schemas.ContactRequest true "Información de contacto"
// @Success 200 {object} map[string]string "Mensaje enviado correctamente"
// @Failure 400 {object} errors.Error "Bad Request - Campos faltantes o inválidos"
// @Failure 422 {object} errors.Error "Unprocessable Entity - Cuerpo de solicitud inválido"
// @Failure 500 {object} errors.Error "Internal Server Error - No se pudo enviar el mensaje"
// @Router /contact [post]
func (a *Api) ContactMessage(c echo.Context) error {
	var req schemas.ContactRequest

	if err := c.Bind(&req); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	if req.Name == "" || req.Email == "" || req.Message == "" {
		return errors.HandleError(errors.ContactError.MissingFields, c)
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		return errors.HandleError(errors.ContactError.InvalidEmailFormat, c)
	}

	if err := a.BllController.Contact.SendMessage(&req); err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Mensaje enviado correctamente",
	})
}
