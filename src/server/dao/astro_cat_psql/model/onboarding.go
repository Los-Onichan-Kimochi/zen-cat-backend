package model

import (
	"time"

	"github.com/google/uuid"
)

type DocumentType string

const (
	DocumentTypeDni           DocumentType = "DNI"
	DocumentTypeForeignerCard DocumentType = "FOREIGNER_CARD"
	DocumentTypePassport      DocumentType = "PASSPORT"
)

type Gender string

const (
	GenderMale   Gender = "MALE"
	GenderFemale Gender = "FEMALE"
	GenderOther  Gender = "OTHER"
)

type Onboarding struct {
	Id uuid.UUID `gorm:"type:uuid;primaryKey"`
	// Documento
	DocumentType   DocumentType `gorm:"type:varchar(50);not null"`
	DocumentNumber string       `gorm:"type:varchar(20);not null"`
	// Contacto
	PhoneNumber string `gorm:"type:varchar(20);not null"`
	// Datos personales adicionales
	BirthDate *time.Time `gorm:"type:date"`
	Gender    *Gender    `gorm:"type:varchar(20)"`
	// Dirección
	City       string `gorm:"type:varchar(100);not null"`
	PostalCode string `gorm:"type:varchar(10);not null"`
	District   string `gorm:"type:varchar(100);not null"`
	Address    string `gorm:"type:varchar(255);not null"`
	AuditFields

	UserId uuid.UUID `gorm:"type:uuid;unique"`
	User   User      `gorm:"foreignKey:UserId"`
}

func (Onboarding) TableName() string {
	return "astro_cat_onboarding"
}
