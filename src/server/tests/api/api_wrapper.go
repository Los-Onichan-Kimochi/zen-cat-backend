package api_test

import (
	"testing"

	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/api"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	testSetup "onichankimochi.com/astro_cat_backend/src/server/tests"
)

type ApiWrapper struct {
	server      *api.Api
	astroCatDB  *gorm.DB
	envSettings *schemas.EnvSettings
}

// Create new api wrapper
func NewApiTestWrapper() *ApiWrapper {
	testLogger := logging.NewLoggerMock()
	envSettings := schemas.NewEnvSettings(testLogger)
	envSettings.DisableAuthForTests = true
	envSettings.EnableSqlLogs = false // Disable SQL logs for testing
	server, astroCatDB := api.NewApi(testLogger, envSettings)

	// Register routes but don't start the HTTP server
	server.RegisterRoutes(envSettings)

	return &ApiWrapper{
		server:      server,
		astroCatDB:  astroCatDB,
		envSettings: envSettings,
	}
}

// Restart astro cat db based on env settings and testing parameter
func (a *ApiWrapper) restartDB(t *testing.T) {
	testSetup.ClearPostgresqlDatabaseTesting(
		a.server.Logger,
		a.astroCatDB,
		a.envSettings,
		t,
	)
}

/*
--------------------------------
	Api test wrappers
--------------------------------
*/

// Create api server wrapper - each test gets its own isolated instance
func NewApiServerTestWrapper(t *testing.T) (*api.Api, *gorm.DB) {
	// Create a fresh wrapper for each test to ensure isolation
	apiTestWrapper := NewApiTestWrapper()
	apiTestWrapper.restartDB(t)
	return apiTestWrapper.server, apiTestWrapper.astroCatDB
}
