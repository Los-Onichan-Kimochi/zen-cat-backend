package controller

import (
	"gorm.io/gorm"
	"onichankimochi.com/zen_cat_backend/src/logging"
	"onichankimochi.com/zen_cat_backend/src/server/bll/adapter"
	"onichankimochi.com/zen_cat_backend/src/server/schemas"
)

type ControllerCollection struct {
	Logger      logging.Logger
	EnvSettings *schemas.EnvSettings
	Community   *Community
}

// Create bll controller collection
func NewControllerCollection(
	logger logging.Logger,
	envSettings *schemas.EnvSettings,
) (*ControllerCollection, *gorm.DB) {
	bllAdapter, zenCatPsqlDB := adapter.NewAdapterCollection(
		logger,
		envSettings,
	)
	community := NewCommunityController(logger, bllAdapter, envSettings)

	return &ControllerCollection{
		Logger:      logger,
		EnvSettings: envSettings,
		Community:   community,
	}, zenCatPsqlDB
}
