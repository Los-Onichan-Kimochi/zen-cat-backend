package session_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestCheckSessionConflictsNoConflict(t *testing.T) {
	/*
		GIVEN: No existing sessions conflict with the requested time
		WHEN:  POST /session/check-conflicts/ is called with non-conflicting time
		THEN:  A HTTP_200_OK status should be returned with no conflicts
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create dependencies
	professional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	local := factories.NewLocalModel(db, factories.LocalModelF{})

	// Create a conflict check request for a time with no existing sessions
	request := schemas.CheckConflictRequest{
		Date:           time.Now().AddDate(0, 0, 1),
		StartTime:      time.Now().Add(8 * time.Hour),
		EndTime:        time.Now().Add(9 * time.Hour),
		ProfessionalId: professional.Id,
		LocalId:        &local.Id,
	}

	requestBody, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/session/check-conflicts/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.ConflictResult
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	assert.False(t, response.HasConflict)
	assert.Empty(t, response.ProfessionalConflicts)
	assert.Empty(t, response.LocalConflicts)
}

func TestCheckSessionConflictsWithProfessionalConflict(t *testing.T) {
	/*
		GIVEN: An existing session conflicts with the professional's time
		WHEN:  POST /session/check-conflicts/ is called with conflicting professional time
		THEN:  A HTTP_200_OK status should be returned with professional conflicts
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create dependencies
	professional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	local1 := factories.NewLocalModel(db, factories.LocalModelF{})
	local2 := factories.NewLocalModel(db, factories.LocalModelF{})

	// Use the same date for both existing session and conflict check
	testDate := time.Now().AddDate(0, 0, 1)
	existingStartTime := testDate.Add(8 * time.Hour)
	existingEndTime := existingStartTime.Add(1 * time.Hour)

	// Create existing session with the professional
	factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId: &professional.Id,
		LocalId:        &local1.Id,
		Date:           &testDate,
		StartTime:      &existingStartTime,
		EndTime:        &existingEndTime,
	})

	// Create a conflict check request for overlapping time with same professional, different local
	request := schemas.CheckConflictRequest{
		Date:           testDate,                                // Same date as existing session
		StartTime:      existingStartTime.Add(30 * time.Minute), // Overlaps
		EndTime:        existingStartTime.Add(90 * time.Minute),
		ProfessionalId: professional.Id,
		LocalId:        &local2.Id, // Different local
	}

	requestBody, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/session/check-conflicts/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.ConflictResult
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	assert.True(t, response.HasConflict)
	assert.NotEmpty(t, response.ProfessionalConflicts)
	assert.Empty(t, response.LocalConflicts) // Different local, so no local conflict
}

func TestCheckSessionConflictsInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  POST /session/check-conflicts/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/session/check-conflicts/", strings.NewReader(`{"invalid": json`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
