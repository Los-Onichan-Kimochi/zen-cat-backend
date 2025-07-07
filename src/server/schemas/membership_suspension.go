package schemas

import (
	"time"

	"github.com/google/uuid"
)

type MembershipSuspension struct {
	Id           uuid.UUID  `json:"id"`
	MembershipId uuid.UUID  `json:"membership_id"`
	SuspendedAt  time.Time  `json:"suspended_at"`
	ResumedAt    *time.Time `json:"resumed_at,omitempty"`
}
