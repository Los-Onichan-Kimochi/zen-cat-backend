package main

import (
	"onichankimochi.com/astro_cat_backend/src/logging"
	daoPostgresql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	testsSetup "onichankimochi.com/astro_cat_backend/src/server/tests"
)

func main() {
	testLogger := logging.NewLoggerMock()
	envSettings := schemas.NewEnvSettings(testLogger)

	_, astroCatPsqlDB := daoPostgresql.NewAstroCatPsqlCollection(testLogger, envSettings)

	testsSetup.ClearPostgresqlDatabase(testLogger, astroCatPsqlDB, envSettings, nil)
}
