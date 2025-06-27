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

func TestUpdateSessionSuccessfully(t *testing.T) {
	/*
		GIVEN: A session record exists in the database
		WHEN:  UpdateSession is called with valid data
		THEN:  The session record should be updated successfully
	*/
	// GIVEN
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Create dependencies
	testProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testLocal := factories.NewLocalModel(db, factories.LocalModelF{})

	// Create session
	sessionDate := time.Now().Add(24 * time.Hour)
	startTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	endTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)
	title := "Original Session"
	capacity := 10
	testSession := factories.NewSessionModel(db, factories.SessionModelF{
		Title:          &title,
		Date:           &sessionDate,
		StartTime:      &startTime,
		EndTime:        &endTime,
		Capacity:       &capacity,
		ProfessionalId: &testProfessional.Id,
		LocalId:        &testLocal.Id,
	})

	// Prepare update data
	newTitle := "Updated Session"
	newCapacity := 15
	updateRequest := schemas.UpdateSessionRequest{
		Title:    &newTitle,
		Capacity: &newCapacity,
	}

	// WHEN
	result, err := controller.UpdateSession(testSession.Id, updateRequest, "test_admin")

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testSession.Id, result.Id)
	assert.Equal(t, newTitle, result.Title)
	assert.Equal(t, newCapacity, result.Capacity)
	assert.Equal(t, testSession.ProfessionalId, result.ProfessionalId) // Should remain unchanged
}

func TestUpdateSessionNotFound(t *testing.T) {
	/*
		GIVEN: No session record exists with the given ID
		WHEN:  UpdateSession is called with non-existent ID
		THEN:  An error should be returned
	*/
	// GIVEN
	controller, _, _ := controllerTest.NewSessionControllerTestWrapper(t)
	nonExistentId := uuid.New()
	updateRequest := schemas.UpdateSessionRequest{}

	// WHEN
	result, err := controller.UpdateSession(nonExistentId, updateRequest, "test_admin")

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Message, "Session not updated")
}

func TestUpdateSessionWithPartialData(t *testing.T) {
	/*
		GIVEN: A session record exists in the database
		WHEN:  UpdateSession is called with only some fields
		THEN:  Only the specified fields should be updated
	*/
	// GIVEN
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Create dependencies
	testProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	// Create session
	sessionDate := time.Now().Add(24 * time.Hour)
	startTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	endTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)
	title := "Original Session"
	capacity := 10
	testSession := factories.NewSessionModel(db, factories.SessionModelF{
		Title:          &title,
		Date:           &sessionDate,
		StartTime:      &startTime,
		EndTime:        &endTime,
		Capacity:       &capacity,
		ProfessionalId: &testProfessional.Id,
	})

	// Prepare update data - only update capacity
	newCapacity := 20
	updateRequest := schemas.UpdateSessionRequest{
		Capacity: &newCapacity,
	}

	// WHEN
	result, err := controller.UpdateSession(testSession.Id, updateRequest, "test_admin")

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testSession.Id, result.Id)
	assert.Equal(t, newCapacity, result.Capacity)
	assert.Equal(t, testSession.Title, result.Title) // Should remain unchanged
}

func TestUpdateSessionEmptyUpdatedBy(t *testing.T) {
	/*
		GIVEN: A session record exists in the database
		WHEN:  UpdateSession is called with empty updatedBy
		THEN:  An error should be returned
	*/
	// GIVEN
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Create dependencies
	testProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	// Create session
	sessionDate := time.Now().Add(24 * time.Hour)
	startTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	endTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)
	title := "Original Session"
	capacity := 10
	testSession := factories.NewSessionModel(db, factories.SessionModelF{
		Title:          &title,
		Date:           &sessionDate,
		StartTime:      &startTime,
		EndTime:        &endTime,
		Capacity:       &capacity,
		ProfessionalId: &testProfessional.Id,
	})

	updateRequest := schemas.UpdateSessionRequest{}

	// WHEN
	result, err := controller.UpdateSession(testSession.Id, updateRequest, "")

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "Invalid updated by value", err.Message)
}
