package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

// @Summary 			Create ServiceProfessional.
// @Description 		Associates a service with a professional.
// @Tags 				ServiceProfessional
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request body schemas.CreateServiceProfessionalRequest true "Service-Professional Association Request"
// @Success 			201 {object} schemas.ServiceProfessional "Created"
// @Failure 			400 {object} errors.Error "Bad Request (e.g., invalid updatedBy)"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found (Service or Professional does not exist)"
// @Failure 			409 {object} errors.Error "Conflict (Association already exists)"
// @Failure 			422 {object} errors.Error "Unprocessable Entity (Invalid UUIDs or request body)"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/service-professional/ [post]
func (a *Api) CreateServiceProfessional(c echo.Context) error {
	updatedBy := "ADMIN"

	var request schemas.CreateServiceProfessionalRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	if request.ServiceId == uuid.Nil || request.ProfessionalId == uuid.Nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidServiceProfessionalId, c)
	}

	response, err := a.BllController.ServiceProfessional.CreateServiceProfessional(request, updatedBy)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusCreated, response)
}

// @Summary 			Get ServiceProfessional.
// @Description 		Retrieves a specific service-professional association.
// @Tags 				ServiceProfessional
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param 				serviceId path string true "Service ID"
// @Param 				professionalId path string true "Professional ID"
// @Success 			200 {object} schemas.ServiceProfessional "OK"
// @Failure 			400 {object} errors.Error "Bad Request (e.g., invalid UUID format)"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found (Association does not exist)"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/service-professional/{serviceId}/{professionalId}/ [get]
func (a *Api) GetServiceProfessional(c echo.Context) error {
	serviceId := c.Param("serviceId")
	professionalId := c.Param("professionalId")

	if serviceId == "" || professionalId == "" {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidServiceProfessionalId, c)
	}

	response, err := a.BllController.ServiceProfessional.GetServiceProfessional(serviceId, professionalId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Delete ServiceProfessional.
// @Description 		Deletes a specific service-professional association.
// @Tags 				ServiceProfessional
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param 				serviceId path string true "Service ID"
// @Param 				professionalId path string true "Professional ID"
// @Success 			204 "No Content"
// @Failure 			400 {object} errors.Error "Bad Request (e.g., invalid UUID format)"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found (Association does not exist)"
// @Failure 			500 {object} errors.Error "Internal Server Error (e.g., deletion failed)"
// @Router 				/service-professional/{serviceId}/{professionalId}/ [delete]
func (a *Api) DeleteServiceProfessional(c echo.Context) error {
	serviceId := c.Param("serviceId")
	professionalId := c.Param("professionalId")

	if serviceId == "" || professionalId == "" {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidServiceProfessionalId, c)
	}

	err := a.BllController.ServiceProfessional.DeleteServiceProfessional(serviceId, professionalId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.NoContent(http.StatusNoContent)
}

// @Summary 			Bulk Create ServiceProfessionals.
// @Description 		Creates multiple service-professional associations.
// @Tags 				ServiceProfessional
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request body schemas.BatchCreateServiceProfessionalRequest true "Bulk Create ServiceProfessionals Request"
// @Success 			201 {object} schemas.ServiceProfessionals "Created"
// @Failure 			400 {object} errors.Error "Bad Request (e.g., invalid updatedBy)"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found (Service or Professional does not exist)"
// @Failure 			409 {object} errors.Error "Conflict (Association already exists)"
// @Failure 			422 {object} errors.Error "Unprocessable Entity (Invalid UUIDs or request body)"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/service-professional/bulk/ [post]
func (a *Api) BulkCreateServiceProfessionals(c echo.Context) error {
	updatedBy := "ADMIN"

	var request schemas.BatchCreateServiceProfessionalRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, err := a.BllController.ServiceProfessional.BulkCreateServiceProfessionals(
		request.ServiceProfessionals,
		updatedBy,
	)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusCreated, response)
}

// @Summary 			Fetch ServiceProfessionals.
// @Description 		Fetch all service-professional associations, filtered by params.
// @Tags 				ServiceProfessional
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param 				serviceId query string false "Service ID"
// @Param 				professionalId query string false "Professional ID"
// @Success 			200 {object} schemas.ServiceProfessionals "OK"
// @Failure 			400 {object} errors.Error "Bad Request (e.g., invalid UUID format)"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found (Service or Professional does not exist)"
// @Failure 			422 {object} errors.Error "Unprocessable Entity (Invalid UUIDs or request body)"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/service-professional/ [get]
func (a *Api) FetchServiceProfessionals(c echo.Context) error {
	serviceId := c.QueryParam("serviceId")
	professionalId := c.QueryParam("professionalId")

	response, err := a.BllController.ServiceProfessional.FetchServiceProfessionals(serviceId, professionalId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Bulk Delete ServiceProfessionals.
// @Description 		Bulk deletes service-professional associations.
// @Tags 				ServiceProfessional
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request	body   schemas.BulkDeleteServiceProfessionalRequest true  "Bulk Delete ServiceProfessional Request"
// @Success 			204 "No Content"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/service-professional/bulk/ [delete]
func (a *Api) BulkDeleteServiceProfessionals(c echo.Context) error {
	var request schemas.BulkDeleteServiceProfessionalRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	if err := a.BllController.ServiceProfessional.BulkDeleteServiceProfessionals(request); err != nil {
		return errors.HandleError(*err, c)
	}

	return c.NoContent(http.StatusNoContent)
}

