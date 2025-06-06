package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

// @Summary 			Create CommunityService.
// @Description 		Associates a community with a service.
// @Tags 				CommunityService
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request body schemas.CreateCommunityServiceRequest true "Community-Service Association Request"
// @Success 			201 {object} schemas.CommunityService "Created"
// @Failure 			400 {object} errors.Error "Bad Request (e.g., invalid updatedBy)"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found (Community or Service does not exist)"
// @Failure 			409 {object} errors.Error "Conflict (Association already exists)"
// @Failure 			422 {object} errors.Error "Unprocessable Entity (Invalid UUIDs or request body)"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/community-service/ [post]
func (a *Api) CreateCommunityService(c echo.Context) error {
	updatedBy := "ADMIN"

	var request schemas.CreateCommunityServiceRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	if request.CommunityId == uuid.Nil || request.ServiceId == uuid.Nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidCommunityServiceId, c)
	}

	response, err := a.BllController.CommunityService.CreateCommunityService(request, updatedBy)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusCreated, response)
}

// @Summary 			Get CommunityService.
// @Description 		Retrieves a specific community-service association.
// @Tags 				CommunityService
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param 				communityId path string true "Community ID"
// @Param 				serviceId path string true "Service ID"
// @Success 			200 {object} schemas.CommunityService "OK"
// @Failure 			400 {object} errors.Error "Bad Request (e.g., invalid UUID format)"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found (Association does not exist)"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/community-service/{communityId}/{serviceId}/ [get]
func (a *Api) GetCommunityService(c echo.Context) error {
	communityId := c.Param("communityId")
	serviceId := c.Param("serviceId")

	if communityId == "" || serviceId == "" {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidCommunityServiceId, c)
	}

	response, err := a.BllController.CommunityService.GetCommunityService(communityId, serviceId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Delete CommunityService.
// @Description 		Deletes a specific community-service association.
// @Tags 				CommunityService
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param 				communityId path string true "Community ID"
// @Param 				serviceId path string true "Service ID"
// @Success 			204 "No Content"
// @Failure 			400 {object} errors.Error "Bad Request (e.g., invalid UUID format)"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found (Association does not exist)"
// @Failure 			500 {object} errors.Error "Internal Server Error (e.g., deletion failed)"
// @Router 				/community-service/{communityId}/{serviceId}/ [delete]
func (a *Api) DeleteCommunityService(c echo.Context) error {
	communityId := c.Param("communityId")
	serviceId := c.Param("serviceId")

	if communityId == "" || serviceId == "" {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidCommunityServiceId, c)
	}

	err := a.BllController.CommunityService.DeleteCommunityService(communityId, serviceId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.NoContent(http.StatusNoContent)
}

// @Summary 			Bulk Create CommunityServices.
// @Description 		Creates multiple community-service associations.
// @Tags 				CommunityService
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request body schemas.BatchCreateCommunityServiceRequest true "Bulk Create CommunityServices Request"
// @Success 			201 {object} schemas.CommunityServices "Created"
// @Failure 			400 {object} errors.Error "Bad Request (e.g., invalid updatedBy)"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found (Community or Service does not exist)"
// @Failure 			409 {object} errors.Error "Conflict (Association already exists)"
// @Failure 			422 {object} errors.Error "Unprocessable Entity (Invalid UUIDs or request body)"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/community-service/bulk-create/ [post]
func (a *Api) BulkCreateCommunityServices(c echo.Context) error {
	updatedBy := "ADMIN"

	var request schemas.BatchCreateCommunityServiceRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, err := a.BllController.CommunityService.BulkCreateCommunityServices(
		request.CommunityServices,
		updatedBy,
	)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusCreated, response)
}

// @Summary 			Bulk Delete CommunityServices.
// @Description 		Bulk deletes community-service associations.
// @Tags 				CommunityService
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request	body   schemas.BulkDeleteCommunityServiceRequest true  "Bulk Delete CommunityService Request"
// @Success 			204 "No Content"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/community-service/bulk-delete/ [delete]
func (a *Api) BulkDeleteCommunityServices(c echo.Context) error {
	var request schemas.BulkDeleteCommunityServiceRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	if err := a.BllController.CommunityService.BulkDeleteCommunityServices(request); err != nil {
		return errors.HandleError(*err, c)
	}

	return c.NoContent(http.StatusNoContent)
}

// @Summary 			Fetch CommunityServices.
// @Description 		Fetch all community-service associations, filtered by params.
// @Tags 				CommunityService
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param 				communityId query string false "Community ID"
// @Param 				serviceId query string false "Service ID"
// @Success 			200 {object} schemas.CommunityServices "OK"
// @Failure 			400 {object} errors.Error "Bad Request (e.g., invalid UUID format)"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found (Community or Service does not exist)"
// @Failure 			422 {object} errors.Error "Unprocessable Entity (Invalid UUIDs or request body)"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/community-service/ [get]
func (a *Api) FetchCommunityServices(c echo.Context) error {
	communityId := c.QueryParam("communityId")
	serviceId := c.QueryParam("serviceId")

	response, err := a.BllController.CommunityService.FetchCommunityServices(communityId, serviceId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}
