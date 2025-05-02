package model

import (
	"github.com/google/uuid"
)

type Community struct {
	Id          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string
	Description string
	ImageUrl    string
	// Add Audit fields
}

func (Community) TableName() string {
	return "zen_cat_community"
}
