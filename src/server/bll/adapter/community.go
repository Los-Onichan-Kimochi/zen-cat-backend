package adapter

import (
	daoPostgresql "onichankimochi.com/zen_cat_backend/src/server/dao/zen_cat_psql/controller"

	"onichankimochi.com/zen_cat_backend/src/logging"
)

type Community struct {
	logger        logging.Logger
	DaoPostgresql *daoPostgresql.ZenCatPsqlCollection
}

// Creates Community adapter
func NewCommunityAdapter(
	logger logging.Logger,
	daoPostgresql *daoPostgresql.ZenCatPsqlCollection,
) *Community {
	return &Community{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

// Postgresql Functions
