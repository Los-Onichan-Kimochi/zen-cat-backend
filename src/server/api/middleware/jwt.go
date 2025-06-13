package middleware

import (
	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/config"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
)

// JWTMiddleware validates JWT tokens for protected endpoints
func (a *Middleware) JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Verificar si está en modo desarrollo
		if config.GetDevMode() {
			// En modo desarrollo, omitir la validación JWT
			a.Logger.Debugln("🔓 Modo desarrollo: Omitiendo validación JWT")
			return next(c)
		}

		// En modo producción, validar JWT token usando la lógica existente
		_, _, authError := a.BllController.Auth.AccessTokenValidation(c)
		if authError != nil {
			return errors.HandleError(*authError, c)
		}

		// Si la validación pasa, continuar al siguiente handler
		return next(c)
	}
}
