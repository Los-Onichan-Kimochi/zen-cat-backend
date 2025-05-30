package controller

import (
	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
	bllAdapter "onichankimochi.com/astro_cat_backend/src/server/bll/adapter"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type Reservation struct {
	logger      logging.Logger
	Adapter     *bllAdapter.AdapterCollection
	EnvSettings *schemas.EnvSettings
}

// Create Reservation controller
func NewReservationController(
	logger logging.Logger,
	adapter *bllAdapter.AdapterCollection,
	envSettings *schemas.EnvSettings,
) *Reservation {
	return &Reservation{
		logger:      logger,
		Adapter:     adapter,
		EnvSettings: envSettings,
	}
}

// Gets a specific reservation.
func (r *Reservation) GetReservation(reservationId uuid.UUID) (*schemas.Reservation, *errors.Error) {
	return r.Adapter.Reservation.GetPostgresqlReservation(reservationId)
}

// Fetch all reservations, filtered by params.
func (r *Reservation) FetchReservations(
	userIds []string,
	sessionIds []string,
	states []string,
) (*schemas.Reservations, *errors.Error) {
	// Validate and convert userIds to UUIDs if provided.
	parsedUserIds := []uuid.UUID{}
	if len(userIds) > 0 {
		for _, id := range userIds {
			parsedId, err := uuid.Parse(id)
			if err != nil {
				return nil, &errors.UnprocessableEntityError.InvalidUserId
			}

			// Validate that the user exists
			_, newErr := r.Adapter.User.GetPostgresqlUser(parsedId)
			if newErr != nil {
				return nil, newErr
			}

			parsedUserIds = append(parsedUserIds, parsedId)
		}
	}

	// Validate and convert sessionIds to UUIDs if provided.
	parsedSessionIds := []uuid.UUID{}
	if len(sessionIds) > 0 {
		for _, id := range sessionIds {
			parsedId, err := uuid.Parse(id)
			if err != nil {
				return nil, &errors.UnprocessableEntityError.InvalidSessionId
			}

			// Validate that the session exists
			_, newErr := r.Adapter.Session.GetPostgresqlSession(parsedId)
			if newErr != nil {
				return nil, newErr
			}

			parsedSessionIds = append(parsedSessionIds, parsedId)
		}
	}

	reservations, err := r.Adapter.Reservation.FetchPostgresqlReservations(parsedUserIds, parsedSessionIds, states)
	if err != nil {
		return nil, err
	}

	return &schemas.Reservations{Reservations: reservations}, nil
}
