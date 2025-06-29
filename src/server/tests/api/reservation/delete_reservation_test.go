package reservation_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestDeleteReservationSuccessfully(t *testing.T) {
	/*
		GIVEN: A reservation exists in the database
		WHEN:  DELETE /reservation/{reservationId}/ is called with a valid reservation ID
		THEN:  A HTTP_204_NO_CONTENT status should be returned and the reservation should be deleted
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test reservation using factory
	reservation := factories.NewReservationModel(db)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/reservation/"+reservation.Id.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Verify the reservation was deleted
	var count int64
	db.Model(&reservation).Where("id = ?", reservation.Id).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestDeleteReservationNotFound(t *testing.T) {
	/*
		GIVEN: No reservation exists with the provided ID
		WHEN:  DELETE /reservation/{reservationId}/ is called with a non-existent reservation ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentReservationId := uuid.New()

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/reservation/"+nonExistentReservationId.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteReservationInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID is provided
		WHEN:  DELETE /reservation/{reservationId}/ is called with an invalid UUID
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidReservationId := "invalid-uuid"

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/reservation/"+invalidReservationId+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
