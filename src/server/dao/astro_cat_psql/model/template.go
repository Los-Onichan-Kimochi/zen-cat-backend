package model

import "github.com/google/uuid"

type Template struct {
	Id   uuid.UUID `gorm:"type:uuid;primaryKey"`
	Link string
	AuditFields
}

func (Template) TableName() string {
	return "astro_cat_template"
}
