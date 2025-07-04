package api

import (
	"net/http"
	"strings"

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

// @Summary 			Get Service with image.
// @Description 		Gets a service given its id with image bytes.
// @Tags 				Service
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               serviceId    path   string  true  "Service ID"
// @Success 			200 {object} schemas.ServiceWithImage "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/service/{serviceId}/image/ [get]
func (a *Api) GetServiceWithImage(c echo.Context) error {
	serviceId, parseErr := uuid.Parse(c.Param("serviceId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidServiceId, c)
	}

	response, err := a.BllController.Service.GetService(serviceId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	var imageBytes *[]byte
	// Try to download image from S3, but don't fail if S3 is not available
	if response.ImageUrl != "" {
		downloadedBytes, s3Err := a.S3Service.DownloadFile(
			schemas.ServiceS3Prefix,
			response.ImageUrl,
		)
		if s3Err == nil {
			imageBytes = &downloadedBytes
		}
		// If S3 fails, we continue without image bytes
	}

	return c.JSON(http.StatusOK, schemas.ServiceWithImage{
		Service:    *response,
		ImageBytes: imageBytes,
	})
}

// @Summary 			Fetch Services.
// @Description 		Fetch all services, filtered by params.
// @Tags 				Service
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param 				ids query []string false "Service IDs"
// @Success 			200 {object} schemas.Services "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/service/ [get]
func (a *Api) FetchServices(c echo.Context) error {
	idsString := c.QueryParam("ids")

	ids := []string{}
	if idsString != "" {
		ids = strings.Split(idsString, ",")
	}

	response, err := a.BllController.Service.FetchServices(ids)
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

	if request.ImageUrl != "" {
		request.ImageUrl = a.S3Service.GenerateImageUrl(request.ImageUrl)
	}

	response, newErr := a.BllController.Service.CreateService(request, updatedBy)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	// Upload image to S3
	if response.ImageUrl != "" && request.ImageBytes != nil {
		a.S3Service.UploadFile(
			schemas.ServiceS3Prefix,
			response.ImageUrl,
			*request.ImageBytes,
		)
		// If S3 fails, we continue without image upload
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

	if request.ImageUrl != nil {
		*request.ImageUrl = a.S3Service.GenerateImageUrl(*request.ImageUrl)
	}

	response, newErr := a.BllController.Service.UpdateService(serviceId, request, updatedBy)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	// Upload image to S3 if it exists
	if request.ImageUrl != nil && request.ImageBytes != nil {
		err := a.S3Service.UploadFile(
			schemas.ServiceS3Prefix,
			response.ImageUrl,
			*request.ImageBytes,
		)
		if err != nil {
			return errors.HandleError(errors.InternalServerError.FailedToUploadImage, c)
		}
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

// @Summary 			Bulk Delete Services.
// @Description 		Bulk delete services given their ids.
// @Tags 				Service
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request	body   schemas.BulkDeleteServiceRequest true  "Bulk Delete Service Request"
// @Success 			204 "No Content"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/service/bulk-delete/ [delete]
func (a *Api) BulkDeleteServices(c echo.Context) error {
	var request schemas.BulkDeleteServiceRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	if err := a.BllController.Service.BulkDeleteServices(request); err != nil {
		return errors.HandleError(*err, c)
	}

	return c.NoContent(http.StatusNoContent)
}
