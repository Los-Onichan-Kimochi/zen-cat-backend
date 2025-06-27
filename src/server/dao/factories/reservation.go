package factories

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type ReservationModelF struct {
	Id               *uuid.UUID
	Name             *string
	ReservationTime  *time.Time
	State            *model.ReservationState
	LastModification *time.Time
	UserId           *uuid.UUID
	SessionId        *uuid.UUID
}

// Create a new reservation on DB
func NewReservationModel(db *gorm.DB, option ...ReservationModelF) *model.Reservation {
	// Create default user if not provided
	user := NewUserModel(db)

	// Create default session if not provided
	session := NewSessionModel(db)

	now := time.Now()

	reservation := &model.Reservation{
		Id:               uuid.New(),
		Name:             "Test Reservation",
		ReservationTime:  now,
		State:            model.ReservationStateConfirmed,
		LastModification: now,
		UserId:           user.Id,
		SessionId:        session.Id,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}

	if len(option) > 0 {
		parameters := option[0]
		if parameters.Id != nil {
			reservation.Id = *parameters.Id
		}
		if parameters.Name != nil {
			reservation.Name = *parameters.Name
		}
		if parameters.ReservationTime != nil {
			reservation.ReservationTime = *parameters.ReservationTime
		}
		if parameters.State != nil {
			reservation.State = *parameters.State
		}
		if parameters.LastModification != nil {
			reservation.LastModification = *parameters.LastModification
		}
		if parameters.UserId != nil {
			reservation.UserId = *parameters.UserId
		}
		if parameters.SessionId != nil {
			reservation.SessionId = *parameters.SessionId
		}
	}

	result := db.Create(reservation)
	if result.Error != nil {
		log.Fatalf("Error when trying to create reservation: %v", result.Error)
	}

	return reservation
}

// Create size number of new reservations on DB
func NewReservationModelBatch(
	db *gorm.DB,
	size int,
	option ...ReservationModelF,
) []*model.Reservation {
	reservations := []*model.Reservation{}
	for i := 0; i < size; i++ {
		var reservation *model.Reservation
		if len(option) > 0 {
			reservation = NewReservationModel(db, option[0])
		} else {
			reservation = NewReservationModel(db)
		}
		reservations = append(reservations, reservation)
	}
	return reservations
}
