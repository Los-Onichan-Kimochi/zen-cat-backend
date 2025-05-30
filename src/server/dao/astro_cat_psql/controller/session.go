package controller

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type Session struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

// Create Session postgresql controller
func NewSessionController(logger logging.Logger, postgresqlDB *gorm.DB) *Session {
	return &Session{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

// Creates a session given its model.
func (s *Session) CreateSession(session *model.Session) error {
	return s.PostgresqlDB.Create(session).Error
}

// Gets a session model given params.
func (s *Session) GetSession(sessionId uuid.UUID) (*model.Session, error) {
	session := &model.Session{}

	result := s.PostgresqlDB.Preload("Professional").Preload("Local").First(&session, "id = ?", sessionId)
	if result.Error != nil {
		return nil, result.Error
	}

	return session, nil
}

// Updates session given fields to update.
func (s *Session) UpdateSession(
	id uuid.UUID,
	title *string,
	date *time.Time,
	startTime *time.Time,
	endTime *time.Time,
	state *string,
	capacity *int,
	sessionLink *string,
	professionalId *uuid.UUID,
	localId *uuid.UUID,
	updatedBy string,
) (*model.Session, error) {
	updateFields := map[string]any{
		"updated_by": updatedBy,
	}
	if title != nil {
		updateFields["title"] = *title
	}
	if date != nil {
		updateFields["date"] = *date
	}
	if startTime != nil {
		updateFields["start_time"] = *startTime
	}
	if endTime != nil {
		updateFields["end_time"] = *endTime
	}
	if state != nil {
		updateFields["state"] = *state
	}
	if capacity != nil {
		updateFields["capacity"] = *capacity
	}
	if sessionLink != nil {
		updateFields["session_link"] = *sessionLink
	}
	if professionalId != nil {
		updateFields["professional_id"] = *professionalId
	}
	if localId != nil {
		updateFields["local_id"] = *localId
	}

	// Check if there are any fields to update
	var session model.Session
	if len(updateFields) == 1 {
		if err := s.PostgresqlDB.Preload("Professional").Preload("Local").First(&session, "id = ?", id).Error; err != nil {
			return nil, err
		}
		return &session, nil
	}

	// Perform the update
	result := s.PostgresqlDB.Model(&session).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(updateFields)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &session, nil
}

// Soft deletes a session given its ID.
func (s *Session) DeleteSession(sessionId uuid.UUID) error {
	result := s.PostgresqlDB.Where("id = ?", sessionId).Delete(&model.Session{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// TODO: Add filters and sorting.
// Fetch all sessions with optional filters.
func (s *Session) FetchSessions(
	professionalIds []uuid.UUID,
	localIds []uuid.UUID,
	states []string,
) ([]*model.Session, error) {
	sessions := []*model.Session{}

	query := s.PostgresqlDB.Model(&model.Session{}).Preload("Professional").Preload("Local")

	if len(professionalIds) > 0 {
		query = query.Where("professional_id IN (?)", professionalIds)
	}
	if len(localIds) > 0 {
		query = query.Where("local_id IN (?)", localIds)
	}
	if len(states) > 0 {
		query = query.Where("state IN (?)", states)
	}

	if err := query.Find(&sessions).Error; err != nil {
		return nil, err
	}

	return sessions, nil
}

// Creates sessions given their models.
func (s *Session) BulkCreateSessions(sessions []*model.Session) error {
	return s.PostgresqlDB.Create(&sessions).Error
}
