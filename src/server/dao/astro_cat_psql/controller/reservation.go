package controller

import (
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
	result := r.PostgresqlDB.Preload("User").Preload("Session").Where("id = ?", reservationId).First(&reservation)
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
