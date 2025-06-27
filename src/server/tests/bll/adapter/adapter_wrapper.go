package adapter_test

import (
	"testing"

	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/bll/adapter"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	testSetup "onichankimochi.com/astro_cat_backend/src/server/tests"
)

type AdapterTestWrapper struct {
	logger         logging.Logger
	testAdapter    *adapter.AdapterCollection
	astroCatPsqlDB *gorm.DB
	envSettings    *schemas.EnvSettings
}

func newAdapterTestWrapper() *AdapterTestWrapper {
	testLogger := logging.NewLoggerMock()
	envSettings := schemas.NewEnvSettings(testLogger)
	envSettings.EnableSqlLogs = false // Disable SQL logs for testing
	testAdapter, astroCatPsqlDB := adapter.NewAdapterCollection(
		testLogger,
		envSettings,
	)

	return &AdapterTestWrapper{
		logger:         testLogger,
		testAdapter:    testAdapter,
		astroCatPsqlDB: astroCatPsqlDB,
		envSettings:    envSettings,
	}
}

// Restart astro cat database based on env settings and testing parameter
func (atw *AdapterTestWrapper) restartDB(t *testing.T) {
	testSetup.ClearPostgresqlDatabaseTesting(
		atw.logger,
		atw.astroCatPsqlDB,
		atw.envSettings,
		t,
	)
}

var adapterTestWrapper = newAdapterTestWrapper()

/*
--------------------------------
	Adapter test wrappers
--------------------------------
*/

// Create new user adapter wrapper
func NewUserAdapterTestWrapper(
	t *testing.T,
) (*adapter.User, *logging.LoggerMock, *gorm.DB) {
	adapterTestWrapper.restartDB(t)
	loggerMock := adapterTestWrapper.logger.(*logging.LoggerMock)
	return adapterTestWrapper.testAdapter.User, loggerMock, adapterTestWrapper.astroCatPsqlDB
}

// Create new community adapter wrapper
func NewCommunityAdapterTestWrapper(
	t *testing.T,
) (*adapter.Community, *logging.LoggerMock, *gorm.DB) {
	adapterTestWrapper.restartDB(t)
	loggerMock := adapterTestWrapper.logger.(*logging.LoggerMock)
	return adapterTestWrapper.testAdapter.Community, loggerMock, adapterTestWrapper.astroCatPsqlDB
}

// Create new plan adapter wrapper
func NewPlanAdapterTestWrapper(
	t *testing.T,
) (*adapter.Plan, *logging.LoggerMock, *gorm.DB) {
	adapterTestWrapper.restartDB(t)
	loggerMock := adapterTestWrapper.logger.(*logging.LoggerMock)
	return adapterTestWrapper.testAdapter.Plan, loggerMock, adapterTestWrapper.astroCatPsqlDB
}

// Create new reservation adapter wrapper
func NewReservationAdapterTestWrapper(
	t *testing.T,
) (*adapter.Reservation, *logging.LoggerMock, *gorm.DB) {
	adapterTestWrapper.restartDB(t)
	loggerMock := adapterTestWrapper.logger.(*logging.LoggerMock)
	return adapterTestWrapper.testAdapter.Reservation, loggerMock, adapterTestWrapper.astroCatPsqlDB
}

// Create new service adapter wrapper
func NewServiceAdapterTestWrapper(
	t *testing.T,
) (*adapter.Service, *logging.LoggerMock, *gorm.DB) {
	adapterTestWrapper.restartDB(t)
	loggerMock := adapterTestWrapper.logger.(*logging.LoggerMock)
	return adapterTestWrapper.testAdapter.Service, loggerMock, adapterTestWrapper.astroCatPsqlDB
}

// Create new session adapter wrapper
func NewSessionAdapterTestWrapper(
	t *testing.T,
) (*adapter.Session, *logging.LoggerMock, *gorm.DB) {
	adapterTestWrapper.restartDB(t)
	loggerMock := adapterTestWrapper.logger.(*logging.LoggerMock)
	return adapterTestWrapper.testAdapter.Session, loggerMock, adapterTestWrapper.astroCatPsqlDB
}

// Create new professional adapter wrapper
func NewProfessionalAdapterTestWrapper(
	t *testing.T,
) (*adapter.Professional, *logging.LoggerMock, *gorm.DB) {
	adapterTestWrapper.restartDB(t)
	loggerMock := adapterTestWrapper.logger.(*logging.LoggerMock)
	return adapterTestWrapper.testAdapter.Professional, loggerMock, adapterTestWrapper.astroCatPsqlDB
}

// Create new local adapter wrapper
func NewLocalAdapterTestWrapper(
	t *testing.T,
) (*adapter.Local, *logging.LoggerMock, *gorm.DB) {
	adapterTestWrapper.restartDB(t)
	loggerMock := adapterTestWrapper.logger.(*logging.LoggerMock)
	return adapterTestWrapper.testAdapter.Local, loggerMock, adapterTestWrapper.astroCatPsqlDB
}

// Create new membership adapter wrapper
func NewMembershipAdapterTestWrapper(
	t *testing.T,
) (*adapter.Membership, *logging.LoggerMock, *gorm.DB) {
	adapterTestWrapper.restartDB(t)
	loggerMock := adapterTestWrapper.logger.(*logging.LoggerMock)
	return adapterTestWrapper.testAdapter.Membership, loggerMock, adapterTestWrapper.astroCatPsqlDB
}

// Create new onboarding adapter wrapper
func NewOnboardingAdapterTestWrapper(
	t *testing.T,
) (*adapter.Onboarding, *logging.LoggerMock, *gorm.DB) {
	adapterTestWrapper.restartDB(t)
	loggerMock := adapterTestWrapper.logger.(*logging.LoggerMock)
	return adapterTestWrapper.testAdapter.Onboarding, loggerMock, adapterTestWrapper.astroCatPsqlDB
}

// Create new community plan adapter wrapper
func NewCommunityPlanAdapterTestWrapper(
	t *testing.T,
) (*adapter.CommunityPlan, *logging.LoggerMock, *gorm.DB) {
	adapterTestWrapper.restartDB(t)
	loggerMock := adapterTestWrapper.logger.(*logging.LoggerMock)
	return adapterTestWrapper.testAdapter.CommunityPlan, loggerMock, adapterTestWrapper.astroCatPsqlDB
}

// Create new community service adapter wrapper
func NewCommunityServiceAdapterTestWrapper(
	t *testing.T,
) (*adapter.CommunityService, *logging.LoggerMock, *gorm.DB) {
	adapterTestWrapper.restartDB(t)
	loggerMock := adapterTestWrapper.logger.(*logging.LoggerMock)
	return adapterTestWrapper.testAdapter.CommunityService, loggerMock, adapterTestWrapper.astroCatPsqlDB
}

// Create new service local adapter wrapper
func NewServiceLocalAdapterTestWrapper(
	t *testing.T,
) (*adapter.ServiceLocal, *logging.LoggerMock, *gorm.DB) {
	adapterTestWrapper.restartDB(t)
	loggerMock := adapterTestWrapper.logger.(*logging.LoggerMock)
	return adapterTestWrapper.testAdapter.ServiceLocal, loggerMock, adapterTestWrapper.astroCatPsqlDB
}

// Create new service professional adapter wrapper
func NewServiceProfessionalAdapterTestWrapper(
	t *testing.T,
) (*adapter.ServiceProfessional, *logging.LoggerMock, *gorm.DB) {
	adapterTestWrapper.restartDB(t)
	loggerMock := adapterTestWrapper.logger.(*logging.LoggerMock)
	return adapterTestWrapper.testAdapter.ServiceProfessional, loggerMock, adapterTestWrapper.astroCatPsqlDB
}

// Create new audit log adapter wrapper
func NewAuditLogAdapterTestWrapper(
	t *testing.T,
) (*adapter.AuditLog, *logging.LoggerMock, *gorm.DB) {
	adapterTestWrapper.restartDB(t)
	loggerMock := adapterTestWrapper.logger.(*logging.LoggerMock)
	return adapterTestWrapper.testAdapter.AuditLog, loggerMock, adapterTestWrapper.astroCatPsqlDB
}
