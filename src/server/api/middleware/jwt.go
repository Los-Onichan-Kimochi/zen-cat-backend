package middleware

import (
	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/config"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

// JWTMiddleware validates JWT tokens for protected endpoints
func (a *Middleware) JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Verificar si está en modo desarrollo
		if config.GetDevMode() || a.EnvSettings.DisableAuthForTests {
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

// AdminOnlyMiddleware validates that the user has administrator role
func (a *Middleware) AdminOnlyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Verificar si está en modo desarrollo
		if config.GetDevMode() || a.EnvSettings.DisableAuthForTests {
			// En modo desarrollo, omitir la validación de roles
			a.Logger.Debugln("🔓 Modo desarrollo: Omitiendo validación de rol admin")
			return next(c)
		}

		// Validar JWT token y obtener credenciales
		_, credentials, authError := a.BllController.Auth.AccessTokenValidation(c)
		if authError != nil {
			return errors.HandleError(*authError, c)
		}

		// Verificar que el usuario tenga rol de administrador
		hasAdminRole := false
		for _, role := range credentials.UserRoles {
			if role == string(schemas.UserRolAdmin) {
				hasAdminRole = true
				break
			}
		}

		if !hasAdminRole {
			return errors.HandleError(errors.ForbiddenError.InsufficientPrivileges, c)
		}

		// Si la validación pasa, continuar al siguiente handler
		return next(c)
	}
}

// ClientOnlyMiddleware validates that the user has client role (for user frontend)
func (a *Middleware) ClientOnlyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Verificar si está en modo desarrollo
		if config.GetDevMode() || a.EnvSettings.DisableAuthForTests {
			// En modo desarrollo, omitir la validación de roles
			a.Logger.Debugln("🔓 Modo desarrollo: Omitiendo validación de rol cliente")
			return next(c)
		}

		// Validar JWT token y obtener credenciales
		_, credentials, authError := a.BllController.Auth.AccessTokenValidation(c)
		if authError != nil {
			return errors.HandleError(*authError, c)
		}

		// Verificar que el usuario tenga rol de cliente
		hasClientRole := false
		for _, role := range credentials.UserRoles {
			if role == string(schemas.UserRolClient) {
				hasClientRole = true
				break
			}
		}

		if !hasClientRole {
			return errors.HandleError(errors.ForbiddenError.InsufficientPrivileges, c)
		}

		// Si la validación pasa, continuar al siguiente handler
		return next(c)
	}
}
