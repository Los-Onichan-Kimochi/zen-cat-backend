package schemas

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Id              uuid.UUID  `json:"id"`
	Title           string     `json:"title"`
	Date            time.Time  `json:"date"`
	StartTime       time.Time  `json:"start_time"`
	EndTime         time.Time  `json:"end_time"`
	State           string     `json:"state"`
	RegisteredCount int        `json:"registered_count"`
	Capacity        int        `json:"capacity"`
	SessionLink     *string    `json:"session_link"`
	ProfessionalId  uuid.UUID  `json:"professional_id"`
	LocalId         *uuid.UUID `json:"local_id"`
}

type Sessions struct {
	Sessions []*Session `json:"sessions"`
}

type CreateSessionRequest struct {
	Title          string     `json:"title"`
	Date           time.Time  `json:"date"`
	StartTime      time.Time  `json:"start_time"`
	EndTime        time.Time  `json:"end_time"`
	Capacity       int        `json:"capacity"`
	SessionLink    *string    `json:"session_link"`
	ProfessionalId uuid.UUID  `json:"professional_id"`
	LocalId        *uuid.UUID `json:"local_id"`
}

type UpdateSessionRequest struct {
	Title          *string    `json:"title"`
	Date           *time.Time `json:"date"`
	StartTime      *time.Time `json:"start_time"`
	EndTime        *time.Time `json:"end_time"`
	State          *string    `json:"state"`
	Capacity       *int       `json:"capacity"`
	SessionLink    *string    `json:"session_link"`
	ProfessionalId *uuid.UUID `json:"professional_id"`
	LocalId        *uuid.UUID `json:"local_id"`
}

type BatchCreateSessionRequest struct {
	Sessions []*CreateSessionRequest `json:"sessions"`
}
