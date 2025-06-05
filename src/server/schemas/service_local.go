package schemas

import "github.com/google/uuid"

type ServiceLocal struct {
	Id        uuid.UUID `json:"id"`
	ServiceId uuid.UUID `json:"service_id"`
	LocalId   uuid.UUID `json:"local_id"`
}

type ServiceLocals struct {
	ServiceLocals []*ServiceLocal `json:"service_locals"`
}

type CreateServiceLocalRequest struct {
	ServiceId uuid.UUID `json:"service_id" validate:"required"`
	LocalId   uuid.UUID `json:"local_id"   validate:"required"`
}

type BatchCreateServiceLocalRequest struct {
	ServiceLocals []*CreateServiceLocalRequest `json:"service_locals"`
}

type DeleteServiceLocalRequest struct {
	ServiceId uuid.UUID `json:"service_id" validate:"required"`
	LocalId      uuid.UUID `json:"local_id"      validate:"required"`
}

type BulkDeleteServiceLocalRequest struct {
	ServiceLocals []*DeleteServiceLocalRequest `json:"service_locals"`
}