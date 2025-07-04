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
	Id               uuid.UUID        `json:"id"`
	Description      string           `json:"description"`
	StartDate        time.Time        `json:"start_date"`
	EndDate          time.Time        `json:"end_date"`
	Status           MembershipStatus `json:"status"`
	ReservationsUsed *int             `json:"reservations_used"`
	CommunityId      uuid.UUID        `json:"community_id"`
	Community        Community        `json:"community"`
	UserId           uuid.UUID        `json:"user_id"`
	User             User             `json:"user"`
	PlanId           uuid.UUID        `json:"plan_id"`
	Plan             Plan             `json:"plan"`
}

type Memberships struct {
	Memberships []*Membership `json:"memberships"`
}

type CreateMembershipRequest struct {
	Description      string           `json:"description"`
	StartDate        time.Time        `json:"start_date"`
	EndDate          time.Time        `json:"end_date"`
	Status           MembershipStatus `json:"status"`
	ReservationsUsed *int             `json:"reservations_used"`
	CommunityId      uuid.UUID        `json:"community_id"`
	UserId           uuid.UUID        `json:"user_id"`
	PlanId           uuid.UUID        `json:"plan_id"`
}

type CreateMembershipForUserRequest struct {
	Description      string           `json:"description"`
	StartDate        time.Time        `json:"start_date"`
	EndDate          time.Time        `json:"end_date"`
	Status           MembershipStatus `json:"status"`
	ReservationsUsed *int             `json:"reservations_used"`
	CommunityId      uuid.UUID        `json:"community_id"`
	PlanId           uuid.UUID        `json:"plan_id"`
}

type UpdateMembershipRequest struct {
	Description      *string           `json:"description"`
	StartDate        *time.Time        `json:"start_date"`
	EndDate          *time.Time        `json:"end_date"`
	Status           *MembershipStatus `json:"status"`
	ReservationsUsed *int              `json:"reservations_used"`
	CommunityId      *uuid.UUID        `json:"community_id"`
	UserId           *uuid.UUID        `json:"user_id"`
	PlanId           *uuid.UUID        `json:"plan_id"`
}
