package reservation_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	adapterTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/adapter"
)

func TestCreateReservationSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid reservation data with existing user and session
		WHEN:  CreatePostgresqlReservation is called
		THEN:  A new reservation is created and returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewReservationAdapterTestWrapper(t)

	user := factories.NewUserModel(db, factories.UserModelF{})
	session := factories.NewSessionModel(db, factories.SessionModelF{})
	membership := factories.NewMembershipModel(db, factories.MembershipModelF{})

	name := "Test Reservation"
	reservationTime := time.Now().AddDate(0, 0, 1) // Tomorrow
	state := "CONFIRMED"
	updatedBy := "test-admin"

	// WHEN
	reservation, err := adapter.CreatePostgresqlReservation(
		name,
		reservationTime,
		state,
		user.Id,
		session.Id,
		&membership.Id,
		updatedBy,
	)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, reservation)
	assert.NotEmpty(t, reservation.Id)
	assert.Equal(t, name, reservation.Name)
	assert.Equal(t, state, reservation.State)
	assert.Equal(t, user.Id, reservation.UserId)
	assert.Equal(t, session.Id, reservation.SessionId)
	assert.Equal(t, reservationTime.Format("2006-01-02 15:04:05"), reservation.ReservationTime.Format("2006-01-02 15:04:05"))
	assert.Equal(t, &membership.Id, reservation.MembershipId)
}

func TestCreateReservationEmptyUpdatedBy(t *testing.T) {
	/*
		GIVEN: Valid reservation data but empty updatedBy
		WHEN:  CreatePostgresqlReservation is called with empty updatedBy
		THEN:  An error is returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewReservationAdapterTestWrapper(t)

	// Create dependencies
	user := factories.NewUserModel(db, factories.UserModelF{})
	session := factories.NewSessionModel(db, factories.SessionModelF{})
	membership := factories.NewMembershipModel(db, factories.MembershipModelF{})

	name := "Test Reservation"
	reservationTime := time.Now()
	state := "CONFIRMED"
	emptyUpdatedBy := ""

	// WHEN
	reservation, err := adapter.CreatePostgresqlReservation(
		name,
		reservationTime,
		state,
		user.Id,
		session.Id,
		&membership.Id,
		emptyUpdatedBy,
	)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, reservation)
	assert.Equal(t, errors.BadRequestError.InvalidUpdatedByValue.Code, err.Code)
}

func TestCreateReservationWithDifferentStates(t *testing.T) {
	/*
		GIVEN: Valid reservation data with different states
		WHEN:  CreatePostgresqlReservation is called
		THEN:  Reservations are created with correct states
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewReservationAdapterTestWrapper(t)

	user := factories.NewUserModel(db, factories.UserModelF{})
	session := factories.NewSessionModel(db, factories.SessionModelF{})
	membership := factories.NewMembershipModel(db, factories.MembershipModelF{})

	states := []string{"PENDING", "CONFIRMED", "CANCELLED"}
	updatedBy := "test-admin"

	for _, state := range states {
		// WHEN
		reservation, err := adapter.CreatePostgresqlReservation(
			"Test Reservation "+state,
			time.Now().AddDate(0, 0, 1),
			state,
			user.Id,
			session.Id,
			&membership.Id,
			updatedBy,
		)

		// THEN
		assert.Nil(t, err)
		assert.NotNil(t, reservation)
		assert.Equal(t, state, reservation.State)
	}
}

func TestCreateReservationWithPastDate(t *testing.T) {
	/*
		GIVEN: Valid reservation data with past date
		WHEN:  CreatePostgresqlReservation is called
		THEN:  A reservation is created (business logic allows past dates)
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewReservationAdapterTestWrapper(t)

	user := factories.NewUserModel(db, factories.UserModelF{})
	session := factories.NewSessionModel(db, factories.SessionModelF{})
	membership := factories.NewMembershipModel(db, factories.MembershipModelF{})

	name := "Past Reservation"
	reservationTime := time.Now().AddDate(0, 0, -1) // Yesterday
	state := "COMPLETED"
	updatedBy := "test-admin"

	// WHEN
	reservation, err := adapter.CreatePostgresqlReservation(
		name,
		reservationTime,
		state,
		user.Id,
		session.Id,
		&membership.Id,
		updatedBy,
	)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, reservation)
	assert.Equal(t, name, reservation.Name)
	assert.Equal(t, state, reservation.State)
}
