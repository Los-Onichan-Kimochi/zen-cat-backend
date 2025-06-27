package reservation_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestBulkDeleteReservationsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple reservations exist in the database
		WHEN:  DELETE /reservation/bulk-delete/ is called with valid reservation IDs
		THEN:  A HTTP_204_NO_CONTENT status should be returned and the reservations should be deleted
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create test reservations using factory
	numReservations := 3
	testReservations := factories.NewReservationModelBatch(db, numReservations)

	// Extract reservation IDs
	reservationIds := make([]string, numReservations)
	for i, reservation := range testReservations {
		reservationIds[i] = reservation.Id.String()
	}

	bulkDeleteRequest := schemas.BulkDeleteReservationRequest{
		Reservations: reservationIds,
	}

	requestBody, _ := json.Marshal(bulkDeleteRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/reservation/bulk-delete/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Verify the reservations were deleted
	for _, reservation := range testReservations {
		var count int64
		db.Model(&reservation).Where("id = ?", reservation.Id).Count(&count)
		assert.Equal(t, int64(0), count)
	}
}

func TestBulkDeleteReservationsEmptyList(t *testing.T) {
	/*
		GIVEN: A bulk reservation deletion request with an empty list
		WHEN:  DELETE /reservation/bulk-delete/ is called with an empty reservations list
		THEN:  A HTTP_204_NO_CONTENT status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	bulkDeleteRequest := schemas.BulkDeleteReservationRequest{
		Reservations: []string{},
	}

	requestBody, _ := json.Marshal(bulkDeleteRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/reservation/bulk-delete/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestBulkDeleteReservationsInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  DELETE /reservation/bulk-delete/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/reservation/bulk-delete/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
