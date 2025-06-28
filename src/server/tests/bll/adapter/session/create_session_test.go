package session_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	adapterTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/adapter"
)

func TestCreateSessionSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid session data with existing professional and local
		WHEN:  CreatePostgresqlSession is called
		THEN:  A new session is created and returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewSessionAdapterTestWrapper(t)

	professional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	local := factories.NewLocalModel(db, factories.LocalModelF{})

	title := "Test Session"
	date := time.Now().AddDate(0, 0, 1) // Tomorrow
	startTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	endTime := time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)
	capacity := 20
	sessionLink := "https://meet.google.com/test"
	updatedBy := "test-admin"

	// WHEN
	session, err := adapter.CreatePostgresqlSession(
		title,
		date,
		startTime,
		endTime,
		capacity,
		&sessionLink,
		professional.Id,
		&local.Id,
		updatedBy,
	)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, session)
	assert.NotEmpty(t, session.Id)
	assert.Equal(t, title, session.Title)
	assert.Equal(t, date.Format("2006-01-02"), session.Date.Format("2006-01-02"))
	assert.Equal(t, capacity, session.Capacity)
	assert.Equal(t, &sessionLink, session.SessionLink)
	assert.Equal(t, professional.Id, session.ProfessionalId)
	assert.Equal(t, &local.Id, session.LocalId)
	assert.Equal(t, 0, session.RegisteredCount)
}

func TestCreateSessionWithoutLocal(t *testing.T) {
	/*
		GIVEN: Valid session data with professional but no local
		WHEN:  CreatePostgresqlSession is called
		THEN:  A new session is created without local
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewSessionAdapterTestWrapper(t)

	professional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	title := "Online Session"
	date := time.Now().AddDate(0, 0, 1)
	startTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	endTime := time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)
	capacity := 50
	sessionLink := "https://meet.google.com/online"
	updatedBy := "test-admin"

	// WHEN
	session, err := adapter.CreatePostgresqlSession(
		title,
		date,
		startTime,
		endTime,
		capacity,
		&sessionLink,
		professional.Id,
		nil, // No local
		updatedBy,
	)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, title, session.Title)
	assert.Equal(t, professional.Id, session.ProfessionalId)
	assert.Nil(t, session.LocalId)
}

func TestCreateSessionWithEmptyUpdatedBy(t *testing.T) {
	/*
		GIVEN: Valid session data but empty updatedBy
		WHEN:  CreatePostgresqlSession is called
		THEN:  An error is returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewSessionAdapterTestWrapper(t)

	professional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	title := "Test Session"
	date := time.Now().AddDate(0, 0, 1)
	startTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	endTime := time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)
	capacity := 20
	updatedBy := ""

	// WHEN
	session, err := adapter.CreatePostgresqlSession(
		title,
		date,
		startTime,
		endTime,
		capacity,
		nil,
		professional.Id,
		nil,
		updatedBy,
	)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, session)
	assert.Equal(t, errors.BadRequestError.InvalidUpdatedByValue, *err)
}
