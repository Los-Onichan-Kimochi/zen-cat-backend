package schemas

import "github.com/google/uuid"

type ServiceProfessional struct {
	ServiceId 		 uuid.UUID `json:"service_id"`
	ProfessionalId   uuid.UUID `json:"professional_id"`
}

type CreateServiceProfessionalRequest struct {
	ServiceId 			uuid.UUID `json:"service_id" validate:"required"`
	ProfessionalId   	uuid.UUID `json:"professional_id"   validate:"required"`
}
