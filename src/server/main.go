package main

import (
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/api"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

func main() {
	logger := logging.NewLogger("AstroCatBackendServer", "Version 1.0", logging.FormatText, 4)
	envSettings := schemas.NewEnvSettings(logger)

	api.RunService(envSettings, logger)
}
