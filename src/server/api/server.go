package api

import (
	"github.com/labstack/echo/v4"
	"onichankimochi.com/zen_cat_backend/src/logging"
	"onichankimochi.com/zen_cat_backend/src/server/bll/controller"
	"onichankimochi.com/zen_cat_backend/src/server/schemas"
)

type Api struct {
	Logger        logging.Logger
	BllController *controller.ControllerCollection
	EnvSettings   *schemas.EnvSettings
	Echo          *echo.Echo
}

// TODO: Stablish a connection to the ZenCat database
/*
Creates a new api server with
- Logger provided by input
- BllController as new bll controller collection
- EnvSettings as new env settings provided by .env file
*/
func NewApi(
	logger logging.Logger,
	envSettings *schemas.EnvSettings,
) *Api {
	// bllController, zenCatPsqlDB := controller.NewControllerCollection(logger, envSettings)

	return &Api{
		Logger: logger,
		// BllController: bllController,
		EnvSettings: envSettings,
		Echo:        echo.New(),
	}
}

// @title ZenCat API
// @version 1.0
// @description ZenCat API sample for clients
// @BasePath /
func RunService(envSettings *schemas.EnvSettings, logger logging.Logger) {
	api := NewApi(logger, envSettings)
	api.RunApi(envSettings)
}
