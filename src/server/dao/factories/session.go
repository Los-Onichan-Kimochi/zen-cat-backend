package factories

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type SessionModelF struct {
	Id              *uuid.UUID
	Title           *string
	Date            *time.Time
	StartTime       *time.Time
	EndTime         *time.Time
	State           *model.SessionState
	RegisteredCount *int
	Capacity        *int
	SessionLink     *string
	ProfessionalId  *uuid.UUID
	LocalId         *uuid.UUID
	CommunityServiceId *uuid.UUID
}

// Create a new session on DB
func NewSessionModel(db *gorm.DB, option ...SessionModelF) *model.Session {
	// Create default professional if not provided
	professional := NewProfessionalModel(db)

	// Create default local if not provided
	local := NewLocalModel(db)
	
	// Create default community service if not provided
	communityService := NewCommunityServiceModel(db)

	now := time.Now()
	startTime := now.Add(1 * time.Hour)
	endTime := now.Add(2 * time.Hour)
	registeredCount := 0
	capacity := 10
	sessionLink := "https://meet.example.com/session"

	session := &model.Session{
		Id:              uuid.New(),
		Title:           "Test Session",
		Date:            now,
		StartTime:       startTime,
		EndTime:         endTime,
		State:           model.SessionStateScheduled,
		RegisteredCount: registeredCount,
		Capacity:        capacity,
		SessionLink:     &sessionLink,
		ProfessionalId:  professional.Id,
		LocalId:         &local.Id,
		CommunityServiceId: &communityService.Id,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}

	if len(option) > 0 {
		parameters := option[0]
		if parameters.Id != nil {
			session.Id = *parameters.Id
		}
		if parameters.Title != nil {
			session.Title = *parameters.Title
		}
		if parameters.Date != nil {
			session.Date = *parameters.Date
		}
		if parameters.StartTime != nil {
			session.StartTime = *parameters.StartTime
		}
		if parameters.EndTime != nil {
			session.EndTime = *parameters.EndTime
		}
		if parameters.State != nil {
			session.State = *parameters.State
		}
		if parameters.RegisteredCount != nil {
			session.RegisteredCount = *parameters.RegisteredCount
		}
		if parameters.Capacity != nil {
			session.Capacity = *parameters.Capacity
		}
		if parameters.SessionLink != nil {
			session.SessionLink = parameters.SessionLink
		}
		if parameters.ProfessionalId != nil {
			session.ProfessionalId = *parameters.ProfessionalId
		}
		if parameters.LocalId != nil {
			session.LocalId = parameters.LocalId
		}
		if parameters.CommunityServiceId != nil {
			session.CommunityServiceId = parameters.CommunityServiceId
		}
	}

	result := db.Create(session)
	if result.Error != nil {
		log.Fatalf("Error when trying to create session: %v", result.Error)
	}

	return session
}

// Create size number of new sessions on DB
func NewSessionModelBatch(
	db *gorm.DB,
	size int,
	option ...SessionModelF,
) []*model.Session {
	sessions := []*model.Session{}
	for i := 0; i < size; i++ {
		var session *model.Session
		if len(option) > 0 {
			session = NewSessionModel(db, option[0])
		} else {
			session = NewSessionModel(db)
		}
		sessions = append(sessions, session)
	}
	return sessions
}
