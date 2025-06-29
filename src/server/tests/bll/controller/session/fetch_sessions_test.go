package session_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestFetchSessionsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple session records exist in the database
		WHEN:  FetchSessions is called
		THEN:  All session records should be returned
	*/
	// GIVEN
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Create dependencies
	testProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	// Create sessions
	sessionDate := time.Now().Add(24 * time.Hour)
	startTime1 := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	endTime1 := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)
	startTime2 := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 14, 0, 0, 0, time.UTC)
	endTime2 := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 15, 0, 0, 0, time.UTC)
	title1 := "Morning Session"
	title2 := "Afternoon Session"
	capacity := 10

	factories.NewSessionModel(db, factories.SessionModelF{
		Title:          &title1,
		Date:           &sessionDate,
		StartTime:      &startTime1,
		EndTime:        &endTime1,
		Capacity:       &capacity,
		ProfessionalId: &testProfessional.Id,
	})

	factories.NewSessionModel(db, factories.SessionModelF{
		Title:          &title2,
		Date:           &sessionDate,
		StartTime:      &startTime2,
		EndTime:        &endTime2,
		Capacity:       &capacity,
		ProfessionalId: &testProfessional.Id,
	})

	// WHEN
	result, err := controller.FetchSessions([]string{}, []string{}, []string{}, []string{})

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, len(result.Sessions), 2)
}

func TestFetchSessionsEmpty(t *testing.T) {
	/*
		GIVEN: No session records exist in the database
		WHEN:  FetchSessions is called
		THEN:  An empty list should be returned
	*/
	// GIVEN
	controller, _, _ := controllerTest.NewSessionControllerTestWrapper(t)

	// WHEN
	result, err := controller.FetchSessions([]string{}, []string{}, []string{}, []string{})

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result.Sessions))
}

func TestFetchSessionsWithFilters(t *testing.T) {
	/*
		GIVEN: Multiple session records exist with different professionals
		WHEN:  FetchSessions is called with professional filter
		THEN:  Only sessions for that professional should be returned
	*/
	// GIVEN
	controller, _, db := controllerTest.NewSessionControllerTestWrapper(t)

	// Create two professionals
	testProfessional1 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testProfessional2 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	// Create sessions for each professional
	sessionDate := time.Now().Add(24 * time.Hour)
	startTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	endTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)
	title1 := "Session 1"
	title2 := "Session 2"
	capacity := 10

	factories.NewSessionModel(db, factories.SessionModelF{
		Title:          &title1,
		Date:           &sessionDate,
		StartTime:      &startTime,
		EndTime:        &endTime,
		Capacity:       &capacity,
		ProfessionalId: &testProfessional1.Id,
	})

	factories.NewSessionModel(db, factories.SessionModelF{
		Title:          &title2,
		Date:           &sessionDate,
		StartTime:      &startTime,
		EndTime:        &endTime,
		Capacity:       &capacity,
		ProfessionalId: &testProfessional2.Id,
	})

	// WHEN - Filter by first professional
	result, err := controller.FetchSessions([]string{testProfessional1.Id.String()}, []string{}, []string{}, []string{})

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Sessions))
	assert.Equal(t, testProfessional1.Id, result.Sessions[0].ProfessionalId)
}
