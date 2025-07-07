package controller

import (
	"time"

	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
	bllAdapter "onichankimochi.com/astro_cat_backend/src/server/bll/adapter"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type Session struct {
	logger      logging.Logger
	Adapter     *bllAdapter.AdapterCollection
	EnvSettings *schemas.EnvSettings
}

// Create Session controller
func NewSessionController(
	logger logging.Logger,
	adapter *bllAdapter.AdapterCollection,
	envSettings *schemas.EnvSettings,
) *Session {
	return &Session{
		logger:      logger,
		Adapter:     adapter,
		EnvSettings: envSettings,
	}
}

// Creates a session.
func (s *Session) CreateSession(
	req schemas.CreateSessionRequest,
	updatedBy string,
) (*schemas.Session, *errors.Error) {
	// Validate session times
	if !req.EndTime.After(req.StartTime) {
		return nil, &errors.BadRequestError.SessionNotCreated
	}

	// Validate that the professional exists
	_, err := s.Adapter.Professional.GetPostgresqlProfessional(req.ProfessionalId)
	if err != nil {
		return nil, err
	}

	// Validate that the local exists if provided
	if req.LocalId != nil {
		_, err := s.Adapter.Local.GetPostgresqlLocal(*req.LocalId)
		if err != nil {
			return nil, err
		}
	}

	// Validate that the community service exists if provided
	if req.CommunityServiceId != nil {
		// Get the community service to validate it exists
		_, err := s.Adapter.CommunityService.GetPostgresqlCommunityServiceById(*req.CommunityServiceId)
		if err != nil {
			return nil, err
		}
	}

	// Check for conflicts
	conflictCheck := schemas.CheckConflictRequest{
		Date:               req.Date,
		StartTime:          req.StartTime,
		EndTime:            req.EndTime,
		ProfessionalId:     req.ProfessionalId,
		LocalId:            req.LocalId,
		CommunityServiceId: req.CommunityServiceId,
	}

	conflictResult, conflictErr := s.CheckConflicts(conflictCheck)
	if conflictErr != nil {
		return nil, conflictErr
	}

	if conflictResult.HasConflict {
		// Create specific error message
		var conflictDetails []string
		if len(conflictResult.ProfessionalConflicts) > 0 {
			conflictDetails = append(conflictDetails, "conflicto de profesional")
		}
		if len(conflictResult.LocalConflicts) > 0 {
			conflictDetails = append(conflictDetails, "conflicto de local")
		}

		return nil, &errors.ConflictError.SessionTimeConflict
	}

	// Create session if no conflicts
	return s.Adapter.Session.CreatePostgresqlSession(
		req.Title,
		req.Date,
		req.StartTime,
		req.EndTime,
		req.Capacity,
		req.SessionLink,
		req.ProfessionalId,
		req.LocalId,
		req.CommunityServiceId,
		updatedBy,
	)
}

// Gets a session.
func (s *Session) GetSession(sessionId uuid.UUID) (*schemas.Session, *errors.Error) {
	return s.Adapter.Session.GetPostgresqlSession(sessionId)
}

// Updates a session.
func (s *Session) UpdateSession(
	sessionId uuid.UUID,
	req schemas.UpdateSessionRequest,
	updatedBy string,
) (*schemas.Session, *errors.Error) {
	// Validate that the professional exists if provided
	if req.ProfessionalId != nil {
		_, err := s.Adapter.Professional.GetPostgresqlProfessional(*req.ProfessionalId)
		if err != nil {
			return nil, err
		}
	}

	// Validate that the local exists if provided
	if req.LocalId != nil && *req.LocalId != uuid.Nil {
		_, err := s.Adapter.Local.GetPostgresqlLocal(*req.LocalId)
		if err != nil {
			return nil, err
		}
	}

	// Validate that the community service exists if provided
	if req.CommunityServiceId != nil {
		// Get the community service to validate it exists
		_, err := s.Adapter.CommunityService.GetPostgresqlCommunityServiceById(*req.CommunityServiceId)
		if err != nil {
			return nil, err
		}
	}

	// Check for conflicts if relevant fields are being updated
	if req.Date != nil || req.StartTime != nil || req.EndTime != nil || req.ProfessionalId != nil || req.LocalId != nil {
		// Get current session for default values
		currentSession, err := s.Adapter.Session.GetPostgresqlSession(sessionId)
		if err != nil {
			return nil, err
		}

		// Prepare data for conflict check
		checkDate := currentSession.Date
		if req.Date != nil {
			checkDate = *req.Date
		}

		checkStartTime := currentSession.StartTime
		if req.StartTime != nil {
			checkStartTime = *req.StartTime
		}

		checkEndTime := currentSession.EndTime
		if req.EndTime != nil {
			checkEndTime = *req.EndTime
		}

		checkProfessionalId := currentSession.ProfessionalId
		if req.ProfessionalId != nil {
			checkProfessionalId = *req.ProfessionalId
		}

		checkLocalId := currentSession.LocalId
		if req.LocalId != nil {
			checkLocalId = req.LocalId
		}

		conflictCheck := schemas.CheckConflictRequest{
			Date:               checkDate,
			StartTime:          checkStartTime,
			EndTime:            checkEndTime,
			ProfessionalId:     checkProfessionalId,
			LocalId:            checkLocalId,
			CommunityServiceId: currentSession.CommunityServiceId,
			ExcludeId:          &sessionId, // Exclude the current session from conflict checks
		}

		conflictResult, conflictErr := s.CheckConflicts(conflictCheck)
		if conflictErr != nil {
			return nil, conflictErr
		}

		if conflictResult.HasConflict {
			// Create specific error message
			var conflictDetails []string
			if len(conflictResult.ProfessionalConflicts) > 0 {
				conflictDetails = append(conflictDetails, "conflicto de profesional")
			}
			if len(conflictResult.LocalConflicts) > 0 {
				conflictDetails = append(conflictDetails, "conflicto de local")
			}

			return nil, &errors.ConflictError.SessionTimeConflict
		}
	}

	// If the session is being cancelled, update all related reservations to ANNULLED
	if req.State != nil && *req.State == "CANCELLED" {
		// Get all reservations for this session
		reservations, err := s.Adapter.Reservation.FetchPostgresqlReservations(
			[]uuid.UUID{},          // No user filter
			[]uuid.UUID{sessionId}, // Filter by this session
			[]string{},             // No state filter
		)
		if err != nil {
			s.logger.Warn("Error fetching reservations for session", "error", err)
			// Continue with session update even if reservation fetch fails
		} else {
			// Update each reservation to ANULLED state
			annulledState := "ANULLED"
			for _, reservation := range reservations {
				_, updateErr := s.Adapter.Reservation.UpdatePostgresqlReservation(
					reservation.Id,
					nil,            // No name change
					nil,            // No reservation time change
					&annulledState, // Change state to ANNULLED
					nil,            // No user change
					nil,            // No session change
					nil,            // No membership change
					updatedBy,
				)
				if updateErr != nil {
					s.logger.Warn("Error updating reservation status", "reservationId", reservation.Id, "error", updateErr)
					// Continue with other reservations even if one fails
				}
			}
			s.logger.Info("Updated reservations to ANNULLED for cancelled session", "sessionId", sessionId, "count", len(reservations))
		}
	}

	return s.Adapter.Session.UpdatePostgresqlSession(
		sessionId,
		req.Title,
		req.Date,
		req.StartTime,
		req.EndTime,
		req.State,
		req.RegisteredCount,
		req.Capacity,
		req.SessionLink,
		req.ProfessionalId,
		req.LocalId,
		req.CommunityServiceId,
		updatedBy,
	)
}

// Soft deletes a session.
func (s *Session) DeleteSession(sessionId uuid.UUID) *errors.Error {
	return s.Adapter.Session.DeletePostgresqlSession(sessionId)
}

// Bulk deletes sessions.
func (s *Session) BulkDeleteSessions(
	bulkDeleteSessionData schemas.BulkDeleteSessionRequest,
) *errors.Error {
	return s.Adapter.Session.BulkDeletePostgresqlSessions(
		bulkDeleteSessionData.Sessions,
	)
}

// Fetch all sessions, filtered by params.
func (s *Session) FetchSessions(
	professionalIds []string,
	localIds []string,
	communityServiceIds []string,
	states []string,
) (*schemas.Sessions, *errors.Error) {
	// Validate and convert professionalIds to UUIDs if provided.
	parsedProfessionalIds := []uuid.UUID{}
	if len(professionalIds) > 0 {
		for _, id := range professionalIds {
			parsedId, err := uuid.Parse(id)
			if err != nil {
				return nil, &errors.UnprocessableEntityError.InvalidProfessionalId
			}

			// Validate that the professional exists
			_, newErr := s.Adapter.Professional.GetPostgresqlProfessional(parsedId)
			if newErr != nil {
				return nil, newErr
			}

			parsedProfessionalIds = append(parsedProfessionalIds, parsedId)
		}
	}

	// Validate and convert localIds to UUIDs if provided.
	parsedLocalIds := []uuid.UUID{}
	if len(localIds) > 0 {
		for _, id := range localIds {
			parsedId, err := uuid.Parse(id)
			if err != nil {
				return nil, &errors.UnprocessableEntityError.InvalidLocalId
			}

			// Validate that the local exists
			_, newErr := s.Adapter.Local.GetPostgresqlLocal(parsedId)
			if newErr != nil {
				return nil, newErr
			}

			parsedLocalIds = append(parsedLocalIds, parsedId)
		}
	}

	// Validate and convert communityServiceIds to UUIDs if provided.
	parsedCommunityServiceIds := []uuid.UUID{}
	if len(communityServiceIds) > 0 {
		for _, id := range communityServiceIds {
			parsedId, err := uuid.Parse(id)
			if err != nil {
				return nil, &errors.UnprocessableEntityError.InvalidCommunityServiceId
			}

			// Validate that the community service exists
			_, newErr := s.Adapter.CommunityService.GetPostgresqlCommunityServiceById(parsedId)
			if newErr != nil {
				return nil, newErr
			}

			parsedCommunityServiceIds = append(parsedCommunityServiceIds, parsedId)
		}
	}

	sessions, err := s.Adapter.Session.FetchPostgresqlSessions(
		parsedProfessionalIds,
		parsedLocalIds,
		parsedCommunityServiceIds,
		states,
	)
	if err != nil {
		return nil, err
	}

	return &schemas.Sessions{Sessions: sessions}, nil
}

// Creates multiple sessions
func (s *Session) BulkCreateSessions(
	createSessionsData []*schemas.CreateSessionRequest,
	updatedBy string,
) (*schemas.Sessions, *errors.Error) {
	// Validate that all professionals and locals exist
	for _, sessionData := range createSessionsData {
		_, err := s.Adapter.Professional.GetPostgresqlProfessional(sessionData.ProfessionalId)
		if err != nil {
			return nil, err
		}
		if sessionData.LocalId != nil {
			_, err := s.Adapter.Local.GetPostgresqlLocal(*sessionData.LocalId)
			if err != nil {
				return nil, err
			}
		}
	}

	// Check conflicts with database and within the batch
	for i, sessionData := range createSessionsData {
		// 1. Check against database
		conflictCheck := schemas.CheckConflictRequest{
			Date:               sessionData.Date,
			StartTime:          sessionData.StartTime,
			EndTime:            sessionData.EndTime,
			ProfessionalId:     sessionData.ProfessionalId,
			LocalId:            sessionData.LocalId,
			CommunityServiceId: sessionData.CommunityServiceId,
		}
		conflictResult, conflictErr := s.CheckConflicts(conflictCheck)
		if conflictErr != nil {
			return nil, conflictErr
		}
		if conflictResult.HasConflict {
			return nil, &errors.ConflictError.SessionTimeConflict
		}
		// 2. Check against rest of batch (internal conflicts)
		for j, other := range createSessionsData {
			if i == j {
				continue
			}
			// Same day
			if !s.isSameDate(sessionData.Date, other.Date) {
				continue
			}
			// Professional conflict
			if sessionData.ProfessionalId == other.ProfessionalId && s.hasTimeOverlap(sessionData.StartTime, sessionData.EndTime, other.StartTime, other.EndTime) {
				return nil, &errors.ConflictError.SessionTimeConflict
			}
			// Local conflict - only one activity allowed per local at a time
			if sessionData.LocalId != nil && other.LocalId != nil && *sessionData.LocalId == *other.LocalId &&
				s.hasTimeOverlap(sessionData.StartTime, sessionData.EndTime, other.StartTime, other.EndTime) {
				return nil, &errors.ConflictError.SessionTimeConflict
			}
		}
	}

	sessions, err := s.Adapter.Session.BulkCreatePostgresqlSessions(createSessionsData, updatedBy)
	if err != nil {
		return nil, err
	}

	return &schemas.Sessions{Sessions: sessions}, nil
}

// Checks for session conflicts
func (s *Session) CheckConflicts(
	req schemas.CheckConflictRequest,
) (*schemas.ConflictResult, *errors.Error) {
	// Get all sessions for the specific date
	sessions, err := s.Adapter.Session.FetchPostgresqlSessions(
		[]uuid.UUID{},
		[]uuid.UUID{},
		[]uuid.UUID{},
		[]string{},
	)
	if err != nil {
		return nil, err
	}

	professionalConflicts := []*schemas.Session{}
	localConflicts := []*schemas.Session{}

	for _, session := range sessions {
		// Skip if it's the session we're excluding (for edit mode)
		if req.ExcludeId != nil && session.Id == *req.ExcludeId {
			continue
		}

		// Skip cancelled or completed sessions
		if session.State == "CANCELLED" || session.State == "COMPLETED" {
			continue
		}

		// Check if it's the same date
		if !s.isSameDate(session.Date, req.Date) {
			continue
		}

		// Check for time overlap
		if s.hasTimeOverlap(session.StartTime, session.EndTime, req.StartTime, req.EndTime) {
			// Check professional conflict
			if session.ProfessionalId == req.ProfessionalId {
				professionalConflicts = append(professionalConflicts, session)
			}

			// Check local conflict - only one activity allowed per local at a time
			if req.LocalId != nil && session.LocalId != nil && *session.LocalId == *req.LocalId {
				localConflicts = append(localConflicts, session)
			}
		}
	}

	hasConflict := len(professionalConflicts) > 0 || len(localConflicts) > 0

	return &schemas.ConflictResult{
		HasConflict:           hasConflict,
		ProfessionalConflicts: professionalConflicts,
		LocalConflicts:        localConflicts,
	}, nil
}

// Gets availability information for a specific date
func (s *Session) GetAvailability(
	req schemas.AvailabilityRequest,
) (*schemas.AvailabilityResult, *errors.Error) {
	// Get all sessions for the specific date
	sessions, err := s.Adapter.Session.FetchPostgresqlSessions(
		[]uuid.UUID{},
		[]uuid.UUID{},
		[]uuid.UUID{},
		[]string{},
	)
	if err != nil {
		return nil, err
	}

	busySlots := []schemas.TimeSlot{}

	for _, session := range sessions {
		// Skip cancelled or completed sessions
		if session.State == "CANCELLED" || session.State == "COMPLETED" {
			continue
		}

		// Skip the excluded session if provided
		if req.ExcludeSessionId != nil && session.Id == *req.ExcludeSessionId {
			continue
		}

		// Check if it's the same date
		if !s.isSameDate(session.Date, req.Date) {
			continue
		}

		// Add busy slot if matches criteria
		slotType := ""
		shouldAdd := false

		if req.ProfessionalId != nil && session.ProfessionalId == *req.ProfessionalId {
			slotType = "professional"
			shouldAdd = true
		}

		if req.LocalId != nil && session.LocalId != nil && *session.LocalId == *req.LocalId {
			slotType = "local"
			shouldAdd = true
		}

		if shouldAdd {
			busySlots = append(busySlots, schemas.TimeSlot{
				Start: session.StartTime.Format("15:04"),
				End:   session.EndTime.Format("15:04"),
				Title: session.Title,
				Type:  slotType,
			})
		}
	}

	isAvailable := len(busySlots) == 0

	return &schemas.AvailabilityResult{
		IsAvailable: isAvailable,
		BusySlots:   busySlots,
	}, nil
}

// Helper function to check if two dates are the same day
func (s *Session) isSameDate(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// Helper function to check if two time ranges overlap
func (s *Session) hasTimeOverlap(start1, end1, start2, end2 time.Time) bool {
	return start1.Before(end2) && end1.After(start2)
}

// Creates multiple sessions
func (s *Session) BatchCreateSessions(
	req schemas.BatchCreateSessionRequest,
	updatedBy string,
) (*schemas.Sessions, *errors.Error) {
	sessions, err := s.Adapter.Session.BulkCreatePostgresqlSessions(
		req.Sessions,
		updatedBy,
	)
	if err != nil {
		return nil, err
	}

	return &schemas.Sessions{Sessions: sessions}, nil
}

// Fetch all sessions by professional ID.
func (s *Session) FetchSessionsByProfessionalId(
	professionalId uuid.UUID,
) (*schemas.Sessions, *errors.Error) {
	sessions, err := s.Adapter.Session.FetchPostgresqlSessions(
		[]uuid.UUID{professionalId},
		nil,
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &schemas.Sessions{Sessions: sessions}, nil
}

// Fetch all sessions by local ID.
func (s *Session) FetchSessionsByLocalId(
	localId uuid.UUID,
) (*schemas.Sessions, *errors.Error) {
	sessions, err := s.Adapter.Session.FetchPostgresqlSessions(
		nil,
		[]uuid.UUID{localId},
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &schemas.Sessions{Sessions: sessions}, nil
}
