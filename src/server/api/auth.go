package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
)

// RefreshToken 		godoc
// @Summary 			Refresh user access token.
// @Description 		Refresh user access token.
// @Tags 				Auth
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Success 			200 {object} schemas.TokenResponse "Ok"
// @Router 				/auth/refresh/ [post]
func (a *Api) RefreshToken(c echo.Context) error {
	accessToken, _, authError := a.BllController.Auth.AccessTokenValidation(c)
	if authError != nil {
		return errors.HandleError(*authError, c)
	}

	token, err := a.BllController.Auth.RefreshToken(accessToken)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, token)
}

// Logout 				godoc
// @Summary 			User logout
// @Description 		Logout user session (audit purposes)
// @Tags 				Auth
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Success 			200 {object} string "Logout successful"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Router 				/auth/logout/ [post]
func (a *Api) Logout(c echo.Context) error {
	// Validate token for audit purposes (we want to know who logged out)
	_, _, authError := a.BllController.Auth.AccessTokenValidation(c)
	if authError != nil {
		// Even if token is invalid, we'll return success for security reasons
		// This prevents information disclosure about token validity
		return c.JSON(http.StatusOK, map[string]string{"message": "Logout successful"})
	}

	// The actual token invalidation is handled client-side
	// This endpoint exists primarily for audit logging purposes
	return c.JSON(http.StatusOK, map[string]string{"message": "Logout successful"})
}
