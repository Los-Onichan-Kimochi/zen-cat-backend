package schemas

import "github.com/google/uuid"

type ServiceLocal struct {
	ServiceId uuid.UUID `json:"service_id"`
	LocalId   uuid.UUID `json:"local_id"`
}

type CreateServiceLocalRequest struct {
	ServiceId uuid.UUID `json:"service_id" validate:"required"`
	LocalId   uuid.UUID `json:"local_id"   validate:"required"`
}
