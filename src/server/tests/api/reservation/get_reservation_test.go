package reservation_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestGetReservationSuccessfully(t *testing.T) {
	/*
		GIVEN: A reservation exists in the database
		WHEN:  GET /reservation/{reservationId} is called with a valid reservation ID
		THEN:  A HTTP_200_OK status should be returned with the reservation data
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test reservation using factory
	reservation := factories.NewReservationModel(db, factories.ReservationModelF{})

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/reservation/"+reservation.Id.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Reservation
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains the correct reservation data
	assert.Equal(t, reservation.Id, response.Id)
	assert.Equal(t, reservation.UserId, response.UserId)
	assert.Equal(t, reservation.SessionId, response.SessionId)
}

func TestGetReservationNotFound(t *testing.T) {
	/*
		GIVEN: No reservation exists with the provided ID
		WHEN:  GET /reservation/{reservationId} is called with a non-existent reservation ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentReservationId := uuid.New()

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/reservation/"+nonExistentReservationId.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetReservationInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID format is provided
		WHEN:  GET /reservation/{reservationId} is called with an invalid UUID
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidReservationId := "invalid-uuid"

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/reservation/"+invalidReservationId+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
