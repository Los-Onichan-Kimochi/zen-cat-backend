package api

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
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
