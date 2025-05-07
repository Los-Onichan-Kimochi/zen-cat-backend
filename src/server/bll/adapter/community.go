package adapter

import (
	daoPostgresql "onichankimochi.com/astro_cat_backend/src/server/dao/zen_cat_psql/controller"

	"onichankimochi.com/astro_cat_backend/src/logging"
)

type Community struct {
	logger        logging.Logger
	DaoPostgresql *daoPostgresql.AstroCatPsqlCollection
}

// Creates Community adapter
func NewCommunityAdapter(
	logger logging.Logger,
	daoPostgresql *daoPostgresql.AstroCatPsqlCollection,
) *Community {
	return &Community{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

// Postgresql Functions
