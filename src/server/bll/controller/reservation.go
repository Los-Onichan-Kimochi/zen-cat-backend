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

// Creates a reservation.
func (r *Reservation) CreateReservation(
	createReservationData schemas.CreateReservationRequest,
	updatedBy string,
) (*schemas.Reservation, *errors.Error) {
	// Validate that the user exists
	_, userErr := r.Adapter.User.GetPostgresqlUser(createReservationData.UserId)
	if userErr != nil {
		return nil, userErr
	}

	// Validate that the session exists
	_, sessionErr := r.Adapter.Session.GetPostgresqlSession(createReservationData.SessionId)
	if sessionErr != nil {
		return nil, sessionErr
	}

	return r.Adapter.Reservation.CreatePostgresqlReservation(
		createReservationData.Name,
		createReservationData.ReservationTime,
		createReservationData.State,
		createReservationData.UserId,
		createReservationData.SessionId,
		updatedBy,
	)
}

// Updates a reservation.
func (r *Reservation) UpdateReservation(
	reservationId uuid.UUID,
	updateReservationData schemas.UpdateReservationRequest,
	updatedBy string,
) (*schemas.Reservation, *errors.Error) {
	// Validate that the user exists if provided
	if updateReservationData.UserId != nil {
		_, userErr := r.Adapter.User.GetPostgresqlUser(*updateReservationData.UserId)
		if userErr != nil {
			return nil, userErr
		}
	}

	// Validate that the session exists if provided
	if updateReservationData.SessionId != nil {
		_, sessionErr := r.Adapter.Session.GetPostgresqlSession(*updateReservationData.SessionId)
		if sessionErr != nil {
			return nil, sessionErr
		}
	}

	return r.Adapter.Reservation.UpdatePostgresqlReservation(
		reservationId,
		updateReservationData.Name,
		updateReservationData.ReservationTime,
		updateReservationData.State,
		updateReservationData.UserId,
		updateReservationData.SessionId,
		updatedBy,
	)
}

// Deletes a reservation.
func (r *Reservation) DeleteReservation(reservationId uuid.UUID) *errors.Error {
	return r.Adapter.Reservation.DeletePostgresqlReservation(reservationId)
}

// Bulk deletes reservations.
func (r *Reservation) BulkDeleteReservations(
	bulkDeleteReservationData schemas.BulkDeleteReservationRequest,
) *errors.Error {
	return r.Adapter.Reservation.BulkDeletePostgresqlReservations(
		bulkDeleteReservationData.Reservations,
	)
}
