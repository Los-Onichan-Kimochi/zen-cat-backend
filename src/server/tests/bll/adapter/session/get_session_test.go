package session_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	adapterTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/adapter"
)

func TestGetSessionSuccessfully(t *testing.T) {
	/*
		GIVEN: A session exists in the database
		WHEN:  GetPostgresqlSession is called with the session ID
		THEN:  The session is returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewSessionAdapterTestWrapper(t)

	session := factories.NewSessionModel(db, factories.SessionModelF{})

	// WHEN
	result, err := adapter.GetPostgresqlSession(session.Id)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, session.Id, result.Id)
	assert.Equal(t, session.Title, result.Title)
	assert.Equal(t, session.Capacity, result.Capacity)
	assert.Equal(t, session.ProfessionalId, result.ProfessionalId)
}

func TestGetSessionNotFound(t *testing.T) {
	/*
		GIVEN: No session exists with the given ID
		WHEN:  GetPostgresqlSession is called with non-existent ID
		THEN:  A not found error is returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewSessionAdapterTestWrapper(t)

	nonExistentId := uuid.New()

	// WHEN
	session, err := adapter.GetPostgresqlSession(nonExistentId)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, session)
	assert.Equal(t, errors.ObjectNotFoundError.SessionNotFound, *err)
}

func TestFetchSessionsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple sessions exist in the database
		WHEN:  FetchPostgresqlSessions is called
		THEN:  All matching sessions are returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewSessionAdapterTestWrapper(t)

	professional1 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	professional2 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	session1 := factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId: &professional1.Id,
	})
	session2 := factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId: &professional2.Id,
	})

	// WHEN - Fetch all sessions
	sessions, err := adapter.FetchPostgresqlSessions(nil, nil, nil, nil)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, sessions)
	assert.GreaterOrEqual(t, len(sessions), 2)

	// Find our created sessions
	foundSession1 := false
	foundSession2 := false
	for _, session := range sessions {
		if session.Id == session1.Id {
			foundSession1 = true
		}
		if session.Id == session2.Id {
			foundSession2 = true
		}
	}
	assert.True(t, foundSession1)
	assert.True(t, foundSession2)
}

func TestFetchSessionsWithProfessionalFilter(t *testing.T) {
	/*
		GIVEN: Sessions exist for different professionals
		WHEN:  FetchPostgresqlSessions is called with professional filter
		THEN:  Only sessions for that professional are returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewSessionAdapterTestWrapper(t)

	professional1 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	professional2 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	session1 := factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId: &professional1.Id,
	})
	factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId: &professional2.Id,
	})

	// WHEN - Filter by professional1
	sessions, err := adapter.FetchPostgresqlSessions([]uuid.UUID{professional1.Id}, nil, nil, nil)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, sessions)
	assert.GreaterOrEqual(t, len(sessions), 1)

	// Verify all returned sessions are for professional1
	for _, session := range sessions {
		assert.Equal(t, professional1.Id, session.ProfessionalId)
	}

	// Verify our session is in the results
	foundSession1 := false
	for _, session := range sessions {
		if session.Id == session1.Id {
			foundSession1 = true
		}
	}
	assert.True(t, foundSession1)
}

func TestFetchSessionsEmpty(t *testing.T) {
	/*
		GIVEN: No sessions exist in the database
		WHEN:  FetchPostgresqlSessions is called
		THEN:  An empty list is returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewSessionAdapterTestWrapper(t)

	// WHEN
	sessions, err := adapter.FetchPostgresqlSessions(nil, nil, nil, nil)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, sessions)
	assert.Equal(t, 0, len(sessions))
}
