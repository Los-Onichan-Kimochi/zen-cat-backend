package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

// @Summary 			Create CommunityPlan.
// @Description 		Associates a community with a plan.
// @Tags 				CommunityPlan
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request body schemas.CreateCommunityPlanRequest true "Community-Plan Association Request"
// @Success 			201 {object} schemas.CommunityPlan "Created"
// @Failure 			400 {object} errors.Error "Bad Request (e.g., invalid updatedBy)"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found (Community or Plan does not exist)"
// @Failure 			409 {object} errors.Error "Conflict (Association already exists)"
// @Failure 			422 {object} errors.Error "Unprocessable Entity (Invalid UUIDs or request body)"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/community-plan/ [post]
func (a *Api) CreateCommunityPlan(c echo.Context) error {
	updatedBy := "ADMIN"

	var request schemas.CreateCommunityPlanRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	if request.CommunityId == uuid.Nil || request.PlanId == uuid.Nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidCommunityPlanId, c)
	}

	response, err := a.BllController.CommunityPlan.CreateCommunityPlan(request, updatedBy)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusCreated, response)
}

// @Summary 			Get CommunityPlan.
// @Description 		Retrieves a specific community-plan association.
// @Tags 				CommunityPlan
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param 				communityId path string true "Community ID"
// @Param 				planId path string true "Plan ID"
// @Success 			200 {object} schemas.CommunityPlan "OK"
// @Failure 			400 {object} errors.Error "Bad Request (e.g., invalid UUID format)"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found (Association does not exist)"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/community-plan/{communityId}/{planId}/ [get]
func (a *Api) GetCommunityPlan(c echo.Context) error {
	communityId := c.Param("communityId")
	planId := c.Param("planId")

	if communityId == "" || planId == "" {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidCommunityPlanId, c)
	}

	response, err := a.BllController.CommunityPlan.GetCommunityPlan(communityId, planId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Delete CommunityPlan.
// @Description 		Deletes a specific community-plan association.
// @Tags 				CommunityPlan
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param 				communityId path string true "Community ID"
// @Param 				planId path string true "Plan ID"
// @Success 			204 "No Content"
// @Failure 			400 {object} errors.Error "Bad Request (e.g., invalid UUID format)"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found (Association does not exist)"
// @Failure 			500 {object} errors.Error "Internal Server Error (e.g., deletion failed)"
// @Router 				/community-plan/{communityId}/{planId}/ [delete]
func (a *Api) DeleteCommunityPlan(c echo.Context) error {
	communityId := c.Param("communityId")
	planId := c.Param("planId")

	if communityId == "" || planId == "" {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidCommunityPlanId, c)
	}

	err := a.BllController.CommunityPlan.DeleteCommunityPlan(communityId, planId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.NoContent(http.StatusNoContent)
}

// @Summary 			Bulk Create CommunityPlans.
// @Description 		Creates multiple community-plan associations.
// @Tags 				CommunityPlan
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request body schemas.BatchCreateCommunityPlanRequest true "Bulk Create CommunityPlans Request"
// @Success 			201 {object} schemas.CommunityPlans "Created"
// @Failure 			400 {object} errors.Error "Bad Request (e.g., invalid updatedBy)"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found (Community or Plan does not exist)"
// @Failure 			409 {object} errors.Error "Conflict (Association already exists)"
// @Failure 			422 {object} errors.Error "Unprocessable Entity (Invalid UUIDs or request body)"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/community-plan/bulk/ [post]
func (a *Api) BulkCreateCommunityPlans(c echo.Context) error {
	updatedBy := "ADMIN"

	var request schemas.BatchCreateCommunityPlanRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, err := a.BllController.CommunityPlan.BulkCreateCommunityPlans(
		request.CommunityPlans,
		updatedBy,
	)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusCreated, response)
}

// @Summary 			Fetch CommunityPlans.
// @Description 		Fetch all community-plan associations, filtered by params.
// @Tags 				CommunityPlan
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param 				communityId query string false "Community ID"
// @Param 				planId query string false "Plan ID"
// @Success 			200 {object} schemas.CommunityPlans "OK"
// @Failure 			400 {object} errors.Error "Bad Request (e.g., invalid UUID format)"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found (Community or Plan does not exist)"
// @Failure 			422 {object} errors.Error "Unprocessable Entity (Invalid UUIDs or request body)"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/community-plan/ [get]
func (a *Api) FetchCommunityPlans(c echo.Context) error {
	communityId := c.QueryParam("communityId")
	planId := c.QueryParam("planId")

	response, err := a.BllController.CommunityPlan.FetchCommunityPlans(communityId, planId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Bulk Delete CommunityPlans.
// @Description 		Bulk deletes community-plan associations.
// @Tags 				CommunityPlan
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request	body   schemas.BulkDeleteCommunityPlanRequest true  "Bulk Delete CommunityPlan Request"
// @Success 			204 "No Content"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/community-plan/bulk/ [delete]
func (a *Api) BulkDeleteCommunityPlans(c echo.Context) error {
	var request schemas.BulkDeleteCommunityPlanRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	if err := a.BllController.CommunityPlan.BulkDeleteCommunityPlans(request); err != nil {
		return errors.HandleError(*err, c)
	}

	return c.NoContent(http.StatusNoContent)
}
