package model

import "github.com/google/uuid"

type ServiceProfessional struct {
	Id                  uuid.UUID `gorm:"type:uuid;primaryKey"`
	ServiceId           uuid.UUID `gorm:"type:uuid"`
	Service             Service   `gorm:"foreignKey:ServiceId;references:Id"`
	ProfessionalId      uuid.UUID `gorm:"type:uuid"`
	Professional        Professional     `gorm:"foreignKey:ProfessionalId;references:Id"`
	AuditFields
}

func (ServiceProfessional) TableName() string {
	return "astro_cat_service_professional"
}