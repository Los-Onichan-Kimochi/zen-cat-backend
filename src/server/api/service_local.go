package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

// @Summary 			Create ServiceLocal.
// @Description 		Associates a service with a local.
// @Tags 				ServiceLocal
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request body schemas.CreateServiceLocalRequest true "Service-Local Association Request"
// @Success 			201 {object} schemas.ServiceLocal "Created"
// @Failure 			400 {object} errors.Error "Bad Request (e.g., invalid updatedBy)"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found (Service or Local does not exist)"
// @Failure 			409 {object} errors.Error "Conflict (Association already exists)"
// @Failure 			422 {object} errors.Error "Unprocessable Entity (Invalid UUIDs or request body)"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/service-local/ [post]
func (a *Api) CreateServiceLocal(c echo.Context) error {
	updatedBy := "ADMIN"

	var request schemas.CreateServiceLocalRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	if request.ServiceId == uuid.Nil || request.LocalId == uuid.Nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidServiceLocalId, c)
	}

	response, err := a.BllController.ServiceLocal.CreateServiceLocal(request, updatedBy)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusCreated, response)
}

// @Summary 			Get ServiceLocal.
// @Description 		Retrieves a specific service-local association.
// @Tags 				ServiceLocal
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param 				serviceId path string true "Service ID"
// @Param 				localId path string true "Local ID"
// @Success 			200 {object} schemas.ServiceLocal "OK"
// @Failure 			400 {object} errors.Error "Bad Request (e.g., invalid UUID format)"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found (Association does not exist)"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/service-local/{serviceId}/{localId}/ [get]
func (a *Api) GetServiceLocal(c echo.Context) error {
	serviceId := c.Param("serviceId")
	localId := c.Param("localId")

	if serviceId == "" || localId == "" {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidServiceLocalId, c)
	}

	response, err := a.BllController.ServiceLocal.GetServiceLocal(serviceId, localId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Delete ServiceLocal.
// @Description 		Deletes a specific service-local association.
// @Tags 				ServiceLocal
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param 				serviceId path string true "Service ID"
// @Param 				localId path string true "Local ID"
// @Success 			204 "No Content"
// @Failure 			400 {object} errors.Error "Bad Request (e.g., invalid UUID format)"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found (Association does not exist)"
// @Failure 			500 {object} errors.Error "Internal Server Error (e.g., deletion failed)"
// @Router 				/service-local/{serviceId}/{localId}/ [delete]
func (a *Api) DeleteServiceLocal(c echo.Context) error {
	serviceId := c.Param("serviceId")
	localId := c.Param("localId")

	if serviceId == "" || localId == "" {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidServiceLocalId, c)
	}

	err := a.BllController.ServiceLocal.DeleteServiceLocal(serviceId, localId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.NoContent(http.StatusNoContent)
}

// @Summary 			Bulk Create ServiceLocals.
// @Description 		Creates multiple service-local associations.
// @Tags 				ServiceLocal
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request body schemas.BatchCreateServiceLocalRequest true "Bulk Create ServiceLocals Request"
// @Success 			201 {object} schemas.ServiceLocals "Created"
// @Failure 			400 {object} errors.Error "Bad Request (e.g., invalid updatedBy)"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found (Service or Local does not exist)"
// @Failure 			409 {object} errors.Error "Conflict (Association already exists)"
// @Failure 			422 {object} errors.Error "Unprocessable Entity (Invalid UUIDs or request body)"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/service-local/bulk/ [post]
func (a *Api) BulkCreateServiceLocals(c echo.Context) error {
	updatedBy := "ADMIN"

	var request schemas.BatchCreateServiceLocalRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, err := a.BllController.ServiceLocal.BulkCreateServiceLocals(
		request.ServiceLocals,
		updatedBy,
	)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusCreated, response)
}

// @Summary 			Fetch ServiceLocals.
// @Description 		Fetch all service-local associations, filtered by params.
// @Tags 				ServiceLocal
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param 				serviceId query string false "Service ID"
// @Param 				localId query string false "Local ID"
// @Success 			200 {object} schemas.ServiceLocals "OK"
// @Failure 			400 {object} errors.Error "Bad Request (e.g., invalid UUID format)"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found (Service or Local does not exist)"
// @Failure 			422 {object} errors.Error "Unprocessable Entity (Invalid UUIDs or request body)"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/service-local/ [get]
func (a *Api) FetchServiceLocals(c echo.Context) error {
	serviceId := c.QueryParam("serviceId")
	localId := c.QueryParam("localId")

	response, err := a.BllController.ServiceLocal.FetchServiceLocals(serviceId, localId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Bulk Delete ServiceLocals.
// @Description 		Bulk deletes service-local associations.
// @Tags 				ServiceLocal
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request	body   schemas.BulkDeleteServiceLocalRequest true  "Bulk Delete ServiceLocal Request"
// @Success 			204 "No Content"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/service-local/bulk/ [delete]
func (a *Api) BulkDeleteServiceLocals(c echo.Context) error {
	var request schemas.BulkDeleteServiceLocalRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	if err := a.BllController.ServiceLocal.BulkDeleteServiceLocals(request); err != nil {
		return errors.HandleError(*err, c)
	}

	return c.NoContent(http.StatusNoContent)
}
