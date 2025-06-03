package api

import (
	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
)

// JWTMiddleware validates JWT tokens for protected endpoints
func (a *Api) JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Validate JWT token using existing auth logic
		_, _, authError := a.BllController.Auth.AccessTokenValidation(c)
		if authError != nil {
			return errors.HandleError(*authError, c)
		}

		// If validation passes, continue to the next handler
		return next(c)
	}
}
