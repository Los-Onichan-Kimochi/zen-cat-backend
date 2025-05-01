package main

import (
	"onichankimochi.com/zen_cat_backend/src/logging"
	"onichankimochi.com/zen_cat_backend/src/server/schemas"
)

func main() {
	logger := logging.NewLogger("ZenCatBackendServer", "Version 1.0", logging.FormatText, 4)
	envSettings := schemas.NewEnvSettings(logger)

	logger.Infoln("ZenCat Backend Server running...")
	logger.Infoln("Main Port:", envSettings.MainPort)
	logger.Infoln("ZenCat Postgres Host:", envSettings.ZenCatPostgresHost)
	logger.Infoln("ZenCat Postgres Port:", envSettings.ZenCatPostgresPort)
	logger.Infoln("ZenCat Postgres User:", envSettings.ZenCatPostgresUser)
}
