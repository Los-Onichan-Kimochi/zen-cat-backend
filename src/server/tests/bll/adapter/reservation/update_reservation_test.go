package reservation_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	adapterTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/adapter"
)

func TestUpdateReservationSuccessfully(t *testing.T) {
	/*
		GIVEN: A reservation exists in the database
		WHEN:  UpdatePostgresqlReservation is called with new data
		THEN:  The reservation is updated and returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewReservationAdapterTestWrapper(t)

	reservation := factories.NewReservationModel(db, factories.ReservationModelF{})

	newName := "Updated Reservation Name"
	newState := "CONFIRMED"
	updatedBy := "test-admin"

	// WHEN
	updatedReservation, err := adapter.UpdatePostgresqlReservation(
		reservation.Id,
		&newName,
		nil, // Don't update reservation time
		&newState,
		nil, // Don't update user
		nil, // Don't update session
		updatedBy,
	)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, updatedReservation)
	assert.Equal(t, reservation.Id, updatedReservation.Id)
	assert.Equal(t, newName, updatedReservation.Name)
	assert.Equal(t, newState, updatedReservation.State)
}

func TestUpdateReservationWithCompleteData(t *testing.T) {
	/*
		GIVEN: A reservation exists in the database
		WHEN:  UpdatePostgresqlReservation is called with all fields
		THEN:  The reservation is completely updated
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewReservationAdapterTestWrapper(t)

	reservation := factories.NewReservationModel(db, factories.ReservationModelF{})
	newUser := factories.NewUserModel(db, factories.UserModelF{})
	newSession := factories.NewSessionModel(db, factories.SessionModelF{})

	newName := "Completely Updated Reservation"
	newReservationTime := time.Now().AddDate(0, 0, 3)
	newState := "RESCHEDULED"
	updatedBy := "test-admin"

	// WHEN
	updatedReservation, err := adapter.UpdatePostgresqlReservation(
		reservation.Id,
		&newName,
		&newReservationTime,
		&newState,
		&newUser.Id,
		&newSession.Id,
		updatedBy,
	)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, updatedReservation)
	assert.Equal(t, reservation.Id, updatedReservation.Id)
	assert.Equal(t, newName, updatedReservation.Name)
	assert.Equal(t, newState, updatedReservation.State)
	assert.Equal(t, newUser.Id, updatedReservation.UserId)
	assert.Equal(t, newSession.Id, updatedReservation.SessionId)
	assert.Equal(t, newReservationTime.Format("2006-01-02 15:04:05"), updatedReservation.ReservationTime.Format("2006-01-02 15:04:05"))
}

func TestUpdateReservationStateTransitions(t *testing.T) {
	/*
		GIVEN: A reservation exists in the database
		WHEN:  UpdatePostgresqlReservation is called with different states
		THEN:  The reservation state is updated correctly
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewReservationAdapterTestWrapper(t)

	reservation := factories.NewReservationModel(db, factories.ReservationModelF{})
	updatedBy := "test-admin"

	states := []string{"PENDING", "CONFIRMED", "CANCELLED", "COMPLETED"}

	for _, state := range states {
		// WHEN
		updatedReservation, err := adapter.UpdatePostgresqlReservation(
			reservation.Id,
			nil, // Don't update name
			nil, // Don't update time
			&state,
			nil, // Don't update user
			nil, // Don't update session
			updatedBy,
		)

		// THEN
		assert.Nil(t, err)
		assert.NotNil(t, updatedReservation)
		assert.Equal(t, state, updatedReservation.State)
	}
}

func TestDeleteReservationSuccessfully(t *testing.T) {
	/*
		GIVEN: A reservation exists in the database
		WHEN:  DeletePostgresqlReservation is called
		THEN:  The reservation is deleted
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewReservationAdapterTestWrapper(t)

	reservation := factories.NewReservationModel(db, factories.ReservationModelF{})

	// WHEN
	err := adapter.DeletePostgresqlReservation(reservation.Id)

	// THEN
	assert.Nil(t, err)

	// Verify reservation is deleted by trying to get it
	deletedReservation, getErr := adapter.GetPostgresqlReservation(reservation.Id)
	assert.NotNil(t, getErr)
	assert.Nil(t, deletedReservation)
	assert.Equal(t, errors.ObjectNotFoundError.ReservationNotFound, *getErr)
}

func TestDeleteReservationNotFound(t *testing.T) {
	/*
		GIVEN: No reservation exists with the given ID
		WHEN:  DeletePostgresqlReservation is called
		THEN:  A not found error is returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewReservationAdapterTestWrapper(t)

	nonExistentId := uuid.New()

	// WHEN
	err := adapter.DeletePostgresqlReservation(nonExistentId)

	// THEN
	assert.NotNil(t, err)
	assert.Equal(t, errors.ObjectNotFoundError.ReservationNotFound, *err)
}

func TestBulkDeleteReservationsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple reservations exist in the database
		WHEN:  BulkDeletePostgresqlReservations is called
		THEN:  All specified reservations are deleted
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewReservationAdapterTestWrapper(t)

	reservation1 := factories.NewReservationModel(db, factories.ReservationModelF{})
	reservation2 := factories.NewReservationModel(db, factories.ReservationModelF{})
	reservation3 := factories.NewReservationModel(db, factories.ReservationModelF{})

	reservationIds := []string{
		reservation1.Id.String(),
		reservation2.Id.String(),
	}

	// WHEN
	err := adapter.BulkDeletePostgresqlReservations(reservationIds)

	// THEN
	assert.Nil(t, err)

	// Verify deleted reservations cannot be found
	_, getErr1 := adapter.GetPostgresqlReservation(reservation1.Id)
	assert.NotNil(t, getErr1)
	assert.Equal(t, errors.ObjectNotFoundError.ReservationNotFound, *getErr1)

	_, getErr2 := adapter.GetPostgresqlReservation(reservation2.Id)
	assert.NotNil(t, getErr2)
	assert.Equal(t, errors.ObjectNotFoundError.ReservationNotFound, *getErr2)

	// Verify non-deleted reservation still exists
	reservation3Result, getErr3 := adapter.GetPostgresqlReservation(reservation3.Id)
	assert.Nil(t, getErr3)
	assert.NotNil(t, reservation3Result)
	assert.Equal(t, reservation3.Id, reservation3Result.Id)
}

func TestBulkDeleteReservationsWithInvalidId(t *testing.T) {
	/*
		GIVEN: Invalid reservation ID in the list
		WHEN:  BulkDeletePostgresqlReservations is called
		THEN:  An error is returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewReservationAdapterTestWrapper(t)

	invalidIds := []string{"invalid-uuid", "another-invalid-id"}

	// WHEN
	err := adapter.BulkDeletePostgresqlReservations(invalidIds)

	// THEN
	assert.NotNil(t, err)
	assert.Equal(t, errors.UnprocessableEntityError.InvalidReservationId, *err)
} 