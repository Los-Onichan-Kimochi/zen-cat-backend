package adapter

import (
	daoPostgresql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"

	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type AdapterCollection struct {
	Logger        logging.Logger
	Community     *Community
	Professional  *Professional
	Local         *Local
	User          *User
	Service       *Service
	Plan          *Plan
	CommunityPlan *CommunityPlan
}

// Create bll adapter collection
func NewAdapterCollection(
	logger logging.Logger,
	envSettings *schemas.EnvSettings,
) (*AdapterCollection, *gorm.DB) {
	daoAstroCatPsql, astroCatPsqlDB := daoPostgresql.NewAstroCatPsqlCollection(logger, envSettings)

	return &AdapterCollection{
		Community:     NewCommunityAdapter(logger, daoAstroCatPsql),
		Professional:  NewProfessionalAdapter(logger, daoAstroCatPsql),
		Local:         NewLocalAdapter(logger, daoAstroCatPsql),
		User:          NewUserAdapter(logger, daoAstroCatPsql),
		Service:       NewServiceAdapter(logger, daoAstroCatPsql),
		Plan:          NewPlanAdapter(logger, daoAstroCatPsql),
		CommunityPlan: NewCommunityPlanAdapter(logger, daoAstroCatPsql),
	}, astroCatPsqlDB
}
