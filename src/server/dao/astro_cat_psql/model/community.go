package model

import (
	"github.com/google/uuid"
)

type Community struct {
	Id                  uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name                string
	Purpose             string
	ImageUrl            string
	NumberSubscriptions int
	AuditFields

	Services []Service `gorm:"many2many:astro_cat_community_service;"`
}

func (Community) TableName() string {
	return "astro_cat_community"
}
