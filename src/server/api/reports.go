package api

import (
	"time"

	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
)

// @Summary      Get Service Report
// @Description  Returns aggregated data for service reservations for admin dashboards.
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Security     JWT
// @Param        from      query   string  false  "Start date (YYYY-MM-DD)"
// @Param        to        query   string  false  "End date (YYYY-MM-DD)"
// @Param        serviceType query string  false  "Service type (e.g., yoga, citas médicas)"
// @Param        groupBy   query   string  false  "Group by (day, week, month)"
// @Success      200 {object} map[string]interface{} "Service report data"
// @Failure      400 {object} errors.Error "Bad Request"
// @Failure      401 {object} errors.Error "Missing or malformed JWT"
// @Failure      403 {object} errors.Error "Forbidden - Admin role required"
// @Router       /reports/services [get]
func (a *Api) GetServiceReport(c echo.Context) error {
	fromStr := c.QueryParam("from")
	toStr := c.QueryParam("to")
	groupBy := c.QueryParam("groupBy")
	if groupBy == "" {
		groupBy = "day"
	}

	var from, to *time.Time
	if fromStr != "" {
		parsed, err := time.Parse("2006-01-02", fromStr)
		if err != nil {
			return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
		}
		from = &parsed
	}
	if toStr != "" {
		parsed, err := time.Parse("2006-01-02", toStr)
		if err != nil {
			return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
		}
		to = &parsed
	}

	report, err := a.BllController.Reservation.GetServiceReport(from, to, groupBy)
	if err != nil {
		return errors.HandleError(*err, c)
	}
	return c.JSON(200, report)
}
