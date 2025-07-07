package model

import (
	"time"

	"github.com/google/uuid"
)

type MembershipSuspension struct {
	Id          uuid.UUID `gorm:"type:uuid;primaryKey"`
	SuspendedAt time.Time
	ResumedAt   *time.Time // Pointer to allow NULL values
	AuditFields

	MembershipId uuid.UUID  `gorm:"type:uuid;not null"`
	Membership   Membership `gorm:"foreignKey:MembershipId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (MembershipSuspension) TableName() string {
	return "astro_cat_membership_suspension"
}
