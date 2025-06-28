package service_local_test

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

func TestBulkDeleteServiceLocalsSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid service-local associations exist
		WHEN:  DELETE /service-local/bulk/ is called with valid associations
		THEN:  A HTTP_200_OK status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create service-local associations using factories
	serviceLocal1 := factories.NewServiceLocalModel(db)
	serviceLocal2 := factories.NewServiceLocalModel(db)

	bulkRequest := schemas.BulkDeleteServiceLocalRequest{
		ServiceLocals: []*schemas.DeleteServiceLocalRequest{
			{
				ServiceId: serviceLocal1.ServiceId,
				LocalId:   serviceLocal1.LocalId,
			},
			{
				ServiceId: serviceLocal2.ServiceId,
				LocalId:   serviceLocal2.LocalId,
			},
		},
	}

	requestBody, _ := json.Marshal(bulkRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/service-local/bulk/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestBulkDeleteServiceLocalsInvalidRequest(t *testing.T) {
	/*
		GIVEN: Invalid request body
		WHEN:  DELETE /service-local/bulk/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/service-local/bulk/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestBulkDeleteServiceLocalsNonExistent(t *testing.T) {
	/*
		GIVEN: Non-existent service-local associations
		WHEN:  DELETE /service-local/bulk/ is called with non-existent associations
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// Create delete request with non-existent associations
	nonExistentServiceId, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
	nonExistentLocalId, _ := uuid.Parse("00000000-0000-0000-0000-000000000001")
	request := schemas.BulkDeleteServiceLocalRequest{
		ServiceLocals: []*schemas.DeleteServiceLocalRequest{
			{
				ServiceId: nonExistentServiceId,
				LocalId:   nonExistentLocalId,
			},
		},
	}

	requestBody, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/service-local/bulk/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
