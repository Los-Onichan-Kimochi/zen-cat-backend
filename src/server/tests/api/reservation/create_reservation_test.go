package reservation_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestCreateReservationSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid data for creating a reservation
		WHEN:  POST /reservation/ is called with valid reservation data
		THEN:  A HTTP_201_CREATED status should be returned with the created reservation
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create dependencies using factories
	user := factories.NewUserModel(db, factories.UserModelF{})
	date := time.Now().Add(24 * time.Hour)
	session := factories.NewSessionModel(db, factories.SessionModelF{
		Date: &date,
	})

	createReservationRequest := schemas.CreateReservationRequest{
		UserId:    user.Id,
		SessionId: session.Id,
	}

	requestBody, _ := json.Marshal(createReservationRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/reservation/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response schemas.Reservation
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response
	assert.NotEmpty(t, response.Id)
	assert.Equal(t, createReservationRequest.UserId, response.UserId)
	assert.Equal(t, createReservationRequest.SessionId, response.SessionId)

	// Verify the reservation was created in the database
	var dbReservation model.Reservation
	err = db.First(&dbReservation, "id = ?", response.Id).Error
	assert.NoError(t, err)
	assert.Equal(t, response.Id, dbReservation.Id)
}

func TestCreateReservationUserNotFound(t *testing.T) {
	/*
		GIVEN: A non-existent user ID
		WHEN:  POST /reservation/ is called with invalid user ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a session for the reservation
	session := factories.NewSessionModel(db, factories.SessionModelF{})

	createReservationRequest := schemas.CreateReservationRequest{
		UserId:    uuid.New(), // Non-existent user
		SessionId: session.Id,
	}

	requestBody, _ := json.Marshal(createReservationRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/reservation/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCreateReservationSessionNotFound(t *testing.T) {
	/*
		GIVEN: A non-existent session ID
		WHEN:  POST /reservation/ is called with invalid session ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a user for the reservation
	user := factories.NewUserModel(db, factories.UserModelF{})

	createReservationRequest := schemas.CreateReservationRequest{
		UserId:    user.Id,
		SessionId: uuid.New(), // Non-existent session
	}

	requestBody, _ := json.Marshal(createReservationRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/reservation/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCreateReservationInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  POST /reservation/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/reservation/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

// Helper function to create string pointers
func strPtr(s string) *string {
	return &s
}
