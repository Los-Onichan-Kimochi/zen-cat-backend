package model

import "github.com/google/uuid"

type ServiceLocal struct {
	Id           uuid.UUID `gorm:"type:uuid;primaryKey"`
	ServiceId    uuid.UUID `gorm:"type:uuid"`
	Service      Service   `gorm:"foreignKey:ServiceId;references:Id"`
	LocalId      uuid.UUID `gorm:"type:uuid"`
	Local        Local     `gorm:"foreignKey:LocalId;references:Id"`
	AuditFields
}

func (ServiceLocal) TableName() string {
	return "astro_cat_service_local"
}
