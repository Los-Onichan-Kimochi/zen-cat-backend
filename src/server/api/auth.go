package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
)

// HealthCheck 			godoc
// @Summary 			Refresh player access token.
// @Description 		Refresh player access token.
// @Tags 				Player
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Success 			200 {object} schemas.TokenResponse "Ok"
// @Router 				/player/refresh/ [post]
func (a *Api) RefreshToken(c echo.Context) error {
	accessToken, _, authError := a.BllController.Auth.AccessTokenValidation(c)
	if authError != nil {
		return errors.HandleError(errors.AuthenticationError.UnauthorizedUser, c)
	}

	token, err := a.BllController.Auth.RefreshToken(accessToken)
	if err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidParsingInteger, c)
	}

	return c.JSON(http.StatusOK, token)
}
