package model

import "github.com/google/uuid"

type Template struct {
	Id   uuid.UUID `gorm:"type:uuid;primaryKey"`
	Link string
	AuditFields

	ProfessionalId uuid.UUID    `gorm:"type:uuid;unique"`
	Professional   Professional `gorm:"foreignKey:ProfessionalId"`
}

func (Template) TableName() string {
	return "astro_cat_template"
}
