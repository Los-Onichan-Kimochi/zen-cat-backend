package model

import "github.com/google/uuid"

type ProfessionalType string

const (
	ProfessionalTypeMedic       ProfessionalType = "MEDIC"
	ProfessionalTypeGymTrainer  ProfessionalType = "GYM_TRAINER"
	ProfessionalTypeYogaTrainer ProfessionalType = "YOGA_TRAINER"
)

type Professional struct {
	Id             uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name           string
	FirstLastName  string
	SecondLastName *string
	Specialty      string
	Email          string
	PhoneNumber    string
	Type           ProfessionalType
	ImageUrl       string
	AuditFields

	Template *Template `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Professional) TableName() string {
	return "astro_cat_professional"
}
