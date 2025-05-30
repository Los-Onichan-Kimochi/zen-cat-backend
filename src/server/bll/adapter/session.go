package adapter

import (
	"time"

	"github.com/google/uuid"
	daoPostgresql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"

	"onichankimochi.com/astro_cat_backend/src/logging"
)

type Session struct {
	logger        logging.Logger
	DaoPostgresql *daoPostgresql.AstroCatPsqlCollection
}

// Creates Session adapter
func NewSessionAdapter(
	logger logging.Logger,
	daoPostgresql *daoPostgresql.AstroCatPsqlCollection,
) *Session {
	return &Session{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

// Creates a session into postgresql DB and returns it.
func (s *Session) CreatePostgresqlSession(
	title string,
	date time.Time,
	startTime time.Time,
	endTime time.Time,
	capacity int,
	sessionLink *string,
	professionalId uuid.UUID,
	localId *uuid.UUID,
	updatedBy string,
) (*schemas.Session, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	sessionModel := &model.Session{
		Id:              uuid.New(),
		Title:           title,
		Date:            date,
		StartTime:       startTime,
		EndTime:         endTime,
		State:           model.SessionStateScheduled, // Default state
		RegisteredCount: 0,                           // Default number of initial registrations
		Capacity:        capacity,
		SessionLink:     sessionLink,
		ProfessionalId:  professionalId,
		LocalId:         localId,
		AuditFields: model.AuditFields{
			UpdatedBy: updatedBy,
		},
	}

	err := s.DaoPostgresql.Session.CreateSession(sessionModel)
	if err != nil {
		return nil, &errors.BadRequestError.SessionNotCreated
	}

	return &schemas.Session{
		Id:              sessionModel.Id,
		Title:           sessionModel.Title,
		Date:            sessionModel.Date,
		StartTime:       sessionModel.StartTime,
		EndTime:         sessionModel.EndTime,
		State:           string(sessionModel.State),
		RegisteredCount: sessionModel.RegisteredCount,
		Capacity:        sessionModel.Capacity,
		SessionLink:     sessionModel.SessionLink,
		ProfessionalId:  sessionModel.ProfessionalId,
		LocalId:         sessionModel.LocalId,
	}, nil
}

// Gets a session from postgresql DB and adapts it to a Session schema.
func (s *Session) GetPostgresqlSession(sessionId uuid.UUID) (*schemas.Session, *errors.Error) {
	sessionModel, err := s.DaoPostgresql.Session.GetSession(sessionId)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.SessionNotFound
	}

	return &schemas.Session{
		Id:              sessionModel.Id,
		Title:           sessionModel.Title,
		Date:            sessionModel.Date,
		StartTime:       sessionModel.StartTime,
		EndTime:         sessionModel.EndTime,
		State:           string(sessionModel.State),
		RegisteredCount: sessionModel.RegisteredCount,
		Capacity:        sessionModel.Capacity,
		SessionLink:     sessionModel.SessionLink,
		ProfessionalId:  sessionModel.ProfessionalId,
		LocalId:         sessionModel.LocalId,
	}, nil
}

// Updates a session given fields in postgresql DB and returns it.
func (s *Session) UpdatePostgresqlSession(
	sessionId uuid.UUID,
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
) (*schemas.Session, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	// Call the DAO with individual parameters following the Local pattern
	sessionModel, err := s.DaoPostgresql.Session.UpdateSession(
		sessionId,
		title,
		date,
		startTime,
		endTime,
		state,
		capacity,
		sessionLink,
		professionalId,
		localId,
		updatedBy,
	)
	if err != nil {
		return nil, &errors.BadRequestError.SessionNotUpdated
	}

	return &schemas.Session{
		Id:              sessionModel.Id,
		Title:           sessionModel.Title,
		Date:            sessionModel.Date,
		StartTime:       sessionModel.StartTime,
		EndTime:         sessionModel.EndTime,
		State:           string(sessionModel.State),
		RegisteredCount: sessionModel.RegisteredCount,
		Capacity:        sessionModel.Capacity,
		SessionLink:     sessionModel.SessionLink,
		ProfessionalId:  sessionModel.ProfessionalId,
		LocalId:         sessionModel.LocalId,
	}, nil
}

// Soft deletes a session from postgresql DB.
func (s *Session) DeletePostgresqlSession(sessionId uuid.UUID) *errors.Error {
	err := s.DaoPostgresql.Session.DeleteSession(sessionId)
	if err != nil {
		return &errors.BadRequestError.SessionNotSoftDeleted
	}

	return nil
}

// Fetch sessions from postgresql DB and adapts them to Session schema.
func (s *Session) FetchPostgresqlSessions(
	professionalIds []uuid.UUID,
	localIds []uuid.UUID,
	states []string,
) ([]*schemas.Session, *errors.Error) {
	sessionModels, err := s.DaoPostgresql.Session.FetchSessions(professionalIds, localIds, states)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.SessionNotFound
	}

	sessions := make([]*schemas.Session, len(sessionModels))
	for i, sessionModel := range sessionModels {
		sessions[i] = &schemas.Session{
			Id:              sessionModel.Id,
			Title:           sessionModel.Title,
			Date:            sessionModel.Date,
			StartTime:       sessionModel.StartTime,
			EndTime:         sessionModel.EndTime,
			State:           string(sessionModel.State),
			RegisteredCount: sessionModel.RegisteredCount,
			Capacity:        sessionModel.Capacity,
			SessionLink:     sessionModel.SessionLink,
			ProfessionalId:  sessionModel.ProfessionalId,
			LocalId:         sessionModel.LocalId,
		}
	}

	return sessions, nil
}

// Creates multiple sessions into postgresql DB and returns them.
func (s *Session) BulkCreatePostgresqlSessions(
	sessionsData []*schemas.CreateSessionRequest,
	updatedBy string,
) ([]*schemas.Session, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	sessionsModel := make([]*model.Session, len(sessionsData))
	for i, sessionData := range sessionsData {
		sessionsModel[i] = &model.Session{
			Id:              uuid.New(),
			Title:           sessionData.Title,
			Date:            sessionData.Date,
			StartTime:       sessionData.StartTime,
			EndTime:         sessionData.EndTime,
			State:           model.SessionStateScheduled, // Default state
			RegisteredCount: 0,                           // Default number of initial registrations
			Capacity:        sessionData.Capacity,
			SessionLink:     sessionData.SessionLink,
			ProfessionalId:  sessionData.ProfessionalId,
			LocalId:         sessionData.LocalId,
			AuditFields: model.AuditFields{
				UpdatedBy: updatedBy,
			},
		}
	}

	if err := s.DaoPostgresql.Session.BulkCreateSessions(sessionsModel); err != nil {
		return nil, &errors.BadRequestError.SessionNotCreated
	}

	sessions := make([]*schemas.Session, len(sessionsModel))
	for i, sessionModel := range sessionsModel {
		sessions[i] = &schemas.Session{
			Id:              sessionModel.Id,
			Title:           sessionModel.Title,
			Date:            sessionModel.Date,
			StartTime:       sessionModel.StartTime,
			EndTime:         sessionModel.EndTime,
			State:           string(sessionModel.State),
			RegisteredCount: sessionModel.RegisteredCount,
			Capacity:        sessionModel.Capacity,
			SessionLink:     sessionModel.SessionLink,
			ProfessionalId:  sessionModel.ProfessionalId,
			LocalId:         sessionModel.LocalId,
		}
	}

	return sessions, nil
}
