package model

import "github.com/google/uuid"

type Local struct {
	Id             uuid.UUID `gorm:"type:uuid;primaryKey"`
	LocalName      string
	StreetName     string
	BuildingNumber string
	District       string
	Province       string
	Region         string
	Reference      string
	Capacity       int
	ImageUrl       string
	AuditFields
}

func (Local) TableName() string {
	return "astro_cat_local"
}
