package session_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestGetSessionSuccessfully(t *testing.T) {
	/*
		GIVEN: A session record exists in the database
		WHEN:  GetSession is called with valid ID
		THEN:  The session record should be returned
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
	title := "Test Session"
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

	// WHEN
	result, err := controller.GetSession(testSession.Id)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testSession.Id, result.Id)
	assert.Equal(t, testSession.Title, result.Title)
	assert.Equal(t, testSession.ProfessionalId, result.ProfessionalId)
	assert.Equal(t, testSession.LocalId, result.LocalId)
}

func TestGetSessionNotFound(t *testing.T) {
	/*
		GIVEN: No session record exists with the given ID
		WHEN:  GetSession is called with non-existent ID
		THEN:  An error should be returned
	*/
	// GIVEN
	controller, _, _ := controllerTest.NewSessionControllerTestWrapper(t)
	nonExistentId := uuid.New()

	// WHEN
	result, err := controller.GetSession(nonExistentId)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Message, "not found")
}

func TestGetSessionWithNilId(t *testing.T) {
	/*
		GIVEN: A nil UUID
		WHEN:  GetSession is called with nil UUID
		THEN:  An error should be returned
	*/
	// GIVEN
	controller, _, _ := controllerTest.NewSessionControllerTestWrapper(t)
	nilId := uuid.Nil

	// WHEN
	result, err := controller.GetSession(nilId)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, result)
}
