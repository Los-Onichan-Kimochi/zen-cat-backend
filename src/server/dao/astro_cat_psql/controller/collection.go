package controller

import (
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	"onichankimochi.com/astro_cat_backend/src/server/utils/psql"
)

type AstroCatPsqlCollection struct {
	Logger    logging.Logger
	Community *Community
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

	// Create Community table
	if err := postgresqlDB.AutoMigrate(&model.Community{}); err != nil {
		panic(err)
	}

	return &AstroCatPsqlCollection{
		Logger:    logger,
		Community: NewCommunityController(logger, postgresqlDB),
	}, postgresqlDB
}
