package model

import "github.com/google/uuid"

type ServiceProfessional struct {
	ServiceId 	uuid.UUID `gorm:"primaryKey"`
	Service   	Service   `gorm:"foreignKey:ServiceId;references:Id"`
	ProfessionalId     uuid.UUID        `gorm:"primaryKey"`
	Professional       Professional     `gorm:"foreignKey:ProfessionalId;references:Id"`
	AuditFields
}

func (ServiceProfessional) TableName() string {
	return "astro_cat_service_professional"
}