package middleware

import (
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/bll/controller"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type AuditMiddleware struct {
	logger        logging.Logger
	bllController *controller.ControllerCollection
}

func NewAuditMiddleware(logger logging.Logger, bllController *controller.ControllerCollection) *AuditMiddleware {
	return &AuditMiddleware{
		logger:        logger,
		bllController: bllController,
	}
}

// AuditMiddleware captures API calls and logs them for audit purposes
func (a *AuditMiddleware) AuditMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Skip audit logging for certain paths
		path := c.Request().URL.Path
		method := c.Request().Method

		// Skip health check, login, swagger, and audit log endpoints
		if shouldSkipAudit(path, method) {
			return next(c)
		}

		// Extract user information from JWT token if available
		var auditContext schemas.AuditContext
		var hasValidUser bool

		// Try to get user credentials from JWT
		if _, credentials, err := a.bllController.Auth.AccessTokenValidation(c); err == nil {
			auditContext = schemas.AuditContext{
				UserId:    credentials.UserId,
				UserEmail: credentials.UserEmail,
				UserRole:  schemas.UserRol(credentials.UserRoles[0]), // Use first role
				IPAddress: getClientIP(c),
				UserAgent: getUserAgent(c),
			}
			hasValidUser = true
		}

		// Execute the request
		err := next(c)

		// Only log if we have a valid user (authenticated requests)
		if hasValidUser {
			// Determine action, entity type, and success status
			action := getActionFromRequest(method, path)
			entityType := getEntityTypeFromPath(path)
			success := err == nil && c.Response().Status < 400

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
			if auditErr := a.bllController.AuditLog.LogAuditEvent(auditContext, event); auditErr != nil {
				a.logger.Error("Failed to log audit event: ", auditErr)
			}
		}

		return err
	}
}

// shouldSkipAudit determines if a request should be skipped from audit logging
func shouldSkipAudit(path, method string) bool {
	skipPaths := []string{
		"/health-check",
		"/swagger",
		"/login",
		"/register",
		"/me",
		"/auth/refresh",
		"/audit-log", // Skip audit log endpoints to avoid infinite loops
	}

	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}

	// Skip GET requests for fetching data (only log CUD operations)
	if method == "GET" {
		return true
	}

	return false
}

// getActionFromRequest determines the audit action based on HTTP method and path
func getActionFromRequest(method, path string) schemas.AuditActionType {
	switch method {
	case "POST":
		if strings.Contains(path, "/bulk-create") || strings.Contains(path, "/bulk") {
			return schemas.AuditActionBulkCreate
		}
		return schemas.AuditActionCreate
	case "PATCH", "PUT":
		return schemas.AuditActionUpdate
	case "DELETE":
		if strings.Contains(path, "/bulk-delete") || strings.Contains(path, "/bulk") {
			return schemas.AuditActionBulkDelete
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
