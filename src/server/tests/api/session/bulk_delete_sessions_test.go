package session_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestBulkDeleteSessionsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple sessions exist in the database
		WHEN:  DELETE /session/bulk-delete/ is called with valid session IDs
		THEN:  A HTTP_204_NO_CONTENT status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create dependencies
	professional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	local := factories.NewLocalModel(db, factories.LocalModelF{})

	// Create test sessions using factory
	numSessions := 3
	testSessions := make([]*model.Session, numSessions)
	for i := 0; i < numSessions; i++ {
		testSessions[i] = factories.NewSessionModel(db, factories.SessionModelF{
			ProfessionalId: &professional.Id,
			LocalId:        &local.Id,
		})
	}

	// Extract session IDs
	sessionIds := make([]uuid.UUID, numSessions)
	for i, session := range testSessions {
		sessionIds[i] = session.Id
	}

	bulkDeleteRequest := schemas.BulkDeleteSessionRequest{
		Sessions: sessionIds,
	}

	requestBody, _ := json.Marshal(bulkDeleteRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/session/bulk-delete/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestBulkDeleteSessionsEmptyList(t *testing.T) {
	/*
		GIVEN: A bulk session deletion request with an empty list
		WHEN:  DELETE /session/bulk-delete/ is called with an empty sessions list
		THEN:  A HTTP_204_NO_CONTENT status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	bulkDeleteRequest := schemas.BulkDeleteSessionRequest{
		Sessions: []uuid.UUID{},
	}

	requestBody, _ := json.Marshal(bulkDeleteRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/session/bulk-delete/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestBulkDeleteSessionsInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  DELETE /session/bulk-delete/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/session/bulk-delete/", strings.NewReader(`{"invalid": json`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestBulkDeleteSessionsNonExistentIds(t *testing.T) {
	/*
		GIVEN: Non-existent session IDs
		WHEN:  DELETE /session/bulk-delete/ is called with non-existent IDs
		THEN:  A HTTP_400_BAD_REQUEST status should be returned (API returns error for non-existent records)
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	nonExistentIds := []uuid.UUID{
		uuid.New(),
		uuid.New(),
	}

	bulkDeleteRequest := schemas.BulkDeleteSessionRequest{
		Sessions: nonExistentIds,
	}

	requestBody, _ := json.Marshal(bulkDeleteRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/session/bulk-delete/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
