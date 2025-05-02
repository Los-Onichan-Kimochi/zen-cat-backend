package controller

import (
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/gorm"
	"onichankimochi.com/zen_cat_backend/src/logging"
	"onichankimochi.com/zen_cat_backend/src/server/dao/zen_cat_psql/model"
	"onichankimochi.com/zen_cat_backend/src/server/schemas"
	"onichankimochi.com/zen_cat_backend/src/server/utils/psql"
)

type ZenCatPsqlCollection struct {
	Logger    logging.Logger
	Community *Community
}

// Create dao controller collection
func NewZenCatPsqlCollection(
	logger logging.Logger,
	envSettings *schemas.EnvSettings,
) (*ZenCatPsqlCollection, *gorm.DB) {
	postgresqlDB, err := psql.CreateConnection(
		envSettings.ZenCatPostgresHost,
		envSettings.ZenCatPostgresUser,
		envSettings.ZenCatPostgresPassword,
		envSettings.ZenCatPostgresName,
		envSettings.ZenCatPostgresPort,
		envSettings.EnableSqlLogs,
	)
	if err != nil {
		logger.Panicln("Failed to connect to ZenCat Postgresql database")
	}

	if err := postgresqlDB.Use(otelgorm.NewPlugin()); err != nil {
		logger.Panicln("Failed to instrument ZenCat Postgresql database")
	}

	// Create Community table
	if err := postgresqlDB.AutoMigrate(&model.Community{}); err != nil {
		panic(err)
	}

	return &ZenCatPsqlCollection{
		Logger:    logger,
		Community: NewCommunityController(logger, postgresqlDB),
	}, postgresqlDB
}
