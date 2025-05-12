package model

import "github.com/google/uuid"

type DocumentType string

const (
	DocumentTypeDni           DocumentType = "DNI"
	DocumentTypeForeignerCard DocumentType = "FOREIGNER_CARD"
)

type Onboarding struct {
	Id             uuid.UUID `gorm:"type:uuid;primaryKey"`
	PhoneNumber    string
	DocumentType   DocumentType
	DocumentNumber string
	StreetName     string
	BuildingNumber string
	District       string
	Province       string
	Region         string
	Reference      string
	AuditFields

	UserId uuid.UUID `gorm:"type:uuid;unique"`
	User   User      `gorm:"foreignKey:UserId"`
}

func (Onboarding) TableName() string {
	return "astro_cat_onboarding"
}
