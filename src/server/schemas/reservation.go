package schemas

import (
	"time"

	"github.com/google/uuid"
)

type Reservation struct {
	Id               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	ReservationTime  time.Time `json:"reservation_time"`
	State            string    `json:"state"`
	LastModification time.Time `json:"last_modification"`
	UserId           uuid.UUID `json:"user_id"`
	SessionId        uuid.UUID `json:"session_id"`
}

type Reservations struct {
	Reservations []*Reservation `json:"reservations"`
}
