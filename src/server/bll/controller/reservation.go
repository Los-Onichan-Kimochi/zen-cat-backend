package controller

import (
	"time"

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
func (r *Reservation) GetReservation(
	reservationId uuid.UUID,
) (*schemas.Reservation, *errors.Error) {
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

	reservations, err := r.Adapter.Reservation.FetchPostgresqlReservations(
		parsedUserIds,
		parsedSessionIds,
		states,
	)
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
	// Validate updatedBy is not empty
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	// Validate required fields
	if createReservationData.Name == "" {
		return nil, &errors.BadRequestError.UserNotCreated // Use existing error for validation
	}

	// Validate that the user exists
	_, userErr := r.Adapter.User.GetPostgresqlUser(createReservationData.UserId)
	if userErr != nil {
		return nil, userErr
	}

	// Validate that the session exists
	session, sessionErr := r.Adapter.Session.GetPostgresqlSession(createReservationData.SessionId)
	if sessionErr != nil {
		return nil, sessionErr
	}

	// Validate that the membership exists if provided
	var membershipToUpdate *schemas.Membership
	if createReservationData.MembershipId != nil {
		membership, membershipErr := r.Adapter.Membership.GetPostgresqlMembership(*createReservationData.MembershipId)
		if membershipErr != nil {
			return nil, membershipErr
		}
		membershipToUpdate = membership
	}

	// Check for user reservation conflicts (user cannot be in two sessions at the same time)
	userReservations, reservationErr := r.Adapter.Reservation.FetchPostgresqlReservations(
		[]uuid.UUID{createReservationData.UserId},
		[]uuid.UUID{},
		[]string{"CONFIRMED"}, // Only CONFIRMED reservations block new reservations
	)
	if reservationErr != nil {
		return nil, reservationErr
	}

	// Check each existing reservation for time conflicts
	for _, existingReservation := range userReservations {
		if existingReservation.SessionId == createReservationData.SessionId {
			continue
		}

		existingSession, existingSessionErr := r.Adapter.Session.GetPostgresqlSession(existingReservation.SessionId)
		if existingSessionErr != nil {
			continue
		}

		if r.isSameDate(session.Date, existingSession.Date) {
			if r.hasTimeOverlap(session.StartTime, session.EndTime, existingSession.StartTime, existingSession.EndTime) {
				return nil, &errors.ConflictError.UserReservationTimeConflict
			}
		}
	}

	// Modify `registered_count` field of the session only if the reservation is CONFIRMED
	if createReservationData.State == "CONFIRMED" {
		session.RegisteredCount++
		_, sessionErr = r.Adapter.Session.UpdatePostgresqlSession(
			createReservationData.SessionId,
			nil,
			nil,
			nil,
			nil,
			nil,
			&session.RegisteredCount,
			nil,
			nil,
			nil,
			nil,
			nil,
			updatedBy,
		)
		if sessionErr != nil {
			r.logger.Error("Failed to increment registered_count for session", sessionErr)
			// Continue with reservation creation even if session update fails
		}
	}

	// Create the reservation
	newReservation, createErr := r.Adapter.Reservation.CreatePostgresqlReservation(
		createReservationData.Name,
		createReservationData.ReservationTime,
		createReservationData.State,
		createReservationData.UserId,
		createReservationData.SessionId,
		createReservationData.MembershipId,
		updatedBy,
	)

	if createErr != nil {
		return nil, createErr
	}

	// Incrementar reservations_used de la membership si existe y si no es ilimitada (null)
	if membershipToUpdate != nil && membershipToUpdate.ReservationsUsed != nil && createReservationData.State == "CONFIRMED" {
		currentUsed := *membershipToUpdate.ReservationsUsed
		newUsed := currentUsed + 1

		_, updateErr := r.Adapter.Membership.UpdatePostgresqlMembership(
			*createReservationData.MembershipId,
			nil,
			nil,
			nil,
			nil,
			&newUsed,
			nil,
			nil,
			nil,
			updatedBy,
		)

		if updateErr != nil {
			r.logger.Error("Failed to update reservations_used for membership", updateErr)
			// No devolvemos error para no afectar la creación de la reserva
		}
	}

	return newReservation, nil
}

// Updates a reservation.
func (r *Reservation) UpdateReservation(
	reservationId uuid.UUID,
	updateReservationData schemas.UpdateReservationRequest,
	updatedBy string,
) (*schemas.Reservation, *errors.Error) {
	// Obtener la reserva actual para comparar estados
	currentReservation, getErr := r.Adapter.Reservation.GetPostgresqlReservation(reservationId)
	if getErr != nil {
		return nil, getErr
	}

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

	// Verificar si hay un cambio de membresía
	var oldMembershipId *uuid.UUID
	var newMembershipId *uuid.UUID

	oldMembershipId = currentReservation.MembershipId
	if updateReservationData.MembershipId != nil {
		newMembershipId = updateReservationData.MembershipId
	} else {
		newMembershipId = oldMembershipId
	}

	// Validate that the membership exists if provided
	var oldMembership *schemas.Membership
	var newMembership *schemas.Membership

	if oldMembershipId != nil {
		membership, membershipErr := r.Adapter.Membership.GetPostgresqlMembership(*oldMembershipId)
		if membershipErr != nil {
			return nil, membershipErr
		}
		oldMembership = membership
	}

	if newMembershipId != nil && (oldMembershipId == nil || *newMembershipId != *oldMembershipId) {
		membership, membershipErr := r.Adapter.Membership.GetPostgresqlMembership(*newMembershipId)
		if membershipErr != nil {
			return nil, membershipErr
		}
		newMembership = membership
	} else if newMembershipId != nil {
		newMembership = oldMembership
	}

	// Verificar si hay un cambio de estado
	oldState := currentReservation.State
	var newState string
	if updateReservationData.State != nil {
		newState = *updateReservationData.State
	} else {
		newState = oldState
	}

	// Realizar la actualización de la reserva
	updatedReservation, updateErr := r.Adapter.Reservation.UpdatePostgresqlReservation(
		reservationId,
		updateReservationData.Name,
		updateReservationData.ReservationTime,
		updateReservationData.State,
		updateReservationData.UserId,
		updateReservationData.SessionId,
		updateReservationData.MembershipId,
		updatedBy,
	)

	if updateErr != nil {
		return nil, updateErr
	}

	// Actualizar el contador de reservas usadas en la membresía según los cambios de estado

	// Caso 1: Estado cambia de confirmado a anulado o cancelado
	if oldState == "CONFIRMED" && (newState == "ANULLED" || newState == "CANCELLED") {
		// Decrementar registered_count de la sesión
		session, sessionErr := r.Adapter.Session.GetPostgresqlSession(currentReservation.SessionId)
		if sessionErr == nil && session.RegisteredCount > 0 {
			session.RegisteredCount--
			_, updateSessionErr := r.Adapter.Session.UpdatePostgresqlSession(
				currentReservation.SessionId,
				nil, // title
				nil, // date
				nil, // start_time
				nil, // end_time
				nil, // state
				&session.RegisteredCount, // decrementar contador
				nil, // capacity
				nil, // session_link
				nil, // professional_id
				nil, // local_id
				nil, // community_service_id
				updatedBy,
			)
			if updateSessionErr != nil {
				r.logger.Error("Failed to decrement registered_count for session", updateSessionErr)
			}
		}

		// Decrementar contador si la membresía tiene un contador (no es ilimitada)
		if oldMembership != nil && oldMembership.ReservationsUsed != nil {
			currentUsed := *oldMembership.ReservationsUsed
			newUsed := currentUsed - 1
			if newUsed < 0 {
				newUsed = 0
			}

			_, membershipErr := r.Adapter.Membership.UpdatePostgresqlMembership(
				*oldMembershipId,
				nil,
				nil,
				nil,
				nil,
				&newUsed,
				nil,
				nil,
				nil,
				updatedBy,
			)

			if membershipErr != nil {
				r.logger.Error("Failed to decrement reservations_used for membership", membershipErr)
				// No devolvemos error para no afectar la actualización de la reserva
			}
		}
	}

	// Caso 2: Estado cambia de anulado o cancelado a confirmado
	if (oldState == "ANULLED" || oldState == "CANCELLED") && newState == "CONFIRMED" {
		// Incrementar registered_count de la sesión
		session, sessionErr := r.Adapter.Session.GetPostgresqlSession(currentReservation.SessionId)
		if sessionErr == nil {
			session.RegisteredCount++
			_, updateSessionErr := r.Adapter.Session.UpdatePostgresqlSession(
				currentReservation.SessionId,
				nil, // title
				nil, // date
				nil, // start_time
				nil, // end_time
				nil, // state
				&session.RegisteredCount, // incrementar contador
				nil, // capacity
				nil, // session_link
				nil, // professional_id
				nil, // local_id
				nil, // community_service_id
				updatedBy,
			)
			if updateSessionErr != nil {
				r.logger.Error("Failed to increment registered_count for session", updateSessionErr)
			}
		}

		// Incrementar contador si la nueva membresía tiene un contador (no es ilimitada)
		if newMembership != nil && newMembership.ReservationsUsed != nil {
			currentUsed := *newMembership.ReservationsUsed
			newUsed := currentUsed + 1

			_, membershipErr := r.Adapter.Membership.UpdatePostgresqlMembership(
				*newMembershipId,
				nil,
				nil,
				nil,
				nil,
				&newUsed,
				nil,
				nil,
				nil,
				updatedBy,
			)

			if membershipErr != nil {
				r.logger.Error("Failed to increment reservations_used for membership", membershipErr)
				// No devolvemos error para no afectar la actualización de la reserva
			}
		}
	}

	// Caso 3: Cambio de membresía manteniendo estado confirmado
	if oldState == "CONFIRMED" && newState == "CONFIRMED" &&
		oldMembershipId != nil && newMembershipId != nil &&
		*oldMembershipId != *newMembershipId {

		// Decrementar contador en la membresía antigua si no es ilimitada
		if oldMembership != nil && oldMembership.ReservationsUsed != nil {
			currentUsed := *oldMembership.ReservationsUsed
			newUsed := currentUsed - 1
			if newUsed < 0 {
				newUsed = 0
			}

			_, membershipErr := r.Adapter.Membership.UpdatePostgresqlMembership(
				*oldMembershipId,
				nil,
				nil,
				nil,
				nil,
				&newUsed,
				nil,
				nil,
				nil,
				updatedBy,
			)

			if membershipErr != nil {
				r.logger.Error("Failed to decrement reservations_used for old membership", membershipErr)
			}
		}

		// Incrementar contador en la nueva membresía si no es ilimitada
		if newMembership != nil && newMembership.ReservationsUsed != nil {
			currentUsed := *newMembership.ReservationsUsed
			newUsed := currentUsed + 1

			_, membershipErr := r.Adapter.Membership.UpdatePostgresqlMembership(
				*newMembershipId,
				nil,
				nil,
				nil,
				nil,
				&newUsed,
				nil,
				nil,
				nil,
				updatedBy,
			)

			if membershipErr != nil {
				r.logger.Error("Failed to increment reservations_used for new membership", membershipErr)
			}
		}
	}

	return updatedReservation, nil
}

// Deletes a reservation.
func (r *Reservation) DeleteReservation(reservationId uuid.UUID) *errors.Error {
	// Obtener la reserva antes de eliminarla para obtener info de la membresía
	currentReservation, getErr := r.Adapter.Reservation.GetPostgresqlReservation(reservationId)
	if getErr != nil {
		return getErr
	}

	// Si la reserva está confirmada, decrementar registered_count y membresía
	if currentReservation.State == "CONFIRMED" {
		// Decrementar registered_count de la sesión
		session, sessionErr := r.Adapter.Session.GetPostgresqlSession(currentReservation.SessionId)
		if sessionErr == nil && session.RegisteredCount > 0 {
			session.RegisteredCount--
			_, updateSessionErr := r.Adapter.Session.UpdatePostgresqlSession(
				currentReservation.SessionId,
				nil, // title
				nil, // date
				nil, // start_time
				nil, // end_time
				nil, // state
				&session.RegisteredCount, // decrementar contador
				nil, // capacity
				nil, // session_link
				nil, // professional_id
				nil, // local_id
				nil, // community_service_id
				"SYSTEM", // En caso de eliminación, usamos SYSTEM como updatedBy
			)
			if updateSessionErr != nil {
				r.logger.Error("Failed to decrement registered_count on reservation delete", updateSessionErr)
			}
		}

		// Decrementar contador de membresía si tiene una asociada
		if currentReservation.MembershipId != nil {
			// Obtener la membresía actual
			membership, membershipErr := r.Adapter.Membership.GetPostgresqlMembership(*currentReservation.MembershipId)
			if membershipErr == nil && membership.ReservationsUsed != nil {
				// Decrementar el contador si no es membresía ilimitada
				currentUsed := *membership.ReservationsUsed
				newUsed := currentUsed - 1
				if newUsed < 0 {
					newUsed = 0
				}

				// Actualizar la membresía
				_, updateErr := r.Adapter.Membership.UpdatePostgresqlMembership(
					*currentReservation.MembershipId,
					nil,
					nil,
					nil,
					nil,
					&newUsed,
					nil,
					nil,
					nil,
					"SYSTEM", // En caso de eliminación, usamos SYSTEM como updatedBy
				)

				if updateErr != nil {
					r.logger.Error("Failed to decrement reservations_used on reservation delete", updateErr)
					// No devolvemos error para no afectar la eliminación de la reserva
				}
			}
		}
	}

	return r.Adapter.Reservation.DeletePostgresqlReservation(reservationId)
}

// Bulk deletes reservations.
func (r *Reservation) BulkDeleteReservations(
	bulkDeleteReservationData schemas.BulkDeleteReservationRequest,
) *errors.Error {
	// Convertir las IDs de string a UUID
	uuidIds := make([]uuid.UUID, 0, len(bulkDeleteReservationData.Reservations))
	for _, idStr := range bulkDeleteReservationData.Reservations {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return &errors.UnprocessableEntityError.InvalidReservationId
		}
		uuidIds = append(uuidIds, id)
	}

	// Para cada reserva confirmada, decrementar el contador de la membresía
	for _, reservationId := range uuidIds {
		// Obtener la reserva antes de eliminarla
		reservation, getErr := r.Adapter.Reservation.GetPostgresqlReservation(reservationId)
		if getErr != nil {
			// Ignorar errores y continuar con las siguientes
			r.logger.Error("Failed to get reservation for bulk delete", getErr)
			continue
		}

		// Si está confirmada, decrementar registered_count y membresía
		if reservation.State == "CONFIRMED" {
			// Decrementar registered_count de la sesión
			session, sessionErr := r.Adapter.Session.GetPostgresqlSession(reservation.SessionId)
			if sessionErr == nil && session.RegisteredCount > 0 {
				session.RegisteredCount--
				_, updateSessionErr := r.Adapter.Session.UpdatePostgresqlSession(
					reservation.SessionId,
					nil, // title
					nil, // date
					nil, // start_time
					nil, // end_time
					nil, // state
					&session.RegisteredCount, // decrementar contador
					nil, // capacity
					nil, // session_link
					nil, // professional_id
					nil, // local_id
					nil, // community_service_id
					"SYSTEM", // En caso de eliminación en bloque, usamos SYSTEM como updatedBy
				)
				if updateSessionErr != nil {
					r.logger.Error("Failed to decrement registered_count on bulk delete", updateSessionErr)
				}
			}

			// Decrementar contador de membresía si tiene una asociada
			if reservation.MembershipId != nil {
				// Obtener la membresía
				membership, membershipErr := r.Adapter.Membership.GetPostgresqlMembership(*reservation.MembershipId)
				if membershipErr == nil && membership.ReservationsUsed != nil {
					// Decrementar el contador si no es ilimitada
					currentUsed := *membership.ReservationsUsed
					newUsed := currentUsed - 1
					if newUsed < 0 {
						newUsed = 0
					}

					// Actualizar la membresía
					_, updateErr := r.Adapter.Membership.UpdatePostgresqlMembership(
						*reservation.MembershipId,
						nil,
						nil,
						nil,
						nil,
						&newUsed,
						nil,
						nil,
						nil,
						"SYSTEM", // En caso de eliminación en bloque, usamos SYSTEM como updatedBy
					)

					if updateErr != nil {
						r.logger.Error("Failed to decrement reservations_used on bulk delete", updateErr)
						// Continuar con la siguiente reserva
					}
				}
			}
		}
	}

	return r.Adapter.Reservation.BulkDeletePostgresqlReservations(
		bulkDeleteReservationData.Reservations,
	)
}

// Helper function to check if two dates are the same day
func (r *Reservation) isSameDate(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// Helper function to check if two time ranges overlap
func (r *Reservation) hasTimeOverlap(start1, end1, start2, end2 time.Time) bool {
	return start1.Before(end2) && end1.After(start2)
}

// GetReservationsByCommunityIdByUserId obtiene las reservas de un usuario para una comunidad específica
func (r *Reservation) GetReservationsByCommunityIdByUserId(
	communityId uuid.UUID,
	userId uuid.UUID,
) (*schemas.Reservations, *errors.Error) {
	// Verificar que la comunidad existe
	_, communityErr := r.Adapter.Community.GetPostgresqlCommunity(communityId)
	if communityErr != nil {
		return nil, communityErr
	}

	// Verificar que el usuario existe
	_, userErr := r.Adapter.User.GetPostgresqlUser(userId)
	if userErr != nil {
		return nil, userErr
	}

	// 1. Obtener los servicios de la comunidad
	communityServices, err := r.Adapter.CommunityService.FetchPostgresqlCommunityServices(
		&communityId,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Si no hay servicios para esta comunidad, devolver una lista vacía
	if len(communityServices) == 0 {
		return &schemas.Reservations{Reservations: []*schemas.Reservation{}}, nil
	}

	// 2. Extraer los IDs de servicios de comunidad
	communityServiceIds := make([]uuid.UUID, len(communityServices))
	for i, cs := range communityServices {
		communityServiceIds[i] = cs.Id
	}

	// 3. Obtener las sesiones asociadas a estos servicios de comunidad
	sessions, sessionsErr := r.Adapter.Session.FetchPostgresqlSessions(
		[]uuid.UUID{}, // No filtrar por profesional
		[]uuid.UUID{}, // No filtrar por local
		communityServiceIds,
		[]string{}, // Todos los estados
	)
	if sessionsErr != nil {
		return nil, sessionsErr
	}

	// Si no hay sesiones, devolver una lista vacía
	if len(sessions) == 0 {
		return &schemas.Reservations{Reservations: []*schemas.Reservation{}}, nil
	}

	// 4. Extraer los IDs de sesión
	sessionIds := make([]uuid.UUID, len(sessions))
	for i, session := range sessions {
		sessionIds[i] = session.Id
	}

	// 5. Obtener las reservas para estas sesiones y este usuario
	reservations, reservationsErr := r.Adapter.Reservation.FetchPostgresqlReservations(
		[]uuid.UUID{userId},
		sessionIds,
		[]string{}, // Todos los estados
	)
	if reservationsErr != nil {
		return nil, reservationsErr
	}

	return &schemas.Reservations{Reservations: reservations}, nil
}

// GetServiceReport obtiene el reporte de servicios para el dashboard admin
type ServiceReportResponse struct {
	Total    int                            `json:"totalReservations"`
	Services []bllAdapter.ServiceReportData `json:"services"`
}

func (r *Reservation) GetServiceReport(from, to *time.Time, groupBy string) (*ServiceReportResponse, *errors.Error) {
	params := bllAdapter.ServiceReportParams{
		From:    from,
		To:      to,
		GroupBy: groupBy,
	}
	total, services, err := r.Adapter.Reservation.GetServiceReport(params)
	if err != nil {
		return nil, &errors.InternalServerError.Default
	}
	return &ServiceReportResponse{
		Total:    total,
		Services: services,
	}, nil
}
