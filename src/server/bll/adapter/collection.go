package adapter

import (
	daoPostgresql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"

	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type AdapterCollection struct {
	Logger    logging.Logger
	Community *Community
	Service   *Service
}

// Create bll adapter collection
func NewAdapterCollection(
	logger logging.Logger,
	envSettings *schemas.EnvSettings,
) (*AdapterCollection, *gorm.DB) {
	daoAstroCatPsql, astroCatPsqlDB := daoPostgresql.NewAstroCatPsqlCollection(logger, envSettings)

	return &AdapterCollection{
		Community: NewCommunityAdapter(logger, daoAstroCatPsql),
		Service:  NewServiceAdapter(logger, daoAstroCatPsql),
	}, astroCatPsqlDB
}
