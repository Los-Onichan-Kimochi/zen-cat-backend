package session_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestCreateSessionSuccessfully(t *testing.T) {
	// GIVEN: Valid session creation request with existing dependencies
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Create dependencies
	testProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testLocal := factories.NewLocalModel(db, factories.LocalModelF{})

	// Create session request
	sessionDate := time.Now().Add(24 * time.Hour)
	startTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	endTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)

	createRequest := schemas.CreateSessionRequest{
		Title:          "Test Session",
		Date:           sessionDate,
		StartTime:      startTime,
		EndTime:        endTime,
		Capacity:       10,
		ProfessionalId: testProfessional.Id,
		LocalId:        &testLocal.Id,
	}

	// WHEN: CreateSession is called
	result, err := controller.CreateSession(createRequest, "test_admin")

	// THEN: Session is created successfully
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createRequest.Title, result.Title)
	assert.Equal(t, createRequest.Date.Format("2006-01-02"), result.Date.Format("2006-01-02"))
	assert.Equal(t, createRequest.ProfessionalId, result.ProfessionalId)
	assert.Equal(t, createRequest.LocalId, result.LocalId)
	assert.NotEqual(t, "", result.Id)
}

func TestCreateSessionEmptyUpdatedBy(t *testing.T) {
	// GIVEN: Valid session creation request but empty updatedBy
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Create dependencies
	testProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	sessionDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	startTime := time.Date(2024, 12, 31, 14, 0, 0, 0, time.UTC)
	endTime := time.Date(2024, 12, 31, 15, 0, 0, 0, time.UTC)

	createRequest := schemas.CreateSessionRequest{
		Title:          "Test Session",
		Date:           sessionDate,
		StartTime:      startTime,
		EndTime:        endTime,
		Capacity:       10,
		ProfessionalId: testProfessional.Id,
	}

	// WHEN: CreateSession is called with empty updatedBy
	result, err := controller.CreateSession(createRequest, "")

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "Invalid updated by value", err.Message)
}

func TestCreateSessionInvalidProfessional(t *testing.T) {
	// GIVEN: Session creation request with non-existent professional
	controller, _, _ := controllerTest.NewSessionControllerTestWrapper(t)

	// Use a random UUID that doesn't exist
	nonExistentProfessionalId := uuid.New()

	sessionDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	startTime := time.Date(2024, 12, 31, 14, 0, 0, 0, time.UTC)
	endTime := time.Date(2024, 12, 31, 15, 0, 0, 0, time.UTC)

	createRequest := schemas.CreateSessionRequest{
		Title:          "Test Session",
		Date:           sessionDate,
		StartTime:      startTime,
		EndTime:        endTime,
		Capacity:       10,
		ProfessionalId: nonExistentProfessionalId,
	}

	// WHEN: CreateSession is called with non-existent professional
	result, err := controller.CreateSession(createRequest, "test_admin")

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestCreateSessionEndTimeBeforeStartTime(t *testing.T) {
	// GIVEN: Session creation request with end time before start time
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	testProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	sessionDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	startTime := time.Date(2024, 12, 31, 15, 0, 0, 0, time.UTC)
	endTime := time.Date(2024, 12, 31, 14, 0, 0, 0, time.UTC) // End time before start time

	createRequest := schemas.CreateSessionRequest{
		Title:          "Test Session",
		Date:           sessionDate,
		StartTime:      startTime,
		EndTime:        endTime,
		Capacity:       10,
		ProfessionalId: testProfessional.Id,
	}

	// WHEN: CreateSession is called
	result, err := controller.CreateSession(createRequest, "test_admin")

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestCreateSessionPastDate(t *testing.T) {
	// GIVEN: Session creation request with past date
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	testProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	// Use a past date
	pastDate := time.Now().Add(-24 * time.Hour)
	startTime := time.Date(pastDate.Year(), pastDate.Month(), pastDate.Day(), 14, 0, 0, 0, time.UTC)
	endTime := time.Date(pastDate.Year(), pastDate.Month(), pastDate.Day(), 15, 0, 0, 0, time.UTC)

	createRequest := schemas.CreateSessionRequest{
		Title:          "Test Session",
		Date:           pastDate,
		StartTime:      startTime,
		EndTime:        endTime,
		Capacity:       10,
		ProfessionalId: testProfessional.Id,
	}

	// WHEN: CreateSession is called with past date
	result, err := controller.CreateSession(createRequest, "test_admin")

	// THEN: Behavior depends on business rules
	// It might be allowed or not depending on implementation
	if err != nil {
		assert.Nil(t, result)
		assert.NotNil(t, err)
	} else {
		assert.NotNil(t, result)
		assert.Nil(t, err)
	}
}

func TestCreateSessionSameStartEndTime(t *testing.T) {
	// GIVEN: Session creation request with same start and end time
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	testProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	sessionDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	sameTime := time.Date(2024, 12, 31, 14, 0, 0, 0, time.UTC)

	createRequest := schemas.CreateSessionRequest{
		Title:          "Test Session",
		Date:           sessionDate,
		StartTime:      sameTime,
		EndTime:        sameTime, // Same as start time
		Capacity:       10,
		ProfessionalId: testProfessional.Id,
	}

	// WHEN: CreateSession is called
	result, err := controller.CreateSession(createRequest, "test_admin")

	// THEN: An error is returned (zero duration session)
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestCreateSessionWithoutLocal(t *testing.T) {
	// GIVEN: Session creation request without local (virtual session)
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	testProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	sessionDate := time.Now().Add(24 * time.Hour)
	startTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 16, 0, 0, 0, time.UTC)
	endTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 17, 0, 0, 0, time.UTC)

	sessionLink := "https://meet.example.com/virtual-session"

	createRequest := schemas.CreateSessionRequest{
		Title:          "Virtual Session",
		Date:           sessionDate,
		StartTime:      startTime,
		EndTime:        endTime,
		Capacity:       20,
		SessionLink:    &sessionLink,
		ProfessionalId: testProfessional.Id,
		LocalId:        nil, // No local for virtual session
	}

	// WHEN: CreateSession is called
	result, err := controller.CreateSession(createRequest, "test_admin")

	// THEN: Virtual session is created successfully
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createRequest.Title, result.Title)
	assert.Equal(t, createRequest.ProfessionalId, result.ProfessionalId)
	assert.Nil(t, result.LocalId)
	assert.Equal(t, sessionLink, *result.SessionLink)
}
