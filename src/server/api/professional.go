package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

// @Summary 			Get Professional.
// @Description 		Gets a professional given its id.
// @Tags 				Professional
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               professionalId    path   string  true  "Professional ID"
// @Success 			200 {object} schemas.Professional "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/professional/{professionalId}/ [get]
func (a *Api) GetProfessional(c echo.Context) error {
	professionalId, parseErr := uuid.Parse(c.Param("professionalId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidProfessionalId, c)
	}

	response, err := a.BllController.Professional.GetProfessional(professionalId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Get Professional with image.
// @Description 		Gets a professional given its id with image bytes.
// @Tags 				Professional
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               professionalId    path   string  true  "Professional ID"
// @Success 			200 {object} schemas.ProfessionalWithImage "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/professional/{professionalId}/image/ [get]
func (a *Api) GetProfessionalWithImage(c echo.Context) error {
	professionalId, parseErr := uuid.Parse(c.Param("professionalId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidProfessionalId, c)
	}

	response, err := a.BllController.Professional.GetProfessional(professionalId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	var imageBytes *[]byte
	// Try to download image from S3, but don't fail if S3 is not available
	if response.ImageUrl != "" {
		downloadedBytes, s3Err := a.S3Service.DownloadFile(
			schemas.ProfessionalS3Prefix,
			response.ImageUrl,
		)
		if s3Err == nil {
			imageBytes = &downloadedBytes
		}
		// If S3 fails, we continue without image bytes
	}

	return c.JSON(http.StatusOK, schemas.ProfessionalWithImage{
		Professional: *response,
		ImageBytes:   imageBytes,
	})
}

// @Summary 			Fetch Professionals.
// @Description 		Fetch all professionals, filtered by params.
// @Tags 				Professional
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Success 			200 {object} schemas.Professionals "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/professional/ [get]
func (a *Api) FetchProfessionals(c echo.Context) error {
	response, err := a.BllController.Professional.FetchProfessionals()
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Create Professional.
// @Description 		Creates a new professional.
// @Tags 				Professional
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request	body   schemas.CreateProfessionalRequest true  "Create Professional Request"
// @Success 			201 {object} schemas.Professional "Created"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/professional/ [post]
func (a *Api) CreateProfessional(c echo.Context) error {
	updateBy := "ADMIN"

	var request schemas.CreateProfessionalRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	if request.ImageUrl != "" {
		request.ImageUrl = a.S3Service.GenerateImageUrl(request.ImageUrl)
	}

	response, err := a.BllController.Professional.CreateProfessional(request, updateBy)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	// Upload image to S3
	if response.ImageUrl != "" && request.ImageBytes != nil {
		a.S3Service.UploadFile(
			schemas.ProfessionalS3Prefix,
			response.ImageUrl,
			*request.ImageBytes,
		)
		// If S3 fails, we continue without image upload
	}

	return c.JSON(http.StatusCreated, response)
}

// @Summary 			Update Professional.
// @Description 		Updates a professional given its id.
// @Tags 				Professional
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               professionalId    path   string  true  "Professional ID"
// @Param               request	body   schemas.UpdateProfessionalRequest true  "Update Professional Request"
// @Success 			200 {object} schemas.Professional "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/professional/{professionalId}/ [patch]
func (a *Api) UpdateProfessional(c echo.Context) error {
	updateBy := "ADMIN"

	professionalId, parseErr := uuid.Parse(c.Param("professionalId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidProfessionalId, c)
	}

	var request schemas.UpdateProfessionalRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	if request.ImageUrl != nil {
		*request.ImageUrl = a.S3Service.GenerateImageUrl(*request.ImageUrl)
	}

	response, err := a.BllController.Professional.UpdateProfessional(
		professionalId,
		request,
		updateBy,
	)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	// Upload image to S3 if it exists
	if request.ImageUrl != nil && request.ImageBytes != nil {
		err := a.S3Service.UploadFile(
			schemas.ProfessionalS3Prefix,
			response.ImageUrl,
			*request.ImageBytes,
		)
		if err != nil {
			return errors.HandleError(errors.InternalServerError.FailedToUploadImage, c)
		}
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Delete Professional.
// @Description 		Deletes a professional given its id.
// @Tags 				Professional
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               professionalId    path   string  true  "Professional ID"
// @Success 			204 {object} schemas.Professional "No Content"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/professional/{professionalId}/ [delete]
func (a *Api) DeleteProfessional(c echo.Context) error {
	professionalId, parseErr := uuid.Parse(c.Param("professionalId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidProfessionalId, c)
	}

	if err := a.BllController.Professional.DeleteProfessional(professionalId); err != nil {
		return errors.HandleError(*err, c)
	}

	return c.NoContent(http.StatusNoContent)
}

// @Summary 			Bulk Create Professionals.
// @Description 		Creates multiple professionals in a batch.
// @Tags 				Professional
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request	body   schemas.BulkCreateProfessionalRequest true  "Bulk Create Professional Request"
// @Success 			201 {object} schemas.Professionals "Created"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/professional/bulk-create/ [post]
func (a *Api) BulkCreateProfessionals(c echo.Context) error {
	var request schemas.BulkCreateProfessionalRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	updateBy := "ADMIN"
	response, err := a.BllController.Professional.BulkCreateProfessionals(
		request.Professionals,
		updateBy,
	)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusCreated, response)
}

// @Summary 			Bulk Delete Professionals.
// @Description 		Bulk deletes professionals given their ids.
// @Tags 				Professional
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request	body   schemas.BulkDeleteProfessionalRequest true  "Bulk Delete Professional Request"
// @Success 			204 {object} schemas.Professional "No Content"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/professional/bulk-delete/ [delete]
func (a *Api) BulkDeleteProfessionals(c echo.Context) error {
	var request schemas.BulkDeleteProfessionalRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	if err := a.BllController.Professional.BulkDeleteProfessionals(request); err != nil {
		return errors.HandleError(*err, c)
	}

	return c.NoContent(http.StatusNoContent)
}
