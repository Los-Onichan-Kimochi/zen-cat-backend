package session_test

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

func TestGetSessionSuccessfully(t *testing.T) {
	/*
		GIVEN: A session exists in the database
		WHEN:  GET /session/{sessionId}/ is called with a valid session ID
		THEN:  A HTTP_200_OK status should be returned with the session data
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

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/session/"+session.Id.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Session
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains the correct session data
	assert.Equal(t, session.Id, response.Id)
	assert.Equal(t, session.Title, response.Title)
	assert.Equal(t, session.Capacity, response.Capacity)
	assert.Equal(t, session.ProfessionalId, response.ProfessionalId)
	assert.Equal(t, session.LocalId, response.LocalId)
	assert.Equal(t, string(session.State), response.State)
}

func TestGetSessionNotFound(t *testing.T) {
	/*
		GIVEN: No session exists with the provided ID
		WHEN:  GET /session/{sessionId}/ is called with a non-existent session ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentSessionId := uuid.New()

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/session/"+nonExistentSessionId.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetSessionInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID format is provided
		WHEN:  GET /session/{sessionId}/ is called with an invalid UUID
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidSessionId := "invalid-uuid"

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/session/"+invalidSessionId+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
