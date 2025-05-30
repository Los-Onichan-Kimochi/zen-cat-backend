package schemas

import "github.com/google/uuid"

type ServiceProfessional struct {
	Id             uuid.UUID `json:"id"`
	ServiceId      uuid.UUID `json:"service_id"`
	ProfessionalId uuid.UUID `json:"professional_id"`
}

type ServiceProfessionals struct {
	ServiceProfessionals []*ServiceProfessional `json:"service_professionals"`
}

type CreateServiceProfessionalRequest struct {
	ServiceId      uuid.UUID `json:"service_id" validate:"required"`
	ProfessionalId uuid.UUID `json:"professional_id"   validate:"required"`
}

type BatchCreateServiceProfessionalRequest struct {
	ServiceProfessionals []*CreateServiceProfessionalRequest `json:"service_professionals"`
}

type DeleteServiceProfessionalRequest struct {
	ServiceId      uuid.UUID `json:"service_id" validate:"required"`
	ProfessionalId uuid.UUID `json:"professional_id"      validate:"required"`
}

type BulkDeleteServiceProfessionalRequest struct {
	ServiceProfessionals []*DeleteServiceProfessionalRequest `json:"service_professionals"`
}
