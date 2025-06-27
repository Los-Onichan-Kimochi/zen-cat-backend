package reservation_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestFetchReservationsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple reservations exist in the database
		WHEN:  GET /reservation/ is called
		THEN:  A HTTP_200_OK status should be returned with all reservations
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create test reservations using factory
	numReservations := 3
	factories.NewReservationModelBatch(db, numReservations)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/reservation/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Reservations
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, len(response.Reservations), numReservations)
}

func TestFetchReservationsWithFiltersSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple reservations exist in the database
		WHEN:  GET /reservation/ is called with user, session and state filters
		THEN:  A HTTP_200_OK status should be returned with the filtered reservations
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create test reservations using factory
	user1 := factories.NewUserModel(db, factories.UserModelF{})
	user2 := factories.NewUserModel(db, factories.UserModelF{})
	session1 := factories.NewSessionModel(db, factories.SessionModelF{})
	session2 := factories.NewSessionModel(db, factories.SessionModelF{})

	factories.NewReservationModel(db, factories.ReservationModelF{UserId: &user1.Id, SessionId: &session1.Id})
	factories.NewReservationModel(db, factories.ReservationModelF{UserId: &user1.Id, SessionId: &session2.Id})
	factories.NewReservationModel(db, factories.ReservationModelF{UserId: &user2.Id, SessionId: &session1.Id})

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/reservation/?userIds="+user1.Id.String()+"&sessionIds="+session1.Id.String(), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Reservations
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Len(t, response.Reservations, 1)
	assert.Equal(t, user1.Id, response.Reservations[0].UserId)
	assert.Equal(t, session1.Id, response.Reservations[0].SessionId)
}

func TestFetchReservationsEmpty(t *testing.T) {
	/*
		GIVEN: No reservations exist in the database
		WHEN:  GET /reservation/ is called
		THEN:  A HTTP_200_OK status should be returned with an empty array
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/reservation/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Reservations
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Empty(t, response.Reservations)
}
