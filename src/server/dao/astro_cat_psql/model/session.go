package model

import (
	"time"

	"github.com/google/uuid"
)

type SessionState string

const (
	SessionStateScheduled   SessionState = "SCHEDULED"
	SessionStateOnGoing     SessionState = "ONGOING"
	SessionStateCompleted   SessionState = "COMPLETED"
	SessionStateCancelled   SessionState = "CANCELLED"
	SessionStateRescheduled SessionState = "RESCHEDULED"
)

type Session struct {
	Id              uuid.UUID `gorm:"type:uuid;primaryKey"`
	Title           string
	Date            time.Time
	StartTime       time.Time
	EndTime         time.Time
	State           SessionState
	RegisteredCount int
	Capacity        int
	SessionLink     *string
	AuditFields

	ProfessionalId uuid.UUID  `gorm:"type:uuid;foreignKey:ProfessionalId"`
	LocalId        *uuid.UUID `gorm:"type:uuid;foreignKey:LocalId"`
}

func (Session) TableName() string {
	return "astro_cat_session"
}
