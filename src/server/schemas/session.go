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
	CommunityServiceId *uuid.UUID `json:"community_service_id"`
}

type Sessions struct {
	Sessions []*Session `json:"sessions"`
}

type CreateSessionRequest struct {
	Title             string     `json:"title"`
	Date              time.Time  `json:"date"`
	StartTime         time.Time  `json:"start_time"`
	EndTime           time.Time  `json:"end_time"`
	Capacity          int        `json:"capacity"`
	SessionLink       *string    `json:"session_link"`
	ProfessionalId    uuid.UUID  `json:"professional_id"`
	LocalId           *uuid.UUID `json:"local_id"`
	CommunityServiceId *uuid.UUID `json:"community_service_id"`
}

type UpdateSessionRequest struct {
	Title           *string    `json:"title"`
	Date            *time.Time `json:"date"`
	StartTime       *time.Time `json:"start_time"`
	EndTime         *time.Time `json:"end_time"`
	State           *string    `json:"state"`
	RegisteredCount *int       `json:"registered_count"`
	Capacity        *int       `json:"capacity"`
	SessionLink     *string    `json:"session_link"`
	ProfessionalId  *uuid.UUID `json:"professional_id"`
	LocalId         *uuid.UUID `json:"local_id"`
	CommunityServiceId *uuid.UUID `json:"community_service_id"`
}

type BatchCreateSessionRequest struct {
	Sessions []*CreateSessionRequest `json:"sessions"`
}

type BulkDeleteSessionRequest struct {
	Sessions []uuid.UUID `json:"sessions"`
}

type CheckConflictRequest struct {
	Date           time.Time  `json:"date"`
	StartTime      time.Time  `json:"start_time"`
	EndTime        time.Time  `json:"end_time"`
	ProfessionalId uuid.UUID  `json:"professional_id"`
	LocalId        *uuid.UUID `json:"local_id"`
	ExcludeId      *uuid.UUID `json:"exclude_id"` // Para excluir sesión en modo edición
}

type ConflictResult struct {
	HasConflict           bool       `json:"has_conflict"`
	ProfessionalConflicts []*Session `json:"professional_conflicts"`
	LocalConflicts        []*Session `json:"local_conflicts"`
}

type AvailabilityRequest struct {
	Date           time.Time  `json:"date"`
	ProfessionalId *uuid.UUID `json:"professional_id"`
	LocalId        *uuid.UUID `json:"local_id"`
}

type TimeSlot struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Title string `json:"title"`
	Type  string `json:"type"` // "professional" | "local"
}

type AvailabilityResult struct {
	IsAvailable bool       `json:"is_available"`
	BusySlots   []TimeSlot `json:"busy_slots"`
}
