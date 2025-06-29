package session_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestDeleteSessionSuccessfully(t *testing.T) {
	/*
		GIVEN: A session exists in the database
		WHEN:  DELETE /session/{sessionId}/ is called with a valid session ID
		THEN:  A HTTP_204_NO_CONTENT status should be returned
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
	req := httptest.NewRequest(http.MethodDelete, "/session/"+session.Id.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestDeleteSessionNotFound(t *testing.T) {
	/*
		GIVEN: A non-existent session ID
		WHEN:  DELETE /session/{id}/ is called with non-existent ID
		THEN:  A HTTP_400_BAD_REQUEST status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentId := uuid.New()

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/session/"+nonExistentId.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteSessionInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID format is provided
		WHEN:  DELETE /session/{sessionId}/ is called with an invalid UUID
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidSessionId := "invalid-uuid"

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/session/"+invalidSessionId+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
