package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/bll/controller"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type Middleware struct {
	Logger        logging.Logger
	BllController *controller.ControllerCollection
	EnvSettings   *schemas.EnvSettings
	Echo          *echo.Echo
}

func NewMiddleware(
	logger logging.Logger,
	bllController *controller.ControllerCollection,
	envSettings *schemas.EnvSettings,
	echo *echo.Echo,
) *Middleware {
	return &Middleware{
		Logger:        logger,
		BllController: bllController,
		EnvSettings:   envSettings,
		Echo:          echo,
	}
}

// AuditMiddleware captures API calls and logs them for audit purposes
func (a *Middleware) AuditMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if a.EnvSettings.DisableAuthForTests {
			return next(c)
		}

		// Obtengo path y method
		path := c.Request().URL.Path
		method := c.Request().Method

		shouldLog := true
		success := true

		// Extract user information from JWT token if available
		var auditContext schemas.AuditContext
		var hasValidUser bool

		// For login/register requests, extract email from request body before processing
		var loginEmail string
		if method == "POST" && (strings.Contains(path, "/login") || strings.Contains(path, "/register")) {
			loginEmail = extractEmailFromLoginRequest(c, path)
		}

		// Try to get user credentials from JWT
		if _, credentials, err := a.BllController.Auth.AccessTokenValidation(c); err == nil {
			auditContext = schemas.AuditContext{
				UserId:    credentials.UserId,
				UserEmail: credentials.UserEmail,
				UserRole:  schemas.UserRol(credentials.UserRoles[0]), // Use first role
				IPAddress: getClientIP(c),
				UserAgent: getUserAgent(c),
			}
			hasValidUser = true
		}
		// Ejecuto la request
		err := next(c)
		code := c.Response().Status
		// Skip health check, login, swagger, and audit log endpoints - but only for successful requests
		shouldSkipForPath := shouldSkipAuditForPath(path, method, code, hasValidUser)

		// Determine success status
		if code >= 300 {
			success = false
		}
		// Log audit event if:
		// 1. We have a valid user (authenticated requests), OR
		// 2. The request failed (error occurred or status code != 2xx)
		// shouldLog := hasValidUser || !success

		// Skip logging if this is a successful request on a path that should be skipped
		if shouldSkipForPath {
			shouldLog = false
		}

		// Additional rule: Skip successful GET requests for unauthenticated users
		// This prevents spam from frontend trying to load protected resources
		if success && method == "GET" && !hasValidUser {
			shouldLog = false
		}

		if shouldLog {
			// For successful login/register, try to get the actual user role
			if !hasValidUser && success && loginEmail != "" && (strings.Contains(path, "/login") || strings.Contains(path, "/register")) {
				// For successful login/register, try to get the user's actual role
				var userRole schemas.UserRol = schemas.UserRolGuest // Default fallback

				// Try to get user by email to get their actual role
				if user, err := a.BllController.User.Adapter.User.GetPostgresqlUserByEmail(loginEmail); err == nil && user != nil {
					userRole = user.Rol
				}

				auditContext = schemas.AuditContext{
					UserId:    uuid.Nil, // We don't have the user ID yet, but we have email
					UserEmail: loginEmail,
					UserRole:  userRole, // Use actual role instead of always GUEST
					IPAddress: getClientIP(c),
					UserAgent: getUserAgent(c),
				}
				hasValidUser = true
			}

			// Create audit context for unauthenticated error cases
			if !hasValidUser {
				// For unauthenticated requests, create a minimal context
				auditContext = schemas.AuditContext{
					UserId:    uuid.Nil, // Use nil UUID for unauthenticated requests
					UserEmail: "sin autenticar",
					UserRole:  schemas.UserRolGuest, // Use guest role for unauthenticated requests
					IPAddress: getClientIP(c),
					UserAgent: getUserAgent(c),
				}
			}

			// Determine action, entity type, and success status
			action := getActionFromRequest(method, path)
			entityType := getEntityTypeFromPath(path)

			// Create audit event
			event := schemas.AuditEvent{
				Action:         action,
				EntityType:     entityType,
				EntityId:       getEntityIdFromPath(path),
				EntityName:     getEntityNameFromPath(path, entityType),
				OldValues:      nil, // We don't capture old values in middleware for simplicity
				NewValues:      nil, // We don't capture new values in middleware for simplicity
				AdditionalInfo: getAdditionalInfo(method, path, c.QueryParams()),
				Success:        success,
				ErrorMessage:   getErrorMessage(err),
			}

			// Log the audit event (don't fail the request if audit logging fails)
			if auditErr := a.BllController.AuditLog.LogAuditEvent(auditContext, event); auditErr != nil {
				a.Logger.Error("Failed to log audit event: ", auditErr)
			}
		}

		return err
	}
}

// extractEmailFromLoginRequest extracts email from login/register request body
func extractEmailFromLoginRequest(c echo.Context, path string) string {
	// Read the request body
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return ""
	}

	// Restore the request body for the actual handler
	c.Request().Body = io.NopCloser(bytes.NewReader(body))

	// Try to parse as login request
	if strings.Contains(path, "/login") {
		var loginReq schemas.LoginRequest
		if err := json.Unmarshal(body, &loginReq); err == nil {
			return loginReq.Email
		}
	}

	// Try to parse as register request
	if strings.Contains(path, "/register") {
		var registerReq schemas.RegisterRequest
		if err := json.Unmarshal(body, &registerReq); err == nil {
			return registerReq.Email
		}
	}

	return ""
}

// shouldSkipAuditForPath determines if a request should be skipped from audit logging
// This is more intelligent than the previous version
func shouldSkipAuditForPath(path, method string, code int, hasValidUser bool) bool {
	// Always log errors (4xx, 5xx) for security monitoring
	if code >= 400 {
		// Exception: Skip 401 errors for refresh token attempts - these are expected
		if code == 401 && strings.HasPrefix(path, "/auth/refresh") {
			return true
		}
		// Exception: Skip 401 errors for audit/error log pages when not authenticated
		// This prevents spam when frontend tries to load these pages before auth
		if code == 401 && (strings.HasPrefix(path, "/audit-log") || strings.HasPrefix(path, "/error-log")) {
			return true
		}
		return false // Log other errors
	}

	// Always log logout events regardless of authentication status
	if strings.Contains(path, "/logout") {
		return false // Never skip logout
	}

	// Paths that should always be skipped for successful requests
	alwaysSkipPaths := []string{
		"/health-check",
		"/swagger",
	}

	for _, skipPath := range alwaysSkipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}

	// Skip successful requests to monitoring endpoints when authenticated
	// (we want to track access, but not spam the logs)
	if code < 300 && hasValidUser {
		monitoringPaths := []string{
			"/audit-log",
			"/error-log",
			"/me",
		}

		for _, monitoringPath := range monitoringPaths {
			if strings.HasPrefix(path, monitoringPath) && method == "GET" {
				return true
			}
		}
	}

	// Skip successful GET requests in general (only log CUD operations)
	// This significantly reduces log volume while maintaining security auditing
	if method == "GET" && code < 300 {
		return true
	}

	return false
}

// getActionFromRequest determines the audit action based on HTTP method and path
func getActionFromRequest(method, path string) schemas.AuditActionType {
	switch method {
	case "POST":
		// Check for authentication-specific paths first
		if strings.HasPrefix(path, "/login") || strings.Contains(path, "/login") {
			return schemas.AuditActionLogin
		}
		if strings.HasPrefix(path, "/logout") || strings.Contains(path, "/logout") {
			return schemas.AuditActionLogout
		}
		if strings.HasPrefix(path, "/register") || strings.Contains(path, "/register") {
			return schemas.AuditActionRegister
		}
		if strings.HasPrefix(path, "/auth/refresh") || strings.Contains(path, "/auth/refresh") {
			return schemas.AuditActionLogin // Refresh is essentially a re-login
		}
		if strings.Contains(path, "/bulk-create") || strings.Contains(path, "/bulk") {
			return schemas.AuditActionBulkCreate
		}
		// Check for reservation-specific actions
		if strings.Contains(path, "/reservation") {
			if strings.Contains(path, "/cancel") {
				return schemas.AuditActionCancelReservation
			}
			return schemas.AuditActionCreateReservation
		}
		// Check for subscription actions
		if strings.Contains(path, "/subscribe") {
			return schemas.AuditActionSubscribe
		}
		if strings.Contains(path, "/unsubscribe") {
			return schemas.AuditActionUnsubscribe
		}
		// Default POST action
		return schemas.AuditActionCreate
	case "PATCH", "PUT":
		if strings.Contains(path, "/profile") || strings.Contains(path, "/user") {
			return schemas.AuditActionUpdateProfile
		}
		return schemas.AuditActionUpdate
	case "DELETE":
		if strings.Contains(path, "/bulk-delete") || strings.Contains(path, "/bulk") {
			return schemas.AuditActionBulkDelete
		}
		if strings.Contains(path, "/reservation") && strings.Contains(path, "/cancel") {
			return schemas.AuditActionCancelReservation
		}
		return schemas.AuditActionDelete
	default:
		return schemas.AuditActionCreate // Default fallback
	}
}

// getEntityTypeFromPath determines the entity type from the request path
func getEntityTypeFromPath(path string) schemas.AuditEntityType {
	if strings.Contains(path, "/community-plan") {
		return schemas.AuditEntityCommunityPlan
	} else if strings.Contains(path, "/community-service") {
		return schemas.AuditEntityCommunityService
	} else if strings.Contains(path, "/service-local") {
		return schemas.AuditEntityServiceLocal
	} else if strings.Contains(path, "/service-professional") {
		return schemas.AuditEntityServiceProfessional
	} else if strings.Contains(path, "/community") {
		return schemas.AuditEntityCommunity
	} else if strings.Contains(path, "/professional") {
		return schemas.AuditEntityProfessional
	} else if strings.Contains(path, "/local") {
		return schemas.AuditEntityLocal
	} else if strings.Contains(path, "/service") {
		return schemas.AuditEntityService
	} else if strings.Contains(path, "/user") {
		return schemas.AuditEntityUser
	} else if strings.Contains(path, "/plan") {
		return schemas.AuditEntityPlan
	} else if strings.Contains(path, "/session") {
		return schemas.AuditEntitySession
	} else if strings.Contains(path, "/reservation") {
		return schemas.AuditEntityReservation
	} else if strings.Contains(path, "/onboarding") {
		return schemas.AuditEntityOnboarding
	}

	return schemas.AuditEntityUser // Default fallback
}

// getEntityIdFromPath extracts entity ID from path if possible
func getEntityIdFromPath(path string) *uuid.UUID {
	// This is a simplified extraction - you might want to make it more sophisticated
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if id, err := uuid.Parse(part); err == nil {
			return &id
		}
	}
	return nil
}

// getEntityNameFromPath gets a descriptive name for the entity
func getEntityNameFromPath(path string, entityType schemas.AuditEntityType) *string {
	if strings.Contains(path, "/bulk") {
		name := "Bulk " + string(entityType) + " operation"
		return &name
	}

	name := string(entityType) + " operation"
	return &name
}

// getAdditionalInfo provides additional context about the request
func getAdditionalInfo(method, path string, queryParams map[string][]string) *string {
	info := "Method: " + method + ", Path: " + path
	if len(queryParams) > 0 {
		info += ", Query params: " + strings.Join(getQueryParamKeys(queryParams), ",")
	}
	return &info
}

// getQueryParamKeys extracts query parameter keys
func getQueryParamKeys(params map[string][]string) []string {
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	return keys
}

// getClientIP extracts the client IP address
func getClientIP(c echo.Context) string {
	// Check various headers for the real IP
	ip := c.Request().Header.Get("X-Forwarded-For")
	if ip != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		if idx := strings.Index(ip, ","); idx != -1 {
			ip = ip[:idx]
		}
		return strings.TrimSpace(ip)
	}

	ip = c.Request().Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	return c.RealIP()
}

// getUserAgent extracts the user agent
func getUserAgent(c echo.Context) *string {
	userAgent := c.Request().Header.Get("User-Agent")
	if userAgent == "" {
		return nil
	}
	return &userAgent
}

// getErrorMessage extracts error message if request failed
func getErrorMessage(err error) *string {
	if err != nil {
		errMsg := err.Error()
		return &errMsg
	}
	return nil
}
