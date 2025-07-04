package main

import (
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/api"
	"onichankimochi.com/astro_cat_backend/src/server/config"
	daoPostgresql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"
	"onichankimochi.com/astro_cat_backend/src/server/scheduler"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

func main() {
	config.InitDevMode()
	logger := logging.NewLogger("AstroCatBackendServer", "Version 1.0", logging.FormatText, 4)
	envSettings := schemas.NewEnvSettings(logger)
	_, db := daoPostgresql.NewAstroCatPsqlCollection(logger, envSettings)
	envSettings.DB = db
	go scheduler.StartDailyReminderJob(envSettings)
	api.RunService(envSettings, logger)
}
