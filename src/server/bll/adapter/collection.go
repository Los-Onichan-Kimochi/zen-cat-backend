package adapter

import (
	daoPostgresql "onichankimochi.com/zen_cat_backend/src/server/dao/zen_cat_psql/controller"

	"gorm.io/gorm"
	"onichankimochi.com/zen_cat_backend/src/logging"
	"onichankimochi.com/zen_cat_backend/src/server/schemas"
)

type AdapterCollection struct {
	Logger    logging.Logger
	Community *Community
}

// Create bll adapter collection
func NewAdapterCollection(
	logger logging.Logger,
	envSettings *schemas.EnvSettings,
) (*AdapterCollection, *gorm.DB) {
	daoZenCatPsql, zenCatPsqlDB := daoPostgresql.NewZenCatPsqlCollection(logger, envSettings)

	return &AdapterCollection{
		Community: NewCommunityAdapter(logger, daoZenCatPsql),
	}, zenCatPsqlDB
}
