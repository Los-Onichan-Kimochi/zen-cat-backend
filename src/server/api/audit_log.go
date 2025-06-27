package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
)

// GetAuditLogs 			godoc
// @Summary 			Get Audit Logs
// @Description 		Retrieve audit logs with optional filtering and pagination
// @Tags 				AuditLog
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param 				userIds query string false "Comma-separated user IDs"
// @Param 				actions query string false "Comma-separated action types"
// @Param 				entityTypes query string false "Comma-separated entity types"
// @Param 				userRoles query string false "Comma-separated user roles"
// @Param 				startDate query string false "Start date (YYYY-MM-DD)"
// @Param 				endDate query string false "End date (YYYY-MM-DD)"
// @Param 				success query string false "Filter by success (true/false)"
// @Param 				page query string false "Page number (default: 1)"
// @Param 				pageSize query string false "Page size (default: 50, max: 200)"
// @Success 			200 {object} schemas.AuditLogs "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/audit-log/ [get]
func (a *Api) GetAuditLogs(c echo.Context) error {
	userIds := c.QueryParam("userIds")
	actions := c.QueryParam("actions")
	entityTypes := c.QueryParam("entityTypes")
	userRoles := c.QueryParam("userRoles")
	startDate := c.QueryParam("startDate")
	endDate := c.QueryParam("endDate")
	success := c.QueryParam("success")
	page := c.QueryParam("page")
	pageSize := c.QueryParam("pageSize")

	response, err := a.BllController.AuditLog.GetAuditLogs(
		userIds,
		actions,
		entityTypes,
		userRoles,
		startDate,
		endDate,
		success,
		page,
		pageSize,
	)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// GetAuditLogById 		godoc
// @Summary 			Get Audit Log by ID
// @Description 		Retrieve a specific audit log by its ID
// @Tags 				AuditLog
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param 				auditLogId path string true "Audit Log ID"
// @Success 			200 {object} schemas.AuditLog "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/audit-log/{auditLogId}/ [get]
func (a *Api) GetAuditLogById(c echo.Context) error {
	auditLogId, parseErr := uuid.Parse(c.Param("auditLogId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, err := a.BllController.AuditLog.GetAuditLogById(auditLogId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// GetAuditStats 			godoc
// @Summary 			Get Audit Statistics
// @Description 		Retrieve audit statistics for the specified time period
// @Tags 				AuditLog
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param 				days query string false "Number of days to include in stats (default: 30)"
// @Success 			200 {object} schemas.AuditStats "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/audit-log/stats/ [get]
func (a *Api) GetAuditStats(c echo.Context) error {
	days := c.QueryParam("days")

	response, err := a.BllController.AuditLog.GetAuditStats(days)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// GetErrorLogs 			godoc
// @Summary 			Get Error Logs
// @Description 		Retrieve error logs with optional filtering and pagination
// @Tags 				ErrorLog
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param 				userIds query string false "Comma-separated user IDs"
// @Param 				actions query string false "Comma-separated action types"
// @Param 				entityTypes query string false "Comma-separated entity types"
// @Param 				userRoles query string false "Comma-separated user roles"
// @Param 				startDate query string false "Start date (YYYY-MM-DD)"
// @Param 				endDate query string false "End date (YYYY-MM-DD)"
// @Param 				page query string false "Page number (default: 1)"
// @Param 				pageSize query string false "Page size (default: 50, max: 200)"
// @Success 			200 {object} schemas.AuditLogs "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/error-log/ [get]
func (a *Api) GetErrorLogs(c echo.Context) error {
	userIds := c.QueryParam("userIds")
	actions := c.QueryParam("actions")
	entityTypes := c.QueryParam("entityTypes")
	userRoles := c.QueryParam("userRoles")
	startDate := c.QueryParam("startDate")
	endDate := c.QueryParam("endDate")
	page := c.QueryParam("page")
	pageSize := c.QueryParam("pageSize")

	// Force success=false for error logs - reuse GetAuditLogs logic
	response, err := a.BllController.AuditLog.GetAuditLogs(
		userIds,
		actions,
		entityTypes,
		userRoles,
		startDate,
		endDate,
		"false", // Always pass success=false for error logs
		page,
		pageSize,
	)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// GetErrorStats 			godoc
// @Summary 			Get Error Statistics
// @Description 		Retrieve error statistics for the specified time period
// @Tags 				ErrorLog
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param 				days query string false "Number of days to include in stats (default: 30)"
// @Success 			200 {object} schemas.AuditStats "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/error-log/stats/ [get]
func (a *Api) GetErrorStats(c echo.Context) error {
	days := c.QueryParam("days")

	response, err := a.BllController.AuditLog.GetErrorStats(days)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// GetErrorLogById 		godoc
// @Summary 			Get Error Log by ID
// @Description 		Retrieve a specific error log by its ID
// @Tags 				ErrorLog
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param 				auditLogId path string true "Audit Log ID"
// @Success 			200 {object} schemas.AuditLog "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/error-log/{auditLogId}/ [get]
func (a *Api) GetErrorLogById(c echo.Context) error {
	auditLogId, parseErr := uuid.Parse(c.Param("auditLogId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, err := a.BllController.AuditLog.GetAuditLogById(auditLogId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	// Ensure this is actually an error log (success = false)
	if response.Success {
		return errors.HandleError(errors.ObjectNotFoundError.AuditLogNotFound, c)
	}

	return c.JSON(http.StatusOK, response)
}

// DeleteOldAuditLogs 	godoc
// @Summary 			Delete Old Audit Logs
// @Description 		Delete audit logs older than the specified number of days
// @Tags 				AuditLog
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param 				days query string false "Number of days to keep (default: 90)"
// @Success 			204 "No Content"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/audit-log/cleanup/ [delete]
func (a *Api) DeleteOldAuditLogs(c echo.Context) error {
	days := c.QueryParam("days")

	if err := a.BllController.AuditLog.DeleteOldAuditLogs(days); err != nil {
		return errors.HandleError(*err, c)
	}

	return c.NoContent(http.StatusNoContent)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func uuidPtr(u uuid.UUID) *uuid.UUID {
	return &u
}
