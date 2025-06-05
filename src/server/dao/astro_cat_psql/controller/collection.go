package controller

import (
	// para agregar data dummy

	//-----------------------
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	"onichankimochi.com/astro_cat_backend/src/server/utils/psql"
)

type AstroCatPsqlCollection struct {
	Logger           logging.Logger
	Community        *Community
	Professional     *Professional
	Local            *Local
	User             *User
	Service          *Service
	Plan             *Plan
	CommunityPlan    *CommunityPlan
	CommunityService *CommunityService
	ServiceLocal      *ServiceLocal
	ServiceProfessional *ServiceProfessional
	Session          *Session
	Reservation      *Reservation
}

// Create dao controller collection
func NewAstroCatPsqlCollection(
	logger logging.Logger,
	envSettings *schemas.EnvSettings,
) (*AstroCatPsqlCollection, *gorm.DB) {
	postgresqlDB, err := psql.CreateConnection(
		envSettings.AstroCatPostgresHost,
		envSettings.AstroCatPostgresUser,
		envSettings.AstroCatPostgresPassword,
		envSettings.AstroCatPostgresName,
		envSettings.AstroCatPostgresPort,
		envSettings.EnableSqlLogs,
	)
	if err != nil {
		logger.Panicln("Failed to connect to AstroCat Postgresql database")
	}

	if err := postgresqlDB.Use(otelgorm.NewPlugin()); err != nil {
		logger.Panicln("Failed to instrument AstroCat Postgresql database")
	}

	createTables(postgresqlDB)

	return &AstroCatPsqlCollection{
		Logger:           logger,
		Community:        NewCommunityController(logger, postgresqlDB),
		Professional:     NewProfessionalController(logger, postgresqlDB),
		Local:            NewLocalController(logger, postgresqlDB),
		User:             NewUserController(logger, postgresqlDB),
		Service:          NewServiceController(logger, postgresqlDB),
		Plan:             NewPlanController(logger, postgresqlDB),
		CommunityPlan:    NewCommunityPlanController(logger, postgresqlDB),
		CommunityService: NewCommunityServiceController(logger, postgresqlDB),
		ServiceLocal: 	  NewServiceLocalController(logger, postgresqlDB),
		ServiceProfessional: NewServiceProfessionalController(logger, postgresqlDB),
		Session:          NewSessionController(logger, postgresqlDB),
		Reservation:      NewReservationController(logger, postgresqlDB),
	}, postgresqlDB
}

// Helper function to create AstroCat tables
func createTables(astroCatPsqlDB *gorm.DB) {
	if err := astroCatPsqlDB.AutoMigrate(&model.Plan{}); err != nil {
		panic(err)
	}
	if err := astroCatPsqlDB.AutoMigrate(&model.Template{}); err != nil {
		panic(err)
	}
	if err := astroCatPsqlDB.AutoMigrate(&model.Local{}); err != nil {
		panic(err)
	}
	if err := astroCatPsqlDB.AutoMigrate(&model.Professional{}); err != nil {
		panic(err)
	}
	if err := astroCatPsqlDB.AutoMigrate(&model.Onboarding{}); err != nil {
		panic(err)
	}
	if err := astroCatPsqlDB.AutoMigrate(&model.User{}); err != nil {
		panic(err)
	}
	if err := astroCatPsqlDB.AutoMigrate(&model.Community{}); err != nil {
		panic(err)
	}
	if err := astroCatPsqlDB.AutoMigrate(&model.Membership{}); err != nil {
		panic(err)
	}
	if err := astroCatPsqlDB.AutoMigrate(&model.Service{}); err != nil {
		panic(err)
	}
	if err := astroCatPsqlDB.AutoMigrate(&model.Session{}); err != nil {
		panic(err)
	}
	if err := astroCatPsqlDB.AutoMigrate(&model.Reservation{}); err != nil {
		panic(err)
	}
	if err := astroCatPsqlDB.AutoMigrate(&model.CommunityService{}); err != nil {
		panic(err)
	}
	if err := astroCatPsqlDB.AutoMigrate(&model.CommunityPlan{}); err != nil {
		panic(err)
	}

	if err := astroCatPsqlDB.AutoMigrate(&model.ServiceLocal{}); err != nil {
		panic(err)
	}
	if err := astroCatPsqlDB.AutoMigrate(&model.ServiceProfessional{}); err != nil {
		panic(err)
	}
}
