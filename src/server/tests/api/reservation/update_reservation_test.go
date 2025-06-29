package reservation_test

import (
	"bytes"
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

func TestUpdateReservationSuccessfully(t *testing.T) {
	/*
		GIVEN: A reservation exists in the database
		WHEN:  PATCH /reservation/{reservationId}/ is called with valid update data
		THEN:  A HTTP_200_OK status should be returned with the updated reservation
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test reservation using factory
	reservation := factories.NewReservationModel(db, factories.ReservationModelF{})
	newSession := factories.NewSessionModel(db, factories.SessionModelF{})

	// Prepare update request
	updateReservationRequest := schemas.UpdateReservationRequest{
		SessionId: &newSession.Id,
	}

	requestBody, _ := json.Marshal(updateReservationRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/reservation/"+reservation.Id.String()+"/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Reservation
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains the updated reservation data
	assert.Equal(t, reservation.Id, response.Id)
	assert.Equal(t, *updateReservationRequest.SessionId, response.SessionId)
}

func TestUpdateReservationNotFound(t *testing.T) {
	/*
		GIVEN: No reservation exists with the provided ID
		WHEN:  PATCH /reservation/{reservationId}/ is called with a non-existent reservation ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentReservationId := uuid.New()

	// Prepare update request
	newSessionId := uuid.New()
	updateReservationRequest := schemas.UpdateReservationRequest{
		SessionId: &newSessionId,
	}

	requestBody, _ := json.Marshal(updateReservationRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/reservation/"+nonExistentReservationId.String()+"/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateReservationInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID is provided
		WHEN:  PATCH /reservation/{reservationId}/ is called with an invalid UUID
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidReservationId := "invalid-uuid"

	// Prepare update request
	newSessionId := uuid.New()
	updateReservationRequest := schemas.UpdateReservationRequest{
		SessionId: &newSessionId,
	}

	requestBody, _ := json.Marshal(updateReservationRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/reservation/"+invalidReservationId+"/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestUpdateReservationInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: A reservation exists but the request body is invalid
		WHEN:  PATCH /reservation/{reservationId}/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test reservation
	reservation := factories.NewReservationModel(db, factories.ReservationModelF{})
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/reservation/"+reservation.Id.String()+"/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
