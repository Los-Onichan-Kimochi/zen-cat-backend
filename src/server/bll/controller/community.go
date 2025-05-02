package controller

import (
	"onichankimochi.com/zen_cat_backend/src/logging"
	bllAdapter "onichankimochi.com/zen_cat_backend/src/server/bll/adapter"
	"onichankimochi.com/zen_cat_backend/src/server/schemas"
)

type Community struct {
	logger      logging.Logger
	Adapter     *bllAdapter.AdapterCollection
	EnvSettings *schemas.EnvSettings
}

// Create Community controller
func NewCommunityController(
	logger logging.Logger,
	adapter *bllAdapter.AdapterCollection,
	envSettings *schemas.EnvSettings,
) *Community {
	return &Community{
		logger:      logger,
		Adapter:     adapter,
		EnvSettings: envSettings,
	}
}
