package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

// @Summary 			Get Community.
// @Description 		Gets a community given its id.
// @Tags 				Community
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               communityId    path   string  true  "Community ID"
// @Success 			200 {object} schemas.Community "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/community/{communityId}/ [get]
func (a *Api) GetCommunity(c echo.Context) error {
	communityId, parseErr := uuid.Parse(c.Param("communityId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidCommunityId, c)
	}

	response, err := a.BllController.Community.GetCommunity(communityId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Get Community with Image.
// @Description 		Gets a community given its id with its image.
// @Tags 				Community
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               communityId    path   string  true  "Community ID"
// @Success 			200 {object} schemas.CommunityWithImage "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/community/{communityId}/image/ [get]
func (a *Api) GetCommunityWithImage(c echo.Context) error {
	communityId, parseErr := uuid.Parse(c.Param("communityId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidCommunityId, c)
	}

	response, err := a.BllController.Community.GetCommunity(communityId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	imageBytes, s3Err := a.S3Service.DownloadFile(schemas.CommunityS3Prefix, response.ImageUrl)
	if s3Err != nil && s3Err.Error() != "NoSuchKey" {
		return errors.HandleError(errors.InternalServerError.FailedToDownloadImage, c)
	}

	return c.JSON(http.StatusOK, schemas.CommunityWithImage{
		Community:  *response,
		ImageBytes: imageBytes,
	})
}

// @Summary 			Fetch Communities.
// @Description 		Fetch all communities, filtered by params.
// @Tags 				Community
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Success 			200 {object} schemas.Communities "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/community/ [get]
func (a *Api) FetchCommunities(c echo.Context) error {
	response, err := a.BllController.Community.FetchCommunities()
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Create Community.
// @Description 		Create the community information.
// @Tags 				Community
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request body schemas.CreateCommunityRequest true "Create Community Request"
// @Success 			201 {object} schemas.Community "Created"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/community/ [post]
func (a *Api) CreateCommunity(c echo.Context) error {
	// TODO: Add access token validation (from here we will get the `updatedBy` param)
	updatedBy := "ADMIN"

	var request schemas.CreateCommunityRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, newErr := a.BllController.Community.CreateCommunity(request, updatedBy)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	// Upload image to S3
	err := a.S3Service.UploadFile(schemas.CommunityS3Prefix, request.ImageUrl, request.ImageBytes)
	if err != nil {
		return errors.HandleError(errors.InternalServerError.FailedToUploadImage, c)
	}

	return c.JSON(http.StatusCreated, response)
}

// @Summary 			Bulk Create Community.
// @Description 		Create multiple communities in a single.
// @Tags 				Community
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request body schemas.BatchCreateCommunityRequest true "Bulk Create Communities Request"
// @Success 			201 {object} schemas.Communities "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/community/bulk-create/ [post]
func (a *Api) BulkCreateCommunities(c echo.Context) error {
	updatedBy := "ADMIN"

	var request schemas.BatchCreateCommunityRequest

	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, newErr := a.BllController.Community.BulkCreateCommunities(
		request.Communities,
		updatedBy,
	)

	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	return c.JSON(http.StatusCreated, response)
}

// @Summary 			Update Community.
// @Description 		Update the community information.
// @Tags 				Community
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               communityId    path   string  true  "Community ID"
// @Param               request body schemas.UpdateCommunityRequest true "Update Community Request"
// @Success 			200 {object} schemas.Community "Ok"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/community/{communityId}/ [patch]
func (a *Api) UpdateCommunity(c echo.Context) error {
	// TODO: Add access token validation (from here we will get the `updatedBy` param)
	updatedBy := "ADMIN"

	communityId, parseErr := uuid.Parse(c.Param("communityId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidCommunityId, c)
	}

	var request schemas.UpdateCommunityRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, newErr := a.BllController.Community.UpdateCommunity(communityId, request, updatedBy)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Delete Community.
// @Description 		Deletes a community given its id.
// @Tags 				Community
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               communityId    path   string  true  "Community ID"
// @Success 			204 "No Content"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/community/{communityId}/ [delete]
func (a *Api) DeleteCommunity(c echo.Context) error {
	communityId, parseErr := uuid.Parse(c.Param("communityId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidCommunityId, c)
	}

	if err := a.BllController.Community.DeleteCommunity(communityId); err != nil {
		return errors.HandleError(*err, c)
	}

	return c.NoContent(http.StatusNoContent)
}

// @Summary 			Bulk Delete Communities.
// @Description 		Bulk deletes communities given their ids.
// @Tags 				Community
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request	body   schemas.BulkDeleteCommunityRequest true  "Bulk Delete Community Request"
// @Success 			204 {object} schemas.Community "No Content"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/community/bulk-delete/ [delete]
func (a *Api) BulkDeleteCommunities(c echo.Context) error {
	var request schemas.BulkDeleteCommunityRequest

	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	if err := a.BllController.Community.BulkDeleteCommunities(request); err != nil {
		return errors.HandleError(*err, c)
	}

	return c.NoContent(http.StatusNoContent)
}
