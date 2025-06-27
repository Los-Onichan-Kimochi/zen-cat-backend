package schemas

import (
	"time"

	"github.com/google/uuid"
)

type DocumentType string

const (
	DocumentTypeDNI           DocumentType = "DNI"
	DocumentTypeForeignerCard DocumentType = "FOREIGNER_CARD"
	DocumentTypePassport      DocumentType = "PASSPORT"
)

type Gender string

const (
	GenderMale   Gender = "MALE"
	GenderFemale Gender = "FEMALE"
	GenderOther  Gender = "OTHER"
)

type IdentificationDocument struct {
	Id             uuid.UUID `json:"id"`
	DocumentType   string    `json:"document_type"`
	DocumentNumber string    `json:"document_number"`
}

type Onboarding struct {
	Id uuid.UUID `json:"id"`
	// Documento
	DocumentType   DocumentType `json:"document_type"`
	DocumentNumber string       `json:"document_number"`
	// Contacto
	PhoneNumber string `json:"phone_number"`
	// Datos personales adicionales
	BirthDate *time.Time `json:"birth_date"`
	Gender    *Gender    `json:"gender"`
	// Dirección
	PostalCode string  `json:"postal_code"`
	Address    string  `json:"address"`
	District   *string `json:"district"`
	Province   *string `json:"province"`
	Region     *string `json:"region"`

	UserId uuid.UUID `json:"user_id"`
}

type Onboardings struct {
	Onboardings []*Onboarding `json:"onboardings"`
	// TODO: Add pagination if needed
}

type CreateOnboardingRequest struct {
	// Documento
	DocumentType   DocumentType `json:"document_type" binding:"required"`
	DocumentNumber string       `json:"document_number" binding:"required"`
	// Contacto
	PhoneNumber string `json:"phone_number" binding:"required"`
	// Datos personales adicionales
	BirthDate *time.Time `json:"birth_date"`
	Gender    *Gender    `json:"gender"`
	// Dirección
	PostalCode string  `json:"postal_code" binding:"required"`
	Address    string  `json:"address" binding:"required"`
	District   *string `json:"district"`
	Province   *string `json:"province"`
	Region     *string `json:"region"`
}

type UpdateOnboardingRequest struct {
	// Documento
	DocumentType   *DocumentType `json:"document_type"`
	DocumentNumber *string       `json:"document_number"`
	// Contacto
	PhoneNumber *string `json:"phone_number"`
	// Datos personales adicionales
	BirthDate *time.Time `json:"birth_date"`
	Gender    *Gender    `json:"gender"`
	// Dirección
	PostalCode *string `json:"postal_code"`
	Address    *string `json:"address"`
	District   *string `json:"district"`
	Province   *string `json:"province"`
	Region     *string `json:"region"`
}
