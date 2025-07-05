package session_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestCreateSessionSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid session data
		WHEN:  POST /session/ is called with valid data
		THEN:  A HTTP_201_CREATED status should be returned with the created session
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create dependencies
	professional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	local := factories.NewLocalModel(db, factories.LocalModelF{})

	sessionRequest := schemas.CreateSessionRequest{
		Title:          "Yoga Session",
		Date:           time.Now().AddDate(0, 0, 1),    // Tomorrow
		StartTime:      time.Now().Add(10 * time.Hour), // 10 hours from now
		EndTime:        time.Now().Add(11 * time.Hour), // 11 hours from now
		Capacity:       20,
		SessionLink:    func() *string { s := "https://zoom.us/session"; return &s }(),
		ProfessionalId: professional.Id,
		LocalId:        &local.Id,
	}

	requestBody, _ := json.Marshal(sessionRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/session/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response schemas.Session
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains the correct data
	assert.NotEmpty(t, response.Id)
	assert.Equal(t, sessionRequest.Title, response.Title)
	assert.Equal(t, sessionRequest.Capacity, response.Capacity)
	assert.Equal(t, sessionRequest.ProfessionalId, response.ProfessionalId)
	assert.Equal(t, sessionRequest.LocalId, response.LocalId)
	assert.Equal(t, sessionRequest.SessionLink, response.SessionLink)
}

func TestCreateSessionInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  POST /session/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/session/", strings.NewReader(`{"invalid": json`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCreateSessionMissingRequiredFields(t *testing.T) {
	/*
		GIVEN: A request body missing required fields
		WHEN:  POST /session/ is called with incomplete data
		THEN:  A HTTP_400_BAD_REQUEST status should be returned (validation fails)
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// Missing required fields (professional_id is zero value, invalid time range)
	sessionRequest := schemas.CreateSessionRequest{
		Title:    "Incomplete Session",
		Capacity: 10,
	}

	requestBody, _ := json.Marshal(sessionRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/session/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateSessionNonExistentProfessional(t *testing.T) {
	/*
		GIVEN: A session with non-existent professional ID
		WHEN:  POST /session/ is called with non-existent professional
		THEN:  A HTTP_404_NOT_FOUND status should be returned (foreign key validation fails)
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	sessionRequest := schemas.CreateSessionRequest{
		Title:          "Session with Non-existent Professional",
		Date:           time.Now().AddDate(0, 0, 1),
		StartTime:      time.Now().Add(8 * time.Hour),
		EndTime:        time.Now().Add(9 * time.Hour),
		Capacity:       15,
		ProfessionalId: uuid.New(), // Non-existent professional
	}

	requestBody, _ := json.Marshal(sessionRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/session/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCreateSessionTimeConflict(t *testing.T) {
	/*
		GIVEN: A session already exists at the same time and professional
		WHEN:  POST /session/ is called with conflicting time
		THEN:  A HTTP_409_CONFLICT status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create dependencies
	professional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	local := factories.NewLocalModel(db, factories.LocalModelF{})

	// Use the same date for both sessions
	testDate := time.Now().AddDate(0, 0, 1)
	existingStartTime := testDate.Add(10 * time.Hour)
	existingEndTime := existingStartTime.Add(1 * time.Hour)

	// Create existing session
	factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId: &professional.Id,
		LocalId:        &local.Id,
		Date:           &testDate,
		StartTime:      &existingStartTime,
		EndTime:        &existingEndTime,
	})

	// Try to create conflicting session
	sessionRequest := schemas.CreateSessionRequest{
		Title:          "Conflicting Session",
		Date:           testDate,
		StartTime:      existingStartTime.Add(30 * time.Minute), // Overlaps with existing
		EndTime:        existingStartTime.Add(90 * time.Minute),
		Capacity:       20,
		ProfessionalId: professional.Id,
		LocalId:        &local.Id,
	}

	requestBody, _ := json.Marshal(sessionRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/session/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusConflict, rec.Code)
}
