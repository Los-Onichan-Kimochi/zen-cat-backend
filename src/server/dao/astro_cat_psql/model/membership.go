package model

import (
	"time"

	"github.com/google/uuid"
)

type MembershipStatus string

const (
	MembershipStatusActive    MembershipStatus = "ACTIVE"
	MembershipStatusExpired   MembershipStatus = "EXPIRED"
	MembershipStatusCancelled MembershipStatus = "CANCELLED"
	MembershipStatusOnHold    MembershipStatus = "ON_HOLD"
)

type Membership struct {
	Id               uuid.UUID `gorm:"type:uuid;primaryKey"`
	Description      string
	StartDate        time.Time
	EndDate          time.Time
	Status           MembershipStatus
	ReservationsUsed *int
	AuditFields

	CommunityId uuid.UUID `gorm:"type:uuid"`
	Community   Community `gorm:"foreignKey:CommunityId;constraint:OnUpdate:CASCADE;"`
	UserId      uuid.UUID `gorm:"type:uuid"`
	User        User      `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE;"`
	PlanId      uuid.UUID `gorm:"type:uuid"`
	Plan        Plan      `gorm:"foreignKey:PlanId;constraint:OnUpdate:CASCADE;"`
}

func (Membership) TableName() string {
	return "astro_cat_membership"
}
