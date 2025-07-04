package adapter

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
	daoPsql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
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
func (r *Reservation) GetPostgresqlReservation(
	reservationId uuid.UUID,
) (*schemas.Reservation, *errors.Error) {
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
		Session: schemas.Session{
			Id:                 reservationModel.Session.Id,
			Title:              reservationModel.Session.Title,
			Date:               reservationModel.Session.Date,
			StartTime:          reservationModel.Session.StartTime,
			EndTime:            reservationModel.Session.EndTime,
			State:              string(reservationModel.Session.State),
			RegisteredCount:    reservationModel.Session.RegisteredCount,
			Capacity:           reservationModel.Session.Capacity,
			SessionLink:        reservationModel.Session.SessionLink,
			ProfessionalId:     reservationModel.Session.ProfessionalId,
			LocalId:            reservationModel.Session.LocalId,
			CommunityServiceId: reservationModel.Session.CommunityServiceId,
		},
		MembershipId: reservationModel.MembershipId,
	}, nil
}

// Fetch all reservations from postgresql DB and adapts them to Reservation schema.
func (r *Reservation) FetchPostgresqlReservations(
	userIds []uuid.UUID,
	sessionIds []uuid.UUID,
	states []string,
) ([]*schemas.Reservation, *errors.Error) {
	reservationModels, err := r.DaoPostgresql.Reservation.FetchReservations(
		userIds,
		sessionIds,
		states,
	)
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
			Session: schemas.Session{
				Id:                 reservationModel.Session.Id,
				Title:              reservationModel.Session.Title,
				Date:               reservationModel.Session.Date,
				StartTime:          reservationModel.Session.StartTime,
				EndTime:            reservationModel.Session.EndTime,
				State:              string(reservationModel.Session.State),
				RegisteredCount:    reservationModel.Session.RegisteredCount,
				Capacity:           reservationModel.Session.Capacity,
				SessionLink:        reservationModel.Session.SessionLink,
				ProfessionalId:     reservationModel.Session.ProfessionalId,
				LocalId:            reservationModel.Session.LocalId,
				CommunityServiceId: reservationModel.Session.CommunityServiceId,
			},
			MembershipId: reservationModel.MembershipId,
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
	membershipId *uuid.UUID,
	updatedBy string,
) (*schemas.Reservation, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	reservationModel, err := r.DaoPostgresql.Reservation.CreateReservation(
		name,
		reservationTime,
		state,
		userId,
		sessionId,
		membershipId,
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
		Session: schemas.Session{
			Id:                 reservationModel.Session.Id,
			Title:              reservationModel.Session.Title,
			Date:               reservationModel.Session.Date,
			StartTime:          reservationModel.Session.StartTime,
			EndTime:            reservationModel.Session.EndTime,
			State:              string(reservationModel.Session.State),
			RegisteredCount:    reservationModel.Session.RegisteredCount,
			Capacity:           reservationModel.Session.Capacity,
			SessionLink:        reservationModel.Session.SessionLink,
			ProfessionalId:     reservationModel.Session.ProfessionalId,
			LocalId:            reservationModel.Session.LocalId,
			CommunityServiceId: reservationModel.Session.CommunityServiceId,
		},
		MembershipId: reservationModel.MembershipId,
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
	membershipId *uuid.UUID,
	updatedBy string,
) (*schemas.Reservation, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	reservationModel, err := r.DaoPostgresql.Reservation.UpdateReservation(
		reservationId,
		name,
		reservationTime,
		state,
		userId,
		sessionId,
		membershipId,
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
		Session: schemas.Session{
			Id:                 reservationModel.Session.Id,
			Title:              reservationModel.Session.Title,
			Date:               reservationModel.Session.Date,
			StartTime:          reservationModel.Session.StartTime,
			EndTime:            reservationModel.Session.EndTime,
			State:              string(reservationModel.Session.State),
			RegisteredCount:    reservationModel.Session.RegisteredCount,
			Capacity:           reservationModel.Session.Capacity,
			SessionLink:        reservationModel.Session.SessionLink,
			ProfessionalId:     reservationModel.Session.ProfessionalId,
			LocalId:            reservationModel.Session.LocalId,
			CommunityServiceId: reservationModel.Session.CommunityServiceId,
		},
		MembershipId: reservationModel.MembershipId,
	}, nil
}

// Deletes a reservation from postgresql DB.
func (r *Reservation) DeletePostgresqlReservation(reservationId uuid.UUID) *errors.Error {
	err := r.DaoPostgresql.Reservation.DeleteReservation(reservationId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &errors.ObjectNotFoundError.ReservationNotFound
		}
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

// GetServiceReport obtiene reservas, hace preload profundo y agrupa por servicio y fecha
type ServiceReportParams struct {
	From    *time.Time
	To      *time.Time
	GroupBy string // "day", "week", "month"
}

type ServiceReportData struct {
	ServiceType string
	Total       int
	Data        []struct {
		Date  string
		Count int
	}
}

func (r *Reservation) GetServiceReport(params ServiceReportParams) (total int, services []ServiceReportData, err error) {
	// Construir la consulta base
	query := r.DaoPostgresql.Reservation.PostgresqlDB.Model(&model.Reservation{})
	query = query.Preload("Session.CommunityService.Service")

	// Filtrar por fecha real de la reserva (reservation_time)
	if params.From != nil {
		query = query.Where("reservation_time >= ?", *params.From)
	}
	if params.To != nil {
		query = query.Where("reservation_time <= ?", *params.To)
	}

	var reservations []model.Reservation
	if err := query.Find(&reservations).Error; err != nil {
		return 0, nil, err
	}

	// Agrupar por fecha real de la reserva
	totals := map[string]int{}
	grouped := map[string]map[string]int{} // serviceType -> date -> count

	for _, res := range reservations {
		serviceName := "Desconocido"
		if res.Session.CommunityService != nil && res.Session.CommunityService.Service.Name != "" {
			serviceName = res.Session.CommunityService.Service.Name
		}
		var dateKey string
		switch params.GroupBy {
		case "month":
			dateKey = res.ReservationTime.Format("2006-01")
		case "week":
			y, w := res.ReservationTime.ISOWeek()
			dateKey = fmt.Sprintf("%d-W%02d", y, w)
		default:
			dateKey = res.ReservationTime.Format("2006-01-02")
		}
		totals[serviceName]++
		if grouped[serviceName] == nil {
			grouped[serviceName] = map[string]int{}
		}
		grouped[serviceName][dateKey]++
		total++
	}

	// Construir la respuesta
	for service, dateMap := range grouped {
		data := []struct {
			Date  string
			Count int
		}{}
		for date, count := range dateMap {
			data = append(data, struct {
				Date  string
				Count int
			}{date, count})
		}
		services = append(services, ServiceReportData{
			ServiceType: service,
			Total:       totals[service],
			Data:        data,
		})
	}
	return total, services, nil
}
