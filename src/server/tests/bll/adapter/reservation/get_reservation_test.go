package reservation_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	adapterTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/adapter"
)

func TestGetReservationSuccessfully(t *testing.T) {
	/*
		GIVEN: A reservation exists in the database
		WHEN:  GetPostgresqlReservation is called with the reservation ID
		THEN:  The reservation is returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewReservationAdapterTestWrapper(t)

	reservation := factories.NewReservationModel(db, factories.ReservationModelF{})

	// WHEN
	result, err := adapter.GetPostgresqlReservation(reservation.Id)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, reservation.Id, result.Id)
	assert.Equal(t, reservation.Name, result.Name)
	assert.Equal(t, reservation.UserId, result.UserId)
	assert.Equal(t, reservation.SessionId, result.SessionId)
}

func TestGetReservationNotFound(t *testing.T) {
	/*
		GIVEN: No reservation exists with the given ID
		WHEN:  GetPostgresqlReservation is called with non-existent ID
		THEN:  A not found error is returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewReservationAdapterTestWrapper(t)

	nonExistentId := uuid.New()

	// WHEN
	reservation, err := adapter.GetPostgresqlReservation(nonExistentId)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, reservation)
	assert.Equal(t, errors.ObjectNotFoundError.ReservationNotFound, *err)
}

func TestFetchReservationsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple reservations exist in the database
		WHEN:  FetchPostgresqlReservations is called
		THEN:  All matching reservations are returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewReservationAdapterTestWrapper(t)

	user1 := factories.NewUserModel(db, factories.UserModelF{})
	user2 := factories.NewUserModel(db, factories.UserModelF{})
	session1 := factories.NewSessionModel(db, factories.SessionModelF{})
	session2 := factories.NewSessionModel(db, factories.SessionModelF{})

	reservation1 := factories.NewReservationModel(db, factories.ReservationModelF{
		UserId:    &user1.Id,
		SessionId: &session1.Id,
	})
	reservation2 := factories.NewReservationModel(db, factories.ReservationModelF{
		UserId:    &user2.Id,
		SessionId: &session2.Id,
	})

	// WHEN - Fetch all reservations
	reservations, err := adapter.FetchPostgresqlReservations(nil, nil, nil)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, reservations)
	assert.GreaterOrEqual(t, len(reservations), 2)

	// Find our created reservations
	foundReservation1 := false
	foundReservation2 := false
	for _, reservation := range reservations {
		if reservation.Id == reservation1.Id {
			foundReservation1 = true
		}
		if reservation.Id == reservation2.Id {
			foundReservation2 = true
		}
	}
	assert.True(t, foundReservation1)
	assert.True(t, foundReservation2)
}

func TestFetchReservationsWithUserFilter(t *testing.T) {
	/*
		GIVEN: Reservations exist for different users
		WHEN:  FetchPostgresqlReservations is called with user filter
		THEN:  Only reservations for that user are returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewReservationAdapterTestWrapper(t)

	user1 := factories.NewUserModel(db, factories.UserModelF{})
	user2 := factories.NewUserModel(db, factories.UserModelF{})
	session := factories.NewSessionModel(db, factories.SessionModelF{})

	reservation1 := factories.NewReservationModel(db, factories.ReservationModelF{
		UserId:    &user1.Id,
		SessionId: &session.Id,
	})
	factories.NewReservationModel(db, factories.ReservationModelF{
		UserId:    &user2.Id,
		SessionId: &session.Id,
	})

	// WHEN - Filter by user1
	reservations, err := adapter.FetchPostgresqlReservations([]uuid.UUID{user1.Id}, nil, nil)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, reservations)
	assert.GreaterOrEqual(t, len(reservations), 1)

	// Verify all returned reservations are for user1
	for _, reservation := range reservations {
		assert.Equal(t, user1.Id, reservation.UserId)
	}

	// Verify our reservation is in the results
	foundReservation1 := false
	for _, reservation := range reservations {
		if reservation.Id == reservation1.Id {
			foundReservation1 = true
		}
	}
	assert.True(t, foundReservation1)
}

func TestFetchReservationsWithSessionFilter(t *testing.T) {
	/*
		GIVEN: Reservations exist for different sessions
		WHEN:  FetchPostgresqlReservations is called with session filter
		THEN:  Only reservations for that session are returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewReservationAdapterTestWrapper(t)

	user := factories.NewUserModel(db, factories.UserModelF{})
	session1 := factories.NewSessionModel(db, factories.SessionModelF{})
	session2 := factories.NewSessionModel(db, factories.SessionModelF{})

	reservation1 := factories.NewReservationModel(db, factories.ReservationModelF{
		UserId:    &user.Id,
		SessionId: &session1.Id,
	})
	factories.NewReservationModel(db, factories.ReservationModelF{
		UserId:    &user.Id,
		SessionId: &session2.Id,
	})

	// WHEN - Filter by session1
	reservations, err := adapter.FetchPostgresqlReservations(nil, []uuid.UUID{session1.Id}, nil)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, reservations)
	assert.GreaterOrEqual(t, len(reservations), 1)

	// Verify all returned reservations are for session1
	for _, reservation := range reservations {
		assert.Equal(t, session1.Id, reservation.SessionId)
	}

	// Verify our reservation is in the results
	foundReservation1 := false
	for _, reservation := range reservations {
		if reservation.Id == reservation1.Id {
			foundReservation1 = true
		}
	}
	assert.True(t, foundReservation1)
}

func TestFetchReservationsWithStateFilter(t *testing.T) {
	/*
		GIVEN: Reservations exist with different states
		WHEN:  FetchPostgresqlReservations is called with state filter
		THEN:  Only reservations with that state are returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewReservationAdapterTestWrapper(t)

	user := factories.NewUserModel(db, factories.UserModelF{})
	session := factories.NewSessionModel(db, factories.SessionModelF{})

	reservation1 := factories.NewReservationModel(db, factories.ReservationModelF{
		UserId:    &user.Id,
		SessionId: &session.Id,
	})
	factories.NewReservationModel(db, factories.ReservationModelF{
		UserId:    &user.Id,
		SessionId: &session.Id,
	})

	// WHEN - Filter by CONFIRMED state
	reservations, err := adapter.FetchPostgresqlReservations(nil, nil, []string{"CONFIRMED"})

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, reservations)
	assert.GreaterOrEqual(t, len(reservations), 1)

	// Verify all returned reservations have CONFIRMED state
	for _, reservation := range reservations {
		assert.Equal(t, "CONFIRMED", reservation.State)
	}

	// Verify our reservation is in the results
	foundReservation1 := false
	for _, reservation := range reservations {
		if reservation.Id == reservation1.Id {
			foundReservation1 = true
		}
	}
	assert.True(t, foundReservation1)
}

func TestFetchReservationsEmpty(t *testing.T) {
	/*
		GIVEN: No reservations exist in the database
		WHEN:  FetchPostgresqlReservations is called
		THEN:  An empty list is returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewReservationAdapterTestWrapper(t)

	// WHEN
	reservations, err := adapter.FetchPostgresqlReservations(nil, nil, nil)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, reservations)
	assert.Equal(t, 0, len(reservations))
}
