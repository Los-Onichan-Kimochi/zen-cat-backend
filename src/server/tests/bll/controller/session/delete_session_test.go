package session_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestDeleteSessionSuccessfully(t *testing.T) {
	/*
		GIVEN: A session record exists in the database
		WHEN:  DeleteSession is called with valid ID
		THEN:  The session record should be deleted successfully
	*/
	// GIVEN
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Create dependencies
	testProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	// Create session
	sessionDate := time.Now().Add(24 * time.Hour)
	startTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	endTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)
	title := "Test Session"
	capacity := 10
	testSession := factories.NewSessionModel(db, factories.SessionModelF{
		Title:          &title,
		Date:           &sessionDate,
		StartTime:      &startTime,
		EndTime:        &endTime,
		Capacity:       &capacity,
		ProfessionalId: &testProfessional.Id,
	})

	// WHEN
	err := controller.DeleteSession(testSession.Id)

	// THEN
	assert.Nil(t, err)

	// Verify the session was deleted
	var deletedSession model.Session
	dbErr := db.Where("id = ?", testSession.Id).First(&deletedSession).Error
	assert.Error(t, dbErr) // Should return error because record doesn't exist
}

func TestDeleteSessionNotFound(t *testing.T) {
	/*
		GIVEN: No session record exists with the given ID
		WHEN:  DeleteSession is called with non-existent ID
		THEN:  An error should be returned
	*/
	// GIVEN
	controller, _, _ := controllerTest.NewSessionControllerTestWrapper(t)
	nonExistentId := uuid.New()

	// WHEN
	err := controller.DeleteSession(nonExistentId)

	// THEN
	assert.NotNil(t, err)
	assert.Contains(t, err.Message, "Session not soft deleted")
}

func TestDeleteSessionWithNilId(t *testing.T) {
	/*
		GIVEN: A nil UUID
		WHEN:  DeleteSession is called with nil UUID
		THEN:  An error should be returned
	*/
	// GIVEN
	controller, _, _ := controllerTest.NewSessionControllerTestWrapper(t)
	nilId := uuid.Nil

	// WHEN
	err := controller.DeleteSession(nilId)

	// THEN
	assert.NotNil(t, err)
}
