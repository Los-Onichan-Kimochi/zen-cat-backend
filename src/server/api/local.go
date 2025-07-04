package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

// @Summary 			Get Local.
// @Description 		Gets a local given its id.
// @Tags 				Local
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               localId    path   string  true  "Local ID"
// @Success 			200 {object} schemas.Local "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/local/{localId}/ [get]
func (a *Api) GetLocal(c echo.Context) error {
	localId, parseErr := uuid.Parse(c.Param("localId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidLocalId, c)
	}

	response, err := a.BllController.Local.GetLocal(localId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Get Local with image.
// @Description 		Gets a local given its id with image bytes.
// @Tags 				Local
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               localId    path   string  true  "Local ID"
// @Success 			200 {object} schemas.LocalWithImage "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/local/{localId}/image/ [get]
func (a *Api) GetLocalWithImage(c echo.Context) error {
	localId, parseErr := uuid.Parse(c.Param("localId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidLocalId, c)
	}

	response, err := a.BllController.Local.GetLocal(localId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	var imageBytes *[]byte
	// Try to download image from S3, but don't fail if S3 is not available
	if response.ImageUrl != "" {
		downloadedBytes, s3Err := a.S3Service.DownloadFile(
			schemas.LocalS3Prefix,
			response.ImageUrl,
		)
		if s3Err == nil {
			imageBytes = &downloadedBytes
		}
		// If S3 fails, we continue without image bytes
	}

	return c.JSON(http.StatusOK, schemas.LocalWithImage{
		Local:      *response,
		ImageBytes: imageBytes,
	})
}

// @Summary 			Fetch Locals.
// @Description 		Fetches all locals.
// @Tags 				Local
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Success 			200 {object} schemas.Locals "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/local/ [get]
func (a *Api) FetchLocals(c echo.Context) error {
	response, err := a.BllController.Local.FetchLocals()
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Create Local.
// @Description 		Create the local information.
// @Tags 				Local
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request body schemas.CreateLocalRequest true "Create Local Request"
// @Success 			201 {object} schemas.Local "Created"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/local/ [post]
func (a *Api) CreateLocal(c echo.Context) error {
	// TODO: Add access token validation (from here we will get the `updatedBy` param)
	updatedBy := "ADMIN"

	var request schemas.CreateLocalRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}
	if request.ImageUrl != "" {
		request.ImageUrl = a.S3Service.GenerateImageUrl(request.ImageUrl)
	}

	response, newErr := a.BllController.Local.CreateLocal(request, updatedBy)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	// Upload image to S3
	if response.ImageUrl != "" && request.ImageBytes != nil {
		a.S3Service.UploadFile(
			schemas.LocalS3Prefix,
			response.ImageUrl,
			*request.ImageBytes,
		)
		// If S3 fails, we continue without image upload
	}

	return c.JSON(http.StatusCreated, response)
}

// @Summary 			Bulk Create Locals.
// @Description 		Creates multiple locals in a batch.
// @Tags 				Local
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request body schemas.BatchCreateLocalRequest true "Bulk Create Locals Request"
// @Success 			201 {object} schemas.Locals "Created"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/local/bulk-create/ [post]
func (a *Api) BulkCreateLocals(c echo.Context) error {
	updatedBy := "ADMIN"

	var request schemas.BatchCreateLocalRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	// Check for missing Locals field which means invalid JSON binding
	if request.Locals == nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, newErr := a.BllController.Local.BulkCreateLocals(
		request.Locals,
		updatedBy,
	)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	return c.JSON(http.StatusCreated, response)
}

// @Summary 			Update Local.
// @Description 		Update the local information.
// @Tags 				Local
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               localId    path   string  true  "Local ID"
// @Param               request body schemas.UpdateLocalRequest true "Update Local Request"
// @Success 			200 {object} schemas.Local "Ok"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/local/{localId}/ [patch]
func (a *Api) UpdateLocal(c echo.Context) error {
	// TODO: Add access token validation (from here we will get the `updatedBy` param)
	updatedBy := "ADMIN"

	localId, parseErr := uuid.Parse(c.Param("localId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidLocalId, c)
	}

	var request schemas.UpdateLocalRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}
	if request.ImageUrl != nil {
		*request.ImageUrl = a.S3Service.GenerateImageUrl(*request.ImageUrl)
	}

	response, newErr := a.BllController.Local.UpdateLocal(localId, request, updatedBy)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	// Upload image to S3 if it exists
	if request.ImageUrl != nil && request.ImageBytes != nil {
		err := a.S3Service.UploadFile(
			schemas.LocalS3Prefix,
			response.ImageUrl,
			*request.ImageBytes,
		)
		if err != nil {
			return errors.HandleError(errors.InternalServerError.FailedToUploadImage, c)
		}
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Delete Local.
// @Description 		Deletes a local given its id.
// @Tags 				Local
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               localId    path   string  true  "Local ID"
// @Success 			204 "No Content"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/local/{localId}/ [delete]
func (a *Api) DeleteLocal(c echo.Context) error {
	localId, parseErr := uuid.Parse(c.Param("localId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidLocalId, c)
	}

	if err := a.BllController.Local.DeleteLocal(localId); err != nil {
		return errors.HandleError(*err, c)
	}

	return c.NoContent(http.StatusNoContent)
}

// @Summary 			Bulk Delete Locals.
// @Description 		Bulk delete locals given their ids.
// @Tags 				Local
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request	body   schemas.BulkDeleteLocalRequest true  "Bulk Delete Local Request"
// @Success 			204 "No Content"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/local/bulk-delete/ [delete]
func (a *Api) BulkDeleteLocals(c echo.Context) error {
	var request schemas.BulkDeleteLocalRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	// Check for missing Locals field which means invalid JSON binding
	if request.Locals == nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	if err := a.BllController.Local.BulkDeleteLocals(request); err != nil {
		return errors.HandleError(*err, c)
	}
	return c.NoContent(http.StatusNoContent)
}
