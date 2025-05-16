package schemas

import (
	"time"

	"github.com/google/uuid"
)

type MembershipStatus string

const (
	MembershipStatusActive    MembershipStatus = "ACTIVE"
	MembershipStatusExpired   MembershipStatus = "EXPIRED"
	MembershipStatusCancelled MembershipStatus = "CANCELLED"
)

type Membership struct {
	Id          uuid.UUID        `json:"id"`
	Description string           `json:"description"`
	StartDate   time.Time        `json:"start_date"`
	EndDate     time.Time        `json:"end_date"`
	Status      MembershipStatus `json:"status"`
	Community   Community        `json:"community"`
	Plan        Plan             `json:"plan"`
}

type Memberships struct {
	Memberships []*Membership `json:"memberships"`
}

type CreateMembershipRequest struct {
	Description string           `json:"description"`
	StartDate   time.Time        `json:"start_date"`
	EndDate     time.Time        `json:"end_date"`
	Status      MembershipStatus `json:"status"`
}

type UpdateMembershipRequest struct {
	Description *string           `json:"description"`
	StartDate   *time.Time        `json:"start_date"`
	EndDate     *time.Time        `json:"end_date"`
	Status      *MembershipStatus `json:"status"`
}
