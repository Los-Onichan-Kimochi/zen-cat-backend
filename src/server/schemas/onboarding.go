package schemas

import "github.com/google/uuid"

type IdentificationDocument struct {
	Id             uuid.UUID `json:"id"`
	DocumentType   string    `json:"document_type"`
	DocumentNumber string    `json:"document_number"`
}

type Onboarding struct {
	Id                     uuid.UUID              `json:"id"`
	PhoneNumber            string                 `json:"phone_number"`
	Address                string                 `json:"address"`
	District               string                 `json:"district"`
	City                   string                 `json:"city"`
	PostalCode             string                 `json:"postal_code"`
	IdentificationDocument IdentificationDocument `json:"identification_document"`
}

type Onboardings struct {
	Onboardings []*Onboarding `json:"onboardings"`
	// TODO: Add pagination if needed
}

type CreateOnboardingRequest struct {
	PhoneNumber            string                 `json:"phone_number"`
	Address                string                 `json:"address"`
	District               string                 `json:"district"`
	City                   string                 `json:"city"`
	PostalCode             string                 `json:"postal_code"`
	IdentificationDocument IdentificationDocument `json:"identification_document"`
}

type UpdateOnboardingRequest struct {
	PhoneNumber            *string                 `json:"phone_number"`
	Address                *string                 `json:"address"`
	District               *string                 `json:"district"`
	City                   *string                 `json:"city"`
	PostalCode             *string                 `json:"postal_code"`
	IdentificationDocument *IdentificationDocument `json:"identification_document"`
}
