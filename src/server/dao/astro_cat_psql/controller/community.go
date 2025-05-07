package controller

import (
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
)

type Community struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

// Create Community postgresql controller
func NewCommunityController(logger logging.Logger, postgresqlDB *gorm.DB) *Community {
	return &Community{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}
