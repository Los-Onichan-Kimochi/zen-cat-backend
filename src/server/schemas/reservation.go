package schemas

import (
	"time"

	"github.com/google/uuid"
)

type Reservation struct {
	Id               uuid.UUID  `json:"id"`
	Name             string     `json:"name"`
	ReservationTime  time.Time  `json:"reservation_time"`
	State            string     `json:"state"`
	LastModification time.Time  `json:"last_modification"`
	UserId           uuid.UUID  `json:"user_id"`
	SessionId        uuid.UUID  `json:"session_id"`
	Session          Session    `json:"session"`
	MembershipId     *uuid.UUID `json:"membership_id,omitempty"`
}

type Reservations struct {
	Reservations []*Reservation `json:"reservations"`
}

type CreateReservationRequest struct {
	Name            string     `json:"name"`
	ReservationTime time.Time  `json:"reservation_time"`
	State           string     `json:"state"`
	UserId          uuid.UUID  `json:"user_id"`
	SessionId       uuid.UUID  `json:"session_id"`
	MembershipId    *uuid.UUID `json:"membership_id,omitempty"`
}

type UpdateReservationRequest struct {
	Name            *string    `json:"name"`
	ReservationTime *time.Time `json:"reservation_time"`
	State           *string    `json:"state"`
	UserId          *uuid.UUID `json:"user_id"`
	SessionId       *uuid.UUID `json:"session_id"`
	MembershipId    *uuid.UUID `json:"membership_id,omitempty"`
}

type BulkDeleteReservationRequest struct {
	Reservations []string `json:"reservations"`
}
