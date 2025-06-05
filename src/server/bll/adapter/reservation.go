package adapter

import (
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
