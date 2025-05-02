package main

import (
	"onichankimochi.com/zen_cat_backend/src/logging"
	"onichankimochi.com/zen_cat_backend/src/server/api"
	"onichankimochi.com/zen_cat_backend/src/server/schemas"
)

func main() {
	logger := logging.NewLogger("ZenCatBackendServer", "Version 1.0", logging.FormatText, 4)
	envSettings := schemas.NewEnvSettings(logger)

	api.RunService(envSettings, logger)
}
