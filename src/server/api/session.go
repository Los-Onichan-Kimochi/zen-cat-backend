package api

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

// @Summary 			Get Session.
// @Description 		Gets a session given its id.
// @Tags 				Session
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               sessionId    path   string  true  "Session ID"
// @Success 			200 {object} schemas.Session "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/session/{sessionId}/ [get]
func (a *Api) GetSession(c echo.Context) error {
	sessionId, parseErr := uuid.Parse(c.Param("sessionId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidSessionId, c)
	}

	response, err := a.BllController.Session.GetSession(sessionId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Fetch Sessions.
// @Description 		Fetch all sessions, filtered by params.
// @Tags 				Session
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param 				professionalIds query []string false "Professional IDs"
// @Param 				localIds query []string false "Local IDs"
// @Param 				communityServiceIds query []string false "Community Service IDs"
// @Param 				states query []string false "Session States"
// @Success 			200 {object} schemas.Sessions "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/session/ [get]
func (a *Api) FetchSessions(c echo.Context) error {
	professionalIdsString := c.QueryParam("professionalIds")
	localIdsString := c.QueryParam("localIds")
	communityServiceIdsString := c.QueryParam("communityServiceIds")
	statesString := c.QueryParam("states")

	professionalIds := []string{}
	if professionalIdsString != "" {
		professionalIds = strings.Split(professionalIdsString, ",")
	}

	localIds := []string{}
	if localIdsString != "" {
		localIds = strings.Split(localIdsString, ",")
	}
	
	communityServiceIds := []string{}
	if communityServiceIdsString != "" {
		communityServiceIds = strings.Split(communityServiceIdsString, ",")
	}

	states := []string{}
	if statesString != "" {
		states = strings.Split(statesString, ",")
	}

	response, err := a.BllController.Session.FetchSessions(professionalIds, localIds, communityServiceIds, states)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Create Session.
// @Description 		Create the session information.
// @Tags 				Session
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request body schemas.CreateSessionRequest true "Create Session Request"
// @Success 			201 {object} schemas.Session "Created"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/session/ [post]
func (a *Api) CreateSession(c echo.Context) error {
	// TODO: Add access token validation (from here we will get the `updatedBy` param)
	updatedBy := "ADMIN"

	var request schemas.CreateSessionRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, newErr := a.BllController.Session.CreateSession(request, updatedBy)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	return c.JSON(http.StatusCreated, response)
}

// @Summary 			Update Session.
// @Description 		Update the session information.
// @Tags 				Session
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               sessionId    path   string  true  "Session ID"
// @Param               request body schemas.UpdateSessionRequest true "Update Session Request"
// @Success 			200 {object} schemas.Session "Ok"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/session/{sessionId}/ [patch]
func (a *Api) UpdateSession(c echo.Context) error {
	// TODO: Add access token validation (from here we will get the `updatedBy` param)
	updatedBy := "ADMIN"

	sessionId, parseErr := uuid.Parse(c.Param("sessionId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidSessionId, c)
	}

	var request schemas.UpdateSessionRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, newErr := a.BllController.Session.UpdateSession(sessionId, request, updatedBy)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Delete Session.
// @Description 		Deletes a session given its id.
// @Tags 				Session
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               sessionId    path   string  true  "Session ID"
// @Success 			204 {object} schemas.Session "No Content"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/session/{sessionId}/ [delete]
func (a *Api) DeleteSession(c echo.Context) error {
	sessionId, parseErr := uuid.Parse(c.Param("sessionId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidSessionId, c)
	}

	if err := a.BllController.Session.DeleteSession(sessionId); err != nil {
		return errors.HandleError(*err, c)
	}

	return c.NoContent(http.StatusNoContent)
}

// @Summary 			Bulk Delete Sessions.
// @Description 		Bulk deletes sessions given their ids.
// @Tags 				Session
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request	body   schemas.BulkDeleteSessionRequest true  "Bulk Delete Session Request"
// @Success 			204 {object} schemas.Session "No Content"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/session/bulk-delete/ [delete]
func (a *Api) BulkDeleteSessions(c echo.Context) error {
	var request schemas.BulkDeleteSessionRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	if err := a.BllController.Session.BulkDeleteSessions(request); err != nil {
		return errors.HandleError(*err, c)
	}

	return c.NoContent(http.StatusNoContent)
}

// @Summary 			Bulk Create Session.
// @Description 		Create multiple sessions in a single.
// @Tags 				Session
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request body schemas.BatchCreateSessionRequest true "Bulk Create Sessions Request"
// @Success 			201 {object} schemas.Sessions "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/session/bulk/ [post]
func (a *Api) BulkCreateSessions(c echo.Context) error {
	updatedBy := "ADMIN"

	var request schemas.BatchCreateSessionRequest

	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, newErr := a.BllController.Session.BulkCreateSessions(
		request.Sessions,
		updatedBy,
	)

	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	return c.JSON(http.StatusCreated, response)
}

// @Summary 			Check Session Conflicts.
// @Description 		Check for time conflicts with existing sessions.
// @Tags 				Session
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request body schemas.CheckConflictRequest true "Check Conflict Request"
// @Success 			200 {object} schemas.ConflictResult "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/session/check-conflicts/ [post]
func (a *Api) CheckSessionConflicts(c echo.Context) error {
	var request schemas.CheckConflictRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, newErr := a.BllController.Session.CheckConflicts(request)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Get Day Availability.
// @Description 		Get availability information for a specific date.
// @Tags 				Session
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request body schemas.AvailabilityRequest true "Availability Request"
// @Success 			200 {object} schemas.AvailabilityResult "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/session/availability/ [post]
func (a *Api) GetDayAvailability(c echo.Context) error {
	var request schemas.AvailabilityRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, newErr := a.BllController.Session.GetAvailability(request)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	return c.JSON(http.StatusOK, response)
}
