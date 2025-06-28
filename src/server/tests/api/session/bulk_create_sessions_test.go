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

func TestBulkCreateSessionsSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid sessions data
		WHEN:  POST /session/bulk/ is called with valid data
		THEN:  A HTTP_201_CREATED status should be returned with the created sessions
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create dependencies
	professional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	local := factories.NewLocalModel(db, factories.LocalModelF{})

	request := schemas.BatchCreateSessionRequest{
		Sessions: []*schemas.CreateSessionRequest{
			{
				Title:          "Morning Yoga",
				Date:           time.Now().AddDate(0, 0, 1),
				StartTime:      time.Now().Add(8 * time.Hour),
				EndTime:        time.Now().Add(9 * time.Hour),
				Capacity:       20,
				SessionLink:    func() *string { s := "https://zoom.us/session1"; return &s }(),
				ProfessionalId: professional.Id,
				LocalId:        &local.Id,
			},
			{
				Title:          "Evening Yoga",
				Date:           time.Now().AddDate(0, 0, 1),
				StartTime:      time.Now().Add(18 * time.Hour),
				EndTime:        time.Now().Add(19 * time.Hour),
				Capacity:       15,
				SessionLink:    func() *string { s := "https://zoom.us/session2"; return &s }(),
				ProfessionalId: professional.Id,
				LocalId:        &local.Id,
			},
		},
	}

	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/session/bulk/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response schemas.Sessions
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response.Sessions, 2)

	// Verify the created sessions
	for i, session := range response.Sessions {
		assert.NotEmpty(t, session.Id)
		assert.Equal(t, request.Sessions[i].Title, session.Title)
		assert.Equal(t, request.Sessions[i].Capacity, session.Capacity)
		assert.Equal(t, request.Sessions[i].ProfessionalId, session.ProfessionalId)
	}
}

func TestBulkCreateSessionsEmptyList(t *testing.T) {
	/*
		GIVEN: An empty sessions list
		WHEN:  POST /session/bulk/ is called with empty list
		THEN:  A HTTP_201_CREATED status should be returned with empty response
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	request := schemas.BatchCreateSessionRequest{
		Sessions: []*schemas.CreateSessionRequest{},
	}

	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/session/bulk/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response schemas.Sessions
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Empty(t, response.Sessions)
}

func TestBulkCreateSessionsInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  POST /session/bulk/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/session/bulk/", strings.NewReader(`{"invalid": json`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestBulkCreateSessionsPartialFailure(t *testing.T) {
	/*
		GIVEN: A mix of valid and invalid session data (non-existent professional)
		WHEN:  POST /session/bulk/ is called with some invalid data
		THEN:  A HTTP_404_NOT_FOUND status should be returned (foreign key validation fails)
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create valid dependencies
	professional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	sessionLink1 := "https://example.com/session1"
	sessionLink2 := "https://example.com/session2"

	request := schemas.BatchCreateSessionRequest{
		Sessions: []*schemas.CreateSessionRequest{
			{
				Title:          "Valid Session",
				Date:           time.Now().AddDate(0, 0, 1),
				StartTime:      time.Now().Add(8 * time.Hour),
				EndTime:        time.Now().Add(9 * time.Hour),
				Capacity:       10,
				SessionLink:    &sessionLink1,
				ProfessionalId: professional.Id,
			},
			{
				// Invalid: non-existent professional
				Title:          "Invalid Session",
				Date:           time.Now().AddDate(0, 0, 2),
				StartTime:      time.Now().Add(10 * time.Hour),
				EndTime:        time.Now().Add(11 * time.Hour),
				Capacity:       15,
				SessionLink:    &sessionLink2,
				ProfessionalId: uuid.New(), // Non-existent professional
			},
		},
	}

	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/session/bulk/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
