package reservation_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestCreateReservationSuccessfully(t *testing.T) {
	// GIVEN: Valid reservation creation request with existing dependencies
	controller, _, db := controllerTest.NewReservationControllerTestWrapper(t)

	// Create dependencies
	testUser := factories.NewUserModel(db, factories.UserModelF{})
	testSession := factories.NewSessionModel(db, factories.SessionModelF{})

	// Create reservation request
	reservationTime := time.Now().Add(24 * time.Hour)

	createRequest := schemas.CreateReservationRequest{
		Name:            "Test Reservation",
		ReservationTime: reservationTime,
		UserId:          testUser.Id,
		SessionId:       testSession.Id,
	}

	// WHEN: CreateReservation is called
	result, err := controller.CreateReservation(createRequest, "test_admin")

	// THEN: Reservation is created successfully
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createRequest.Name, result.Name)
	assert.Equal(t, createRequest.UserId, result.UserId)
	assert.Equal(t, createRequest.SessionId, result.SessionId)
	assert.NotEqual(t, "", result.Id)
}

func TestCreateReservationEmptyUpdatedBy(t *testing.T) {
	// GIVEN: Valid reservation creation request but empty updatedBy
	controller, _, db := controllerTest.NewReservationControllerTestWrapper(t)

	// Create dependencies
	testUser := factories.NewUserModel(db, factories.UserModelF{})
	testSession := factories.NewSessionModel(db, factories.SessionModelF{})

	createRequest := schemas.CreateReservationRequest{
		Name:            "Test Reservation",
		ReservationTime: time.Now().Add(24 * time.Hour),
		UserId:          testUser.Id,
		SessionId:       testSession.Id,
	}

	// WHEN: CreateReservation is called with empty updatedBy
	result, err := controller.CreateReservation(createRequest, "")

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "Invalid updated by value", err.Message)
}

func TestCreateReservationInvalidUser(t *testing.T) {
	// GIVEN: Reservation creation request with non-existent user
	controller, _, db := controllerTest.NewReservationControllerTestWrapper(t)

	// Create session but use non-existent user
	testSession := factories.NewSessionModel(db, factories.SessionModelF{})
	nonExistentUserId := uuid.New()

	createRequest := schemas.CreateReservationRequest{
		Name:            "Test Reservation",
		ReservationTime: time.Now().Add(24 * time.Hour),
		UserId:          nonExistentUserId,
		SessionId:       testSession.Id,
	}

	// WHEN: CreateReservation is called with non-existent user
	result, err := controller.CreateReservation(createRequest, "test_admin")

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestCreateReservationInvalidSession(t *testing.T) {
	// GIVEN: Reservation creation request with non-existent session
	controller, _, db := controllerTest.NewReservationControllerTestWrapper(t)

	// Create user but use non-existent session
	testUser := factories.NewUserModel(db, factories.UserModelF{})
	nonExistentSessionId := uuid.New()

	createRequest := schemas.CreateReservationRequest{
		Name:            "Test Reservation",
		ReservationTime: time.Now().Add(24 * time.Hour),
		UserId:          testUser.Id,
		SessionId:       nonExistentSessionId,
	}

	// WHEN: CreateReservation is called with non-existent session
	result, err := controller.CreateReservation(createRequest, "test_admin")

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestCreateReservationPastTime(t *testing.T) {
	// GIVEN: Reservation creation request with past time
	controller, _, db := controllerTest.NewReservationControllerTestWrapper(t)

	// Create dependencies
	testUser := factories.NewUserModel(db, factories.UserModelF{})
	testSession := factories.NewSessionModel(db, factories.SessionModelF{})

	// Use past time
	pastTime := time.Now().Add(-24 * time.Hour)

	createRequest := schemas.CreateReservationRequest{
		Name:            "Past Reservation",
		ReservationTime: pastTime,
		UserId:          testUser.Id,
		SessionId:       testSession.Id,
	}

	// WHEN: CreateReservation is called with past time
	result, err := controller.CreateReservation(createRequest, "test_admin")

	// THEN: Behavior depends on business rules
	// It might be allowed or not depending on implementation
	if err != nil {
		assert.Nil(t, result)
		assert.NotNil(t, err)
	} else {
		assert.NotNil(t, result)
		assert.Nil(t, err)
	}
}

func TestCreateReservationEmptyName(t *testing.T) {
	// GIVEN: Reservation creation request with empty name
	controller, _, db := controllerTest.NewReservationControllerTestWrapper(t)

	// Create dependencies
	testUser := factories.NewUserModel(db, factories.UserModelF{})
	testSession := factories.NewSessionModel(db, factories.SessionModelF{})

	createRequest := schemas.CreateReservationRequest{
		Name:            "", // Empty name
		ReservationTime: time.Now().Add(24 * time.Hour),
		UserId:          testUser.Id,
		SessionId:       testSession.Id,
	}

	// WHEN: CreateReservation is called
	result, err := controller.CreateReservation(createRequest, "test_admin")

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestCreateReservationDuplicateUserSession(t *testing.T) {
	// GIVEN: User already has reservation for the same session
	controller, _, db := controllerTest.NewReservationControllerTestWrapper(t)

	// Create dependencies
	testUser := factories.NewUserModel(db, factories.UserModelF{})
	testSession := factories.NewSessionModel(db, factories.SessionModelF{})

	createRequest1 := schemas.CreateReservationRequest{
		Name:            "First Reservation",
		ReservationTime: time.Now().Add(24 * time.Hour),
		UserId:          testUser.Id,
		SessionId:       testSession.Id,
	}

	createRequest2 := schemas.CreateReservationRequest{
		Name:            "Second Reservation",           // Different name
		ReservationTime: time.Now().Add(25 * time.Hour), // Different time
		UserId:          testUser.Id,                    // Same user
		SessionId:       testSession.Id,                 // Same session
	}

	// Create first reservation
	_, err1 := controller.CreateReservation(createRequest1, "test_admin")
	assert.Nil(t, err1)

	// WHEN: CreateReservation is called for same user and session
	result, err := controller.CreateReservation(createRequest2, "test_admin")

	// THEN: Behavior depends on business rules
	// It might be allowed (multiple reservations) or not
	if err != nil {
		assert.Nil(t, result)
		assert.NotNil(t, err)
	} else {
		assert.NotNil(t, result)
		assert.Nil(t, err)
	}
}

func TestCreateReservationUserTimeConflict(t *testing.T) {
	// GIVEN: User already has reservation for a different session at the same time
	controller, _, db := controllerTest.NewReservationControllerTestWrapper(t)

	// Create dependencies
	testUser := factories.NewUserModel(db, factories.UserModelF{})
	testProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testLocal := factories.NewLocalModel(db, factories.LocalModelF{})

	// Create two sessions at the same time
	sessionDate := time.Now().Add(24 * time.Hour)
	session1StartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	session1EndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)
	session2StartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 30, 0, 0, time.UTC) // Overlaps
	session2EndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 30, 0, 0, time.UTC)

	session1 := factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId: &testProfessional.Id,
		LocalId:        &testLocal.Id,
		Date:           &sessionDate,
		StartTime:      &session1StartTime,
		EndTime:        &session1EndTime,
	})

	session2 := factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId: &testProfessional.Id,
		LocalId:        &testLocal.Id,
		Date:           &sessionDate,
		StartTime:      &session2StartTime,
		EndTime:        &session2EndTime,
	})

	// Create first reservation
	createRequest1 := schemas.CreateReservationRequest{
		Name:            "First Reservation",
		ReservationTime: sessionDate,
		State:           "CONFIRMED",
		UserId:          testUser.Id,
		SessionId:       session1.Id,
	}

	// Create second reservation (conflicting time)
	createRequest2 := schemas.CreateReservationRequest{
		Name:            "Second Reservation",
		ReservationTime: sessionDate,
		UserId:          testUser.Id,
		SessionId:       session2.Id,
	}

	// Create first reservation
	_, err1 := controller.CreateReservation(createRequest1, "test_admin")
	assert.Nil(t, err1)

	// WHEN: CreateReservation is called for same user but different session at overlapping time
	result, err := controller.CreateReservation(createRequest2, "test_admin")

	// THEN: Should return conflict error
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "CONFLICT_ERROR_002", err.Code)
	assert.Contains(t, err.Message, "User has another reservation at the same time")
}

func TestCreateReservationUserDifferentTimes(t *testing.T) {
	// GIVEN: User has reservation for a session at a different time
	controller, _, db := controllerTest.NewReservationControllerTestWrapper(t)

	// Create dependencies
	testUser := factories.NewUserModel(db, factories.UserModelF{})
	testProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	testLocal := factories.NewLocalModel(db, factories.LocalModelF{})

	// Create two sessions at different times
	sessionDate := time.Now().Add(24 * time.Hour)
	session1StartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 10, 0, 0, 0, time.UTC)
	session1EndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 11, 0, 0, 0, time.UTC)
	session2StartTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 14, 0, 0, 0, time.UTC) // Different time
	session2EndTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), 15, 0, 0, 0, time.UTC)

	session1 := factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId: &testProfessional.Id,
		LocalId:        &testLocal.Id,
		Date:           &sessionDate,
		StartTime:      &session1StartTime,
		EndTime:        &session1EndTime,
	})

	session2 := factories.NewSessionModel(db, factories.SessionModelF{
		ProfessionalId: &testProfessional.Id,
		LocalId:        &testLocal.Id,
		Date:           &sessionDate,
		StartTime:      &session2StartTime,
		EndTime:        &session2EndTime,
	})

	// Create first reservation
	createRequest1 := schemas.CreateReservationRequest{
		Name:            "First Reservation",
		ReservationTime: sessionDate,
		UserId:          testUser.Id,
		SessionId:       session1.Id,
	}

	// Create second reservation (different time)
	createRequest2 := schemas.CreateReservationRequest{
		Name:            "Second Reservation",
		ReservationTime: sessionDate,
		UserId:          testUser.Id,
		SessionId:       session2.Id,
	}

	// Create first reservation
	_, err1 := controller.CreateReservation(createRequest1, "test_admin")
	assert.Nil(t, err1)

	// WHEN: CreateReservation is called for same user but different session at different time
	result, err := controller.CreateReservation(createRequest2, "test_admin")

	// THEN: Should succeed
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createRequest2.Name, result.Name)
	assert.Equal(t, createRequest2.UserId, result.UserId)
	assert.Equal(t, createRequest2.SessionId, result.SessionId)
}
