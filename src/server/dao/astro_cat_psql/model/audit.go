package model

import (
	"time"

	"gorm.io/gorm"
)

type AuditFields struct {
	CreatedAt time.Time      `gorm:"autoCreateTime"` // Creation date
	UpdatedAt time.Time      `gorm:"autoUpdateTime"` // Last update date
	DeletedAt gorm.DeletedAt `gorm:"index"`          // Delete date (soft delete)
	UpdatedBy string         `gorm:"size:255"`       // Admin user who updated the record
}
