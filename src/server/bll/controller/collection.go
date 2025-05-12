package controller

import (
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/bll/adapter"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type ControllerCollection struct {
	Logger       logging.Logger
	EnvSettings  *schemas.EnvSettings
	Community    *Community
	Professional *Professional
}

// Create bll controller collection
func NewControllerCollection(
	logger logging.Logger,
	envSettings *schemas.EnvSettings,
) (*ControllerCollection, *gorm.DB) {
	bllAdapter, astroCatPsqlDB := adapter.NewAdapterCollection(
		logger,
		envSettings,
	)
	community := NewCommunityController(logger, bllAdapter, envSettings)
	professional := NewProfessionalController(logger, bllAdapter, envSettings)

	return &ControllerCollection{
		Logger:       logger,
		EnvSettings:  envSettings,
		Community:    community,
		Professional: professional,
	}, astroCatPsqlDB
}
