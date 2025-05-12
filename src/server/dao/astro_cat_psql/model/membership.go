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
)

type Membership struct {
	Id          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Description string
	StartDate   time.Time
	EndDate     time.Time
	Status      MembershipStatus
	AuditFields

	CommunityId uuid.UUID `gorm:"type:uuid;foreignKey:CommunityId"`
	UserId      uuid.UUID `gorm:"type:uuid;foreignKey:UserId"`
	PlanId      uuid.UUID `gorm:"type:uuid;foreignKey:PlanId"`
}

func (Membership) TableName() string {
	return "astro_cat_membership"
}
