package model

import "github.com/google/uuid"

type ServiceLocal struct {
	ServiceId 	uuid.UUID `gorm:"primaryKey"`
	Service   	Service   `gorm:"foreignKey:ServiceId;references:Id"`
	LocalId     uuid.UUID `gorm:"primaryKey"`
	Local       Local     `gorm:"foreignKey:LocalId;references:Id"`
	AuditFields
}

func (ServiceLocal) TableName() string {
	return "astro_cat_service_local"
}