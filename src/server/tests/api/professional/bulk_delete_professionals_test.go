package professional_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestBulkDeleteProfessionalsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple professionals exist in the database
		WHEN:  DELETE /professional/bulk-delete/ is called with valid professional IDs
		THEN:  A HTTP_204_NO_CONTENT status should be returned and the professionals should be deleted
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create test professionals using factory
	numProfessionals := 3
	testProfessionals := factories.NewProfessionalModelBatch(db, numProfessionals)

	// Extract professional IDs
	professionalIds := make([]uuid.UUID, numProfessionals)
	for i, professional := range testProfessionals {
		professionalIds[i] = professional.Id
	}

	bulkDeleteRequest := schemas.BulkDeleteProfessionalRequest{
		Professionals: professionalIds,
	}

	requestBody, _ := json.Marshal(bulkDeleteRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/professional/bulk-delete/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestBulkDeleteProfessionalsEmptyList(t *testing.T) {
	/*
		GIVEN: A bulk professional deletion request with an empty list
		WHEN:  DELETE /professional/bulk-delete/ is called with an empty professionals list
		THEN:  A HTTP_204_NO_CONTENT status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	bulkDeleteRequest := schemas.BulkDeleteProfessionalRequest{
		Professionals: []uuid.UUID{},
	}

	requestBody, _ := json.Marshal(bulkDeleteRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/professional/bulk-delete/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestBulkDeleteProfessionalsInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  DELETE /professional/bulk-delete/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/professional/bulk-delete/", strings.NewReader(`{"invalid": json`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestBulkDeleteProfessionalsNonExistentIds(t *testing.T) {
	/*
		GIVEN: Non-existent professional IDs
		WHEN:  DELETE /professional/bulk-delete/ is called with non-existent IDs
		THEN:  A HTTP_400_BAD_REQUEST status should be returned (API returns error for non-existent records)
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	nonExistentIds := []uuid.UUID{
		uuid.New(),
		uuid.New(),
	}

	bulkDeleteRequest := schemas.BulkDeleteProfessionalRequest{
		Professionals: nonExistentIds,
	}

	requestBody, _ := json.Marshal(bulkDeleteRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/professional/bulk-delete/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
