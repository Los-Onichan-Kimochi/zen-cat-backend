package controller

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type Reservation struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

// Create Reservation postgresql controller
func NewReservationController(logger logging.Logger, postgresqlDB *gorm.DB) *Reservation {
	return &Reservation{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

// Gets a specific reservation by ID.
func (r *Reservation) GetReservation(reservationId uuid.UUID) (*model.Reservation, error) {
	var reservation model.Reservation
	result := r.PostgresqlDB.Preload("User").
		Preload("Session").
		Where("id = ?", reservationId).
		First(&reservation)
	if result.Error != nil {
		return nil, result.Error
	}

	return &reservation, nil
}

// Fetch all reservations with optional filters.
func (r *Reservation) FetchReservations(
	userIds []uuid.UUID,
	sessionIds []uuid.UUID,
	states []string,
) ([]*model.Reservation, error) {
	reservations := []*model.Reservation{}

	query := r.PostgresqlDB.Model(&model.Reservation{}).Preload("User").Preload("Session")

	if len(userIds) > 0 {
		query = query.Where("user_id IN (?)", userIds)
	}
	if len(sessionIds) > 0 {
		query = query.Where("session_id IN (?)", sessionIds)
	}
	if len(states) > 0 {
		query = query.Where("state IN (?)", states)
	}

	if err := query.Find(&reservations).Error; err != nil {
		return nil, err
	}

	return reservations, nil
}

// Creates a new reservation.
func (r *Reservation) CreateReservation(
	name string,
	reservationTime time.Time,
	state string,
	userId uuid.UUID,
	sessionId uuid.UUID,
	updatedBy string,
) (*model.Reservation, error) {
	reservation := model.Reservation{
		Id:               uuid.New(),
		Name:             name,
		ReservationTime:  reservationTime,
		State:            model.ReservationState(state),
		LastModification: time.Now(),
		UserId:           userId,
		SessionId:        sessionId,
	}

	if err := r.PostgresqlDB.Create(&reservation).Error; err != nil {
		return nil, err
	}

	// Reload with preloaded relationships
	if err := r.PostgresqlDB.Preload("User").Preload("Session").First(&reservation, reservation.Id).Error; err != nil {
		return nil, err
	}

	return &reservation, nil
}

// Updates an existing reservation.
func (r *Reservation) UpdateReservation(
	reservationId uuid.UUID,
	name *string,
	reservationTime *time.Time,
	state *string,
	userId *uuid.UUID,
	sessionId *uuid.UUID,
	updatedBy string,
) (*model.Reservation, error) {
	var reservation model.Reservation
	if err := r.PostgresqlDB.First(&reservation, reservationId).Error; err != nil {
		return nil, err
	}

	// Update fields if provided
	if name != nil {
		reservation.Name = *name
	}
	if reservationTime != nil {
		reservation.ReservationTime = *reservationTime
	}
	if state != nil {
		reservation.State = model.ReservationState(*state)
	}
	if userId != nil {
		reservation.UserId = *userId
	}
	if sessionId != nil {
		reservation.SessionId = *sessionId
	}
	reservation.LastModification = time.Now()

	if err := r.PostgresqlDB.Save(&reservation).Error; err != nil {
		return nil, err
	}

	// Reload with preloaded relationships
	if err := r.PostgresqlDB.Preload("User").Preload("Session").First(&reservation, reservation.Id).Error; err != nil {
		return nil, err
	}

	return &reservation, nil
}

// Deletes a reservation.
func (r *Reservation) DeleteReservation(reservationId uuid.UUID) error {
	result := r.PostgresqlDB.Where("id = ?", reservationId).Delete(&model.Reservation{})
	if result.Error != nil {
		return result.Error
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Bulk deletes reservations.
func (r *Reservation) BulkDeleteReservations(reservationIds []uuid.UUID) error {
	result := r.PostgresqlDB.Where("id IN (?)", reservationIds).Delete(&model.Reservation{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
