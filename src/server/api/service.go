package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

// @Summary 			Get Service.
// @Description 		Gets a service given its id.
// @Tags 				Service
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               serviceId    path   string  true  "Service ID"
// @Success 			200 {object} schemas.Service "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/service/{serviceId}/ [get]
func (a *Api) GetService(c echo.Context) error {
	serviceId, parseErr := uuid.Parse(c.Param("serviceId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidServiceId, c)
	}

	response, err := a.BllController.Service.GetService(serviceId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Fetch Services.
// @Description 		Fetch all services, filtered by params.
// @Tags 				Service
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Success 			200 {object} schemas.Services "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/service/ [get]
func (a *Api) FetchServices(c echo.Context) error {
	response, err := a.BllController.Service.FetchServices()
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Create Service.
// @Description 		Create the service information.
// @Tags 				Service
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request body schemas.CreateServiceRequest true "Create Service Request"
// @Success 			201 {object} schemas.Service "Created"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/service/ [post]
func (a *Api) CreateService(c echo.Context) error {
	// TODO: Add access token validation (from here we will get the `updatedBy` param)
	updatedBy := "ADMIN"

	var request schemas.CreateServiceRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, newErr := a.BllController.Service.CreateService(request, updatedBy)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	return c.JSON(http.StatusCreated, response)
}

// @Summary 			Update Service.
// @Description 		Update the service information.
// @Tags 				Service
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               serviceId    path   string  true  "Service ID"
// @Param               request body schemas.UpdateServiceRequest true "Update Service Request"
// @Success 			200 {object} schemas.Service "Ok"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/service/{serviceId}/ [patch]
func (a *Api) UpdateService(c echo.Context) error {
	// TODO: Add access token validation (from here we will get the `updatedBy` param)
	updatedBy := "ADMIN"

	serviceId, parseErr := uuid.Parse(c.Param("serviceId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidServiceId, c)
	}

	var request schemas.UpdateServiceRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, newErr := a.BllController.Service.UpdateService(serviceId, request, updatedBy)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	return c.JSON(http.StatusOK, response)
}


// @Summary 			Delete Service.
// @Description 		Deletes a service given its id.
// @Tags 				Service
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               serviceId    path   string  true  "Service ID"
// @Success 			204 "No Content"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/service/{serviceId}/ [delete]
func (a *Api) DeleteService(c echo.Context) error {
	serviceId, parseErr := uuid.Parse(c.Param("serviceId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidServiceId, c)
	}

	if err := a.BllController.Service.DeleteService(serviceId); err != nil {
		return errors.HandleError(*err, c)
	}

	return c.NoContent(http.StatusNoContent)
}

