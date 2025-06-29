package session_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	adapterTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/adapter"
)

func TestUpdateSessionSuccessfully(t *testing.T) {
	/*
		GIVEN: A session exists in the database
		WHEN:  UpdatePostgresqlSession is called with new data
		THEN:  The session is updated and returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewSessionAdapterTestWrapper(t)

	session := factories.NewSessionModel(db, factories.SessionModelF{})

	newTitle := "Updated Session Title"
	newCapacity := 30
	updatedBy := "test-admin"

	// WHEN
	updatedSession, err := adapter.UpdatePostgresqlSession(
		session.Id,
		&newTitle,
		nil, // Don't update date
		nil, // Don't update start time
		nil, // Don't update end time
		nil, // Don't update state
		nil, // Don't update registered count
		&newCapacity,
		nil, // Don't update session link
		nil, // Don't update professional
		nil, // Don't update local
		nil, // Don't update community service
		updatedBy,
	)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, updatedSession)
	assert.Equal(t, session.Id, updatedSession.Id)
	assert.Equal(t, newTitle, updatedSession.Title)
	assert.Equal(t, newCapacity, updatedSession.Capacity)
}

func TestUpdateSessionWithEmptyUpdatedBy(t *testing.T) {
	/*
		GIVEN: A session exists in the database
		WHEN:  UpdatePostgresqlSession is called with empty updatedBy
		THEN:  An error is returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewSessionAdapterTestWrapper(t)

	session := factories.NewSessionModel(db, factories.SessionModelF{})

	newTitle := "Updated Session Title"
	updatedBy := ""

	// WHEN
	updatedSession, err := adapter.UpdatePostgresqlSession(
		session.Id,
		&newTitle,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		updatedBy,
	)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, updatedSession)
	assert.Equal(t, errors.BadRequestError.InvalidUpdatedByValue, *err)
}

func TestUpdateSessionWithCompleteData(t *testing.T) {
	/*
		GIVEN: A session exists in the database
		WHEN:  UpdatePostgresqlSession is called with all fields
		THEN:  The session is completely updated
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewSessionAdapterTestWrapper(t)

	session := factories.NewSessionModel(db, factories.SessionModelF{})
	newProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	newLocal := factories.NewLocalModel(db, factories.LocalModelF{})
	newCommunityService := factories.NewCommunityServiceModel(db, factories.CommunityServiceModelF{})

	newTitle := "Completely Updated Session"
	newDate := time.Now().AddDate(0, 0, 2)
	newStartTime := time.Date(2024, 1, 1, 14, 0, 0, 0, time.UTC)
	newEndTime := time.Date(2024, 1, 1, 15, 30, 0, 0, time.UTC)
	newState := "ACTIVE"
	newRegisteredCount := 5
	newCapacity := 25
	newSessionLink := "https://updated-link.com"
	updatedBy := "test-admin"

	// WHEN
	updatedSession, err := adapter.UpdatePostgresqlSession(
		session.Id,
		&newTitle,
		&newDate,
		&newStartTime,
		&newEndTime,
		&newState,
		&newRegisteredCount,
		&newCapacity,
		&newSessionLink,
		&newProfessional.Id,
		&newLocal.Id,
		&newCommunityService.Id,
		updatedBy,
	)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, updatedSession)
	assert.Equal(t, session.Id, updatedSession.Id)
	assert.Equal(t, newTitle, updatedSession.Title)
	assert.Equal(t, newCapacity, updatedSession.Capacity)
	assert.Equal(t, newRegisteredCount, updatedSession.RegisteredCount)
	assert.Equal(t, &newSessionLink, updatedSession.SessionLink)
	assert.Equal(t, newProfessional.Id, updatedSession.ProfessionalId)
	assert.Equal(t, &newLocal.Id, updatedSession.LocalId)
	assert.Equal(t, &newCommunityService.Id, updatedSession.CommunityServiceId)
}

func TestDeleteSessionSuccessfully(t *testing.T) {
	/*
		GIVEN: A session exists in the database
		WHEN:  DeletePostgresqlSession is called
		THEN:  The session is soft deleted
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewSessionAdapterTestWrapper(t)

	session := factories.NewSessionModel(db, factories.SessionModelF{})

	// WHEN
	err := adapter.DeletePostgresqlSession(session.Id)

	// THEN
	assert.Nil(t, err)

	// Verify session is deleted by trying to get it
	deletedSession, getErr := adapter.GetPostgresqlSession(session.Id)
	assert.NotNil(t, getErr)
	assert.Nil(t, deletedSession)
	assert.Equal(t, errors.ObjectNotFoundError.SessionNotFound, *getErr)
}

func TestDeleteSessionNotFound(t *testing.T) {
	/*
		GIVEN: No session exists with the given ID
		WHEN:  DeletePostgresqlSession is called
		THEN:  An error is returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewSessionAdapterTestWrapper(t)

	nonExistentId := uuid.New()

	// WHEN
	err := adapter.DeletePostgresqlSession(nonExistentId)

	// THEN
	assert.NotNil(t, err)
	assert.Equal(t, errors.BadRequestError.SessionNotSoftDeleted, *err)
}
