package api

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

// @Summary 			Get Plan.
// @Description 		Gets a plan given its id.
// @Tags 				Plan
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               planId    path   string  true  "Plan ID"
// @Success 			200 {object} schemas.Plan "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/plan/{planId}/ [get]
func (a *Api) GetPlan(c echo.Context) error {
	planId, parseErr := uuid.Parse(c.Param("planId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidPlanId, c)
	}

	response, err := a.BllController.Plan.GetPlan(planId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Fetch Plans.
// @Description 		Fetch all plans, filtered by params.
// @Tags 				Plan
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               ids query []string false "Plan IDs"
// @Success 			200 {object} schemas.Plans "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/plan/ [get]
func (a *Api) FetchPlans(c echo.Context) error {
	idsString := c.QueryParam("ids")

	ids := []string{}
	if idsString != "" {
		ids = strings.Split(idsString, ",")
	}

	response, err := a.BllController.Plan.FetchPlans(ids)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Create Plan.
// @Description 		Create the plan information.
// @Tags 				Plan
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request body schemas.CreatePlanRequest true "Create Plan Request"
// @Success 			201 {object} schemas.Plan "Created"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/plan/ [post]
func (a *Api) CreatePlan(c echo.Context) error {
	// TODO: Add access token validation (from here we will get the `updatedBy` param)
	updatedBy := "ADMIN" // Placeholder for actual user from token

	var request schemas.CreatePlanRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, newErr := a.BllController.Plan.CreatePlan(request, updatedBy)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	return c.JSON(http.StatusCreated, response)
}

// @Summary 			Update Plan.
// @Description 		Update the plan information.
// @Tags 				Plan
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               planId    path   string  true  "Plan ID"
// @Param               request body schemas.UpdatePlanRequest true "Update Plan Request"
// @Success 			200 {object} schemas.Plan "Ok"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/plan/{planId}/ [patch]
func (a *Api) UpdatePlan(c echo.Context) error {
	// TODO: Add access token validation (from here we will get the `updatedBy` param)
	updatedBy := "ADMIN" // Placeholder for actual user from token

	planId, parseErr := uuid.Parse(c.Param("planId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidPlanId, c)
	}

	var request schemas.UpdatePlanRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, newErr := a.BllController.Plan.UpdatePlan(planId, request, updatedBy)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Delete Plan.
// @Description 		Deletes a plan given its id.
// @Tags 				Plan
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               planId    path   string  true  "Plan ID"
// @Success 			204 "No Content"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/plan/{planId}/ [delete]
func (a *Api) DeletePlan(c echo.Context) error {
	planId, parseErr := uuid.Parse(c.Param("planId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidPlanId, c)
	}

	if err := a.BllController.Plan.DeletePlan(planId); err != nil {
		return errors.HandleError(*err, c)
	}

	return c.NoContent(http.StatusNoContent)
}
