package session_test

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

func TestFetchSessionsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple sessions exist in the database
		WHEN:  GET /session/ is called
		THEN:  A HTTP_200_OK status should be returned with the sessions data
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create test sessions using factory
	numSessions := 3
	testSessions := factories.NewSessionModelBatch(db, numSessions)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/session/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Sessions
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains the correct number of sessions
	assert.GreaterOrEqual(t, len(response.Sessions), numSessions)

	// Verify that our created sessions are in the response
	foundSessions := make(map[string]bool)
	for _, session := range response.Sessions {
		foundSessions[session.Id.String()] = true
	}

	for _, testSession := range testSessions {
		assert.True(t, foundSessions[testSession.Id.String()],
			"Created session ID %s not found in response", testSession.Id.String())
	}
}

func TestFetchSessionsWithProfessionalIdFilter(t *testing.T) {
	/*
		GIVEN: Multiple sessions exist with different professionals
		WHEN:  GET /session/?professionalIds={id} is called with a specific professional ID
		THEN:  A HTTP_200_OK status should be returned with filtered sessions
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create professionals
	professional1 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	professional2 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	local := factories.NewLocalModel(db, factories.LocalModelF{})

	// Create sessions for different professionals
	session1 := factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId: &professional1.Id,
		LocalId:        &local.Id,
	})
	session2 := factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId: &professional2.Id,
		LocalId:        &local.Id,
	})

	// WHEN - Filter by professional1 ID
	req := httptest.NewRequest(http.MethodGet, "/session/?professionalIds="+professional1.Id.String(), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Sessions
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Should contain session1 but not session2
	foundSession1 := false
	foundSession2 := false
	for _, session := range response.Sessions {
		if session.Id == session1.Id {
			foundSession1 = true
		}
		if session.Id == session2.Id {
			foundSession2 = true
		}
	}

	assert.True(t, foundSession1, "Session for professional1 should be found")
	assert.False(t, foundSession2, "Session for professional2 should not be found")
}

func TestFetchSessionsWithLocalIdFilter(t *testing.T) {
	/*
		GIVEN: Multiple sessions exist with different locals
		WHEN:  GET /session/?localIds={id} is called with a specific local ID
		THEN:  A HTTP_200_OK status should be returned with filtered sessions
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create locals and professional
	local1 := factories.NewLocalModel(db, factories.LocalModelF{})
	local2 := factories.NewLocalModel(db, factories.LocalModelF{})
	professional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	// Create sessions for different locals
	session1 := factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId: &professional.Id,
		LocalId:        &local1.Id,
	})
	session2 := factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId: &professional.Id,
		LocalId:        &local2.Id,
	})

	// WHEN - Filter by local1 ID
	req := httptest.NewRequest(http.MethodGet, "/session/?localIds="+local1.Id.String(), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Sessions
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Should contain session1 but not session2
	foundSession1 := false
	foundSession2 := false
	for _, session := range response.Sessions {
		if session.Id == session1.Id {
			foundSession1 = true
		}
		if session.Id == session2.Id {
			foundSession2 = true
		}
	}

	assert.True(t, foundSession1, "Session for local1 should be found")
	assert.False(t, foundSession2, "Session for local2 should not be found")
}

func TestFetchSessionsEmpty(t *testing.T) {
	/*
		GIVEN: No sessions exist in the database
		WHEN:  GET /session/ is called
		THEN:  A HTTP_200_OK status should be returned with an empty sessions list
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/session/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Sessions
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Should return empty array, not null
	assert.NotNil(t, response.Sessions)
}
