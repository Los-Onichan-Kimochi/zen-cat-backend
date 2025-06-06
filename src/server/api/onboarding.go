package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

// @Summary 			Get Onboarding by ID.
// @Description 		Gets an onboarding given its id.
// @Tags 				Onboarding
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               onboardingId    path   string  true  "Onboarding ID"
// @Success 			200 {object} schemas.Onboarding "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/onboarding/{onboardingId}/ [get]
func (a *Api) GetOnboarding(c echo.Context) error {
	onboardingId, parseErr := uuid.Parse(c.Param("onboardingId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidOnboardingId, c)
	}

	response, err := a.BllController.Onboarding.GetOnboarding(onboardingId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Get Onboarding by User ID.
// @Description 		Gets an onboarding given its user id.
// @Tags 				Onboarding
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               userId    path   string  true  "User ID"
// @Success 			200 {object} schemas.Onboarding "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/onboarding/user/{userId}/ [get]
func (a *Api) GetOnboardingByUserId(c echo.Context) error {
	userId, parseErr := uuid.Parse(c.Param("userId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidUserId, c)
	}

	response, err := a.BllController.Onboarding.GetOnboardingByUserId(userId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Fetch Onboardings.
// @Description 		Fetch all onboardings.
// @Tags 				Onboarding
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Success 			200 {object} schemas.Onboardings "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/onboarding/ [get]
func (a *Api) FetchOnboardings(c echo.Context) error {
	response, err := a.BllController.Onboarding.FetchOnboardings()
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Create Onboarding for User.
// @Description 		Creates a new onboarding for a specific user.
// @Tags 				Onboarding
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               userId    path   string  true  "User ID"
// @Param               createOnboardingRequest    body   schemas.CreateOnboardingRequest  true  "Create Onboarding Request"
// @Success 			201 {object} schemas.Onboarding "Created"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/onboarding/user/{userId}/ [post]
func (a *Api) CreateOnboardingForUser(c echo.Context) error {
	updatedBy := "USER" // Podría ser obtenido del JWT token en producción

	userId, parseErr := uuid.Parse(c.Param("userId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidUserId, c)
	}

	var request schemas.CreateOnboardingRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, err := a.BllController.Onboarding.CreateOnboardingForUser(userId, request, updatedBy)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusCreated, response)
}

// @Summary 			Update Onboarding by ID.
// @Description 		Updates an onboarding given its id.
// @Tags 				Onboarding
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               onboardingId    path   string  true  "Onboarding ID"
// @Param               updateOnboardingRequest    body   schemas.UpdateOnboardingRequest  true  "Update Onboarding Request"
// @Success 			200 {object} schemas.Onboarding "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/onboarding/{onboardingId}/ [patch]
func (a *Api) UpdateOnboarding(c echo.Context) error {
	updatedBy := "USER" // Podría ser obtenido del JWT token en producción

	onboardingId, parseErr := uuid.Parse(c.Param("onboardingId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidOnboardingId, c)
	}

	var request schemas.UpdateOnboardingRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, err := a.BllController.Onboarding.UpdateOnboarding(onboardingId, request, updatedBy)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Update Onboarding by User ID.
// @Description 		Updates an onboarding given its user id.
// @Tags 				Onboarding
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               userId    path   string  true  "User ID"
// @Param               updateOnboardingRequest    body   schemas.UpdateOnboardingRequest  true  "Update Onboarding Request"
// @Success 			200 {object} schemas.Onboarding "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/onboarding/user/{userId}/ [patch]
func (a *Api) UpdateOnboardingByUserId(c echo.Context) error {
	updatedBy := "USER" // Podría ser obtenido del JWT token en producción

	userId, parseErr := uuid.Parse(c.Param("userId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidUserId, c)
	}

	var request schemas.UpdateOnboardingRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, err := a.BllController.Onboarding.UpdateOnboardingByUserId(userId, request, updatedBy)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Delete Onboarding by ID.
// @Description 		Deletes an onboarding given its id.
// @Tags 				Onboarding
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               onboardingId    path   string  true  "Onboarding ID"
// @Success 			204 "No Content"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/onboarding/{onboardingId}/ [delete]
func (a *Api) DeleteOnboarding(c echo.Context) error {
	onboardingId, parseErr := uuid.Parse(c.Param("onboardingId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidOnboardingId, c)
	}

	if err := a.BllController.Onboarding.DeleteOnboarding(onboardingId); err != nil {
		return errors.HandleError(*err, c)
	}

	return c.NoContent(http.StatusNoContent)
}

// @Summary 			Delete Onboarding by User ID.
// @Description 		Deletes an onboarding given its user id.
// @Tags 				Onboarding
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               userId    path   string  true  "User ID"
// @Success 			204 "No Content"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/onboarding/user/{userId}/ [delete]
func (a *Api) DeleteOnboardingByUserId(c echo.Context) error {
	userId, parseErr := uuid.Parse(c.Param("userId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidUserId, c)
	}

	if err := a.BllController.Onboarding.DeleteOnboardingByUserId(userId); err != nil {
		return errors.HandleError(*err, c)
	}

	return c.NoContent(http.StatusNoContent)
}
