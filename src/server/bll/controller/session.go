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

	return s.Adapter.Session.CreatePostgresqlSession(
		req.Title,
		req.Date,
		req.StartTime,
		req.EndTime,
		req.Capacity,
		req.SessionLink,
		req.ProfessionalId,
		req.LocalId,
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

	return s.Adapter.Session.UpdatePostgresqlSession(
		sessionId,
		req.Title,
		req.Date,
		req.StartTime,
		req.EndTime,
		req.State,
		req.Capacity,
		req.SessionLink,
		req.ProfessionalId,
		req.LocalId,
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

	sessions, err := s.Adapter.Session.FetchPostgresqlSessions(parsedProfessionalIds, parsedLocalIds, states)
	if err != nil {
		return nil, err
	}

	return &schemas.Sessions{Sessions: sessions}, nil
}

// Creates multiple sessions
func (s *Session) BulkCreateSessions(
	createSessionsData []*schemas.CreateSessionRequest,
	updatedBy string,
) ([]*schemas.Session, *errors.Error) {
	// Validate all professionals and locals exist
	for _, sessionData := range createSessionsData {
		// Validate that the professional exists
		_, err := s.Adapter.Professional.GetPostgresqlProfessional(sessionData.ProfessionalId)
		if err != nil {
			return nil, err
		}

		// Validate that the local exists if provided
		if sessionData.LocalId != nil {
			_, err := s.Adapter.Local.GetPostgresqlLocal(*sessionData.LocalId)
			if err != nil {
				return nil, err
			}
		}
	}

	return s.Adapter.Session.BulkCreatePostgresqlSessions(createSessionsData, updatedBy)
}

// Checks for session conflicts
func (s *Session) CheckConflicts(req schemas.CheckConflictRequest) (*schemas.ConflictResult, *errors.Error) {
	// Get all sessions for the specific date
	sessions, err := s.Adapter.Session.FetchPostgresqlSessions([]uuid.UUID{}, []uuid.UUID{}, []string{})
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

			// Check local conflict (if both sessions are presential)
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
func (s *Session) GetAvailability(req schemas.AvailabilityRequest) (*schemas.AvailabilityResult, *errors.Error) {
	// Get all sessions for the specific date
	sessions, err := s.Adapter.Session.FetchPostgresqlSessions([]uuid.UUID{}, []uuid.UUID{}, []string{})
	if err != nil {
		return nil, err
	}

	busySlots := []schemas.TimeSlot{}

	for _, session := range sessions {
		// Skip cancelled or completed sessions
		if session.State == "CANCELLED" || session.State == "COMPLETED" {
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
