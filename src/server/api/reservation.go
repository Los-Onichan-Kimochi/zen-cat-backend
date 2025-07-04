package api

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

// @Summary 			Get Reservation.
// @Description 		Gets a reservation given its id.
// @Tags 				Reservation
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               reservationId    path   string  true  "Reservation ID"
// @Success 			200 {object} schemas.Reservation "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/reservation/{reservationId}/ [get]
func (a *Api) GetReservation(c echo.Context) error {
	reservationId, parseErr := uuid.Parse(c.Param("reservationId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidReservationId, c)
	}

	response, err := a.BllController.Reservation.GetReservation(reservationId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Fetch Reservations.
// @Description 		Fetch all reservations, filtered by params.
// @Tags 				Reservation
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param 				userIds query []string false "User IDs"
// @Param 				sessionIds query []string false "Session IDs"
// @Param 				states query []string false "Reservation States"
// @Success 			200 {object} schemas.Reservations "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/reservation/ [get]
func (a *Api) FetchReservations(c echo.Context) error {
	userIdsString := c.QueryParam("userIds")
	sessionIdsString := c.QueryParam("sessionIds")
	statesString := c.QueryParam("states")

	userIds := []string{}
	if userIdsString != "" {
		userIds = strings.Split(userIdsString, ",")
	}

	sessionIds := []string{}
	if sessionIdsString != "" {
		sessionIds = strings.Split(sessionIdsString, ",")
	}

	states := []string{}
	if statesString != "" {
		states = strings.Split(statesString, ",")
	}

	response, err := a.BllController.Reservation.FetchReservations(userIds, sessionIds, states)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Create Reservation.
// @Description 		Create a new reservation.
// @Tags 				Reservation
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request body schemas.CreateReservationRequest true "Create Reservation Request"
// @Success 			201 {object} schemas.Reservation "Created"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/reservation/ [post]
func (a *Api) CreateReservation(c echo.Context) error {
	// TODO: Add access token validation (from here we will get the `updatedBy` param)
	updatedBy := "ADMIN"

	var request schemas.CreateReservationRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, newErr := a.BllController.Reservation.CreateReservation(request, updatedBy)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	return c.JSON(http.StatusCreated, response)
}

// @Summary 			Update Reservation.
// @Description 		Update an existing reservation.
// @Tags 				Reservation
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               reservationId    path   string  true  "Reservation ID"
// @Param               request body schemas.UpdateReservationRequest true "Update Reservation Request"
// @Success 			200 {object} schemas.Reservation "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/reservation/{reservationId}/ [patch]
func (a *Api) UpdateReservation(c echo.Context) error {
	// TODO: Add access token validation (from here we will get the `updatedBy` param)
	updatedBy := "ADMIN"

	reservationId, parseErr := uuid.Parse(c.Param("reservationId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidReservationId, c)
	}

	var request schemas.UpdateReservationRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, newErr := a.BllController.Reservation.UpdateReservation(reservationId, request, updatedBy)
	if newErr != nil {
		return errors.HandleError(*newErr, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Delete Reservation.
// @Description 		Delete a reservation given its id.
// @Tags 				Reservation
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               reservationId    path   string  true  "Reservation ID"
// @Success 			204 "No Content"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/reservation/{reservationId}/ [delete]
func (a *Api) DeleteReservation(c echo.Context) error {
	reservationId, parseErr := uuid.Parse(c.Param("reservationId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidReservationId, c)
	}

	if err := a.BllController.Reservation.DeleteReservation(reservationId); err != nil {
		return errors.HandleError(*err, c)
	}

	return c.NoContent(http.StatusNoContent)
}

// @Summary 			Bulk Delete Reservations.
// @Description 		Bulk delete reservations given their ids.
// @Tags 				Reservation
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request	body   schemas.BulkDeleteReservationRequest true  "Bulk Delete Reservation Request"
// @Success 			204 "No Content"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/reservation/bulk-delete/ [delete]
func (a *Api) BulkDeleteReservations(c echo.Context) error {
	var request schemas.BulkDeleteReservationRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	if err := a.BllController.Reservation.BulkDeleteReservations(request); err != nil {
		return errors.HandleError(*err, c)
	}

	return c.NoContent(http.StatusNoContent)
}

// @Summary 			Fetch Reservations by Community ID and User ID.
// @Description 		Fetch all reservations for a specific community and user.
// @Tags 				Reservation
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param 				communityId path string true "Community ID"
// @Param 				userId path string true "User ID"
// @Success 			200 {object} schemas.Reservations "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/reservation/{communityId}/{userId}/ [get]
func (a *Api) GetReservationsByCommunityIdByUserId(c echo.Context) error {
	communityId, parseErr := uuid.Parse(c.Param("communityId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidCommunityId, c)
	}

	userId, parseErr := uuid.Parse(c.Param("userId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidUserId, c)
	}

	response, err := a.BllController.Reservation.GetReservationsByCommunityIdByUserId(communityId, userId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}
