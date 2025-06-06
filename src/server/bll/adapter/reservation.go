package adapter

import (
	"time"

	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
	daoPsql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type Reservation struct {
	logger        logging.Logger
	DaoPostgresql *daoPsql.AstroCatPsqlCollection
}

// Create Reservation adapter
func NewReservationAdapter(
	logger logging.Logger,
	daoPostgresql *daoPsql.AstroCatPsqlCollection,
) *Reservation {
	return &Reservation{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

// Gets a specific reservation and adapts it.
func (r *Reservation) GetPostgresqlReservation(reservationId uuid.UUID) (*schemas.Reservation, *errors.Error) {
	reservationModel, err := r.DaoPostgresql.Reservation.GetReservation(reservationId)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.ReservationNotFound
	}

	return &schemas.Reservation{
		Id:               reservationModel.Id,
		Name:             reservationModel.Name,
		ReservationTime:  reservationModel.ReservationTime,
		State:            string(reservationModel.State),
		LastModification: reservationModel.LastModification,
		UserId:           reservationModel.UserId,
		SessionId:        reservationModel.SessionId,
	}, nil
}

// Fetch all reservations from postgresql DB and adapts them to Reservation schema.
func (r *Reservation) FetchPostgresqlReservations(
	userIds []uuid.UUID,
	sessionIds []uuid.UUID,
	states []string,
) ([]*schemas.Reservation, *errors.Error) {
	reservationModels, err := r.DaoPostgresql.Reservation.FetchReservations(userIds, sessionIds, states)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.ReservationNotFound
	}

	reservations := make([]*schemas.Reservation, len(reservationModels))
	for i, reservationModel := range reservationModels {
		reservations[i] = &schemas.Reservation{
			Id:               reservationModel.Id,
			Name:             reservationModel.Name,
			ReservationTime:  reservationModel.ReservationTime,
			State:            string(reservationModel.State),
			LastModification: reservationModel.LastModification,
			UserId:           reservationModel.UserId,
			SessionId:        reservationModel.SessionId,
		}
	}

	return reservations, nil
}

// Creates a reservation in postgresql DB and adapts it to Reservation schema.
func (r *Reservation) CreatePostgresqlReservation(
	name string,
	reservationTime time.Time,
	state string,
	userId uuid.UUID,
	sessionId uuid.UUID,
	updatedBy string,
) (*schemas.Reservation, *errors.Error) {
	reservationModel, err := r.DaoPostgresql.Reservation.CreateReservation(
		name,
		reservationTime,
		state,
		userId,
		sessionId,
		updatedBy,
	)
	if err != nil {
		return nil, &errors.InternalServerError.Default
	}

	return &schemas.Reservation{
		Id:               reservationModel.Id,
		Name:             reservationModel.Name,
		ReservationTime:  reservationModel.ReservationTime,
		State:            string(reservationModel.State),
		LastModification: reservationModel.LastModification,
		UserId:           reservationModel.UserId,
		SessionId:        reservationModel.SessionId,
	}, nil
}

// Updates a reservation in postgresql DB and adapts it to Reservation schema.
func (r *Reservation) UpdatePostgresqlReservation(
	reservationId uuid.UUID,
	name *string,
	reservationTime *time.Time,
	state *string,
	userId *uuid.UUID,
	sessionId *uuid.UUID,
	updatedBy string,
) (*schemas.Reservation, *errors.Error) {
	reservationModel, err := r.DaoPostgresql.Reservation.UpdateReservation(
		reservationId,
		name,
		reservationTime,
		state,
		userId,
		sessionId,
		updatedBy,
	)
	if err != nil {
		return nil, &errors.InternalServerError.Default
	}

	return &schemas.Reservation{
		Id:               reservationModel.Id,
		Name:             reservationModel.Name,
		ReservationTime:  reservationModel.ReservationTime,
		State:            string(reservationModel.State),
		LastModification: reservationModel.LastModification,
		UserId:           reservationModel.UserId,
		SessionId:        reservationModel.SessionId,
	}, nil
}

// Deletes a reservation from postgresql DB.
func (r *Reservation) DeletePostgresqlReservation(reservationId uuid.UUID) *errors.Error {
	err := r.DaoPostgresql.Reservation.DeleteReservation(reservationId)
	if err != nil {
		return &errors.InternalServerError.Default
	}

	return nil
}

// Bulk deletes reservations from postgresql DB.
func (r *Reservation) BulkDeletePostgresqlReservations(reservationIds []string) *errors.Error {
	// Convert string IDs to UUIDs
	uuidIds := make([]uuid.UUID, len(reservationIds))
	for i, id := range reservationIds {
		parsedId, parseErr := uuid.Parse(id)
		if parseErr != nil {
			return &errors.UnprocessableEntityError.InvalidReservationId
		}
		uuidIds[i] = parsedId
	}

	err := r.DaoPostgresql.Reservation.BulkDeleteReservations(uuidIds)
	if err != nil {
		return &errors.InternalServerError.Default
	}

	return nil
}
