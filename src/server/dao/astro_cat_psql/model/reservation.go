package model

import (
	"time"

	"github.com/google/uuid"
)

type ReservationState string

const (
	ReservationStateDone      ReservationState = "DONE"
	ReservationStateConfirmed ReservationState = "CONFIRMED"
	ReservationStateCancelled ReservationState = "CANCELLED"
	ReservationStateAnulled   ReservationState = "ANULLED"
)

type Reservation struct {
	Id               uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name             string
	ReservationTime  time.Time
	State            ReservationState
	LastModification time.Time
	AuditFields

	UserId       uuid.UUID   `gorm:"type:uuid"`
	User         User        `gorm:"foreignKey:UserId"`
	SessionId    uuid.UUID   `gorm:"type:uuid"`
	Session      Session     `gorm:"foreignKey:SessionId"`
	MembershipId *uuid.UUID  `gorm:"type:uuid"`
	Membership   *Membership `gorm:"foreignKey:MembershipId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (Reservation) TableName() string {
	return "astro_cat_reservation"
}
