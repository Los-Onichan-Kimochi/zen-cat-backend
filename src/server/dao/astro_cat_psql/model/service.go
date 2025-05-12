package model

import "github.com/google/uuid"

type Service struct {
	Id          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string
	Description string
	ImageUrl    string
	IsVirtual   bool

	Locals        []Local        `gorm:"many2many:astro_cat_service_local;"`
	Professionals []Professional `gorm:"many2many:astro_cat_service_professional;"`
	AuditFields
}

func (Service) TableName() string {
	return "astro_cat_service"
}
