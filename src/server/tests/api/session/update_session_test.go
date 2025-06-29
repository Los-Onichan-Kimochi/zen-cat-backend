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
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestUpdateSessionSuccessfully(t *testing.T) {
	/*
		GIVEN: A session exists in the database
		WHEN:  PATCH /session/{sessionId}/ is called with valid update data
		THEN:  A HTTP_200_OK status should be returned with the updated session
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create dependencies
	professional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	local := factories.NewLocalModel(db, factories.LocalModelF{})

	// Create a test session using factory
	session := factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId: &professional.Id,
		LocalId:        &local.Id,
	})

	// Prepare update data
	newTitle := "Updated Yoga Session"
	newCapacity := 25
	updateRequest := schemas.UpdateSessionRequest{
		Title:    &newTitle,
		Capacity: &newCapacity,
	}

	requestBody, _ := json.Marshal(updateRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/session/"+session.Id.String()+"/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Session
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains the updated data
	assert.Equal(t, session.Id, response.Id)
	assert.Equal(t, newTitle, response.Title)
	assert.Equal(t, newCapacity, response.Capacity)
	// Other fields should remain unchanged
	assert.Equal(t, session.ProfessionalId, response.ProfessionalId)
}

func TestUpdateSessionNotFound(t *testing.T) {
	/*
		GIVEN: A non-existent session ID
		WHEN:  PATCH /session/{id}/ is called with non-existent ID
		THEN:  A HTTP_400_BAD_REQUEST status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentId := uuid.New()

	updateRequest := schemas.UpdateSessionRequest{
		Title: func() *string { s := "Updated Session"; return &s }(),
	}

	requestBody, _ := json.Marshal(updateRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/session/"+nonExistentId.String()+"/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateSessionInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID format is provided
		WHEN:  PATCH /session/{sessionId}/ is called with an invalid UUID
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidSessionId := "invalid-uuid"

	updateRequest := schemas.UpdateSessionRequest{
		Title: func() *string { s := "Updated Session"; return &s }(),
	}

	requestBody, _ := json.Marshal(updateRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/session/"+invalidSessionId+"/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestUpdateSessionInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  PATCH /session/{sessionId}/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	professional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	local := factories.NewLocalModel(db, factories.LocalModelF{})
	session := factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId: &professional.Id,
		LocalId:        &local.Id,
	})

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/session/"+session.Id.String()+"/", strings.NewReader(`{"invalid": json`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
