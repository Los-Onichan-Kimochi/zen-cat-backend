package controller_test

import (
	"testing"

	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/bll/controller"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	testSetup "onichankimochi.com/astro_cat_backend/src/server/tests"
)

type ControllerTestWrapper struct {
	logger         logging.Logger
	testController *controller.ControllerCollection
	astroCatPsqlDB *gorm.DB
	envSettings    *schemas.EnvSettings
}

func newControllerTestWrapper() *ControllerTestWrapper {
	testLogger := logging.NewLoggerMock()
	envSettings := schemas.NewEnvSettings(testLogger)
	envSettings.EnableSqlLogs = false // Disable SQL logs for testing
	testController, astroCatPsqlDB := controller.NewControllerCollection(
		testLogger,
		envSettings,
	)

	return &ControllerTestWrapper{
		logger:         testLogger,
		testController: testController,
		astroCatPsqlDB: astroCatPsqlDB,
		envSettings:    envSettings,
	}
}

// Restart astro cat database based on env settings and testing parameter
func (ctw *ControllerTestWrapper) restartDB(t *testing.T) {
	testSetup.ClearPostgresqlDatabase(
		ctw.logger,
		ctw.astroCatPsqlDB,
		ctw.envSettings,
		t,
	)
}

var controllerTestWrapper = newControllerTestWrapper()

/*
--------------------------------
	Controller test wrappers
--------------------------------
*/

// Create new user controller wrapper
func NewUserControllerTestWrapper(
	t *testing.T,
) (*controller.User, *logging.LoggerMock, *gorm.DB) {
	controllerTestWrapper.restartDB(t)
	loggerMock := controllerTestWrapper.logger.(*logging.LoggerMock)
	return controllerTestWrapper.testController.User, loggerMock,
		controllerTestWrapper.astroCatPsqlDB
}

// Create new community controller wrapper
func NewCommunityControllerTestWrapper(
	t *testing.T,
) (*controller.Community, *logging.LoggerMock, *gorm.DB) {
	controllerTestWrapper.restartDB(t)
	loggerMock := controllerTestWrapper.logger.(*logging.LoggerMock)
	return controllerTestWrapper.testController.Community, loggerMock,
		controllerTestWrapper.astroCatPsqlDB
}

// Create new plan controller wrapper
func NewPlanControllerTestWrapper(
	t *testing.T,
) (*controller.Plan, *logging.LoggerMock, *gorm.DB) {
	controllerTestWrapper.restartDB(t)
	loggerMock := controllerTestWrapper.logger.(*logging.LoggerMock)
	return controllerTestWrapper.testController.Plan, loggerMock,
		controllerTestWrapper.astroCatPsqlDB
}

// Create new reservation controller wrapper
func NewReservationControllerTestWrapper(
	t *testing.T,
) (*controller.Reservation, *logging.LoggerMock, *gorm.DB) {
	controllerTestWrapper.restartDB(t)
	loggerMock := controllerTestWrapper.logger.(*logging.LoggerMock)
	return controllerTestWrapper.testController.Reservation, loggerMock,
		controllerTestWrapper.astroCatPsqlDB
}

// Create new service controller wrapper
func NewServiceControllerTestWrapper(
	t *testing.T,
) (*controller.Service, *logging.LoggerMock, *gorm.DB) {
	controllerTestWrapper.restartDB(t)
	loggerMock := controllerTestWrapper.logger.(*logging.LoggerMock)
	return controllerTestWrapper.testController.Service, loggerMock,
		controllerTestWrapper.astroCatPsqlDB
}

// Create new session controller wrapper
func NewSessionControllerTestWrapper(
	t *testing.T,
) (*controller.Session, *logging.LoggerMock, *gorm.DB) {
	controllerTestWrapper.restartDB(t)
	loggerMock := controllerTestWrapper.logger.(*logging.LoggerMock)
	return controllerTestWrapper.testController.Session, loggerMock,
		controllerTestWrapper.astroCatPsqlDB
}

// Create new auth controller wrapper
func NewAuthControllerTestWrapper(
	t *testing.T,
) (*controller.Auth, *logging.LoggerMock, *gorm.DB) {
	controllerTestWrapper.restartDB(t)
	loggerMock := controllerTestWrapper.logger.(*logging.LoggerMock)
	return controllerTestWrapper.testController.Auth, loggerMock,
		controllerTestWrapper.astroCatPsqlDB
}

// Create new professional controller wrapper
func NewProfessionalControllerTestWrapper(
	t *testing.T,
) (*controller.Professional, *logging.LoggerMock, *gorm.DB) {
	controllerTestWrapper.restartDB(t)
	loggerMock := controllerTestWrapper.logger.(*logging.LoggerMock)
	return controllerTestWrapper.testController.Professional, loggerMock,
		controllerTestWrapper.astroCatPsqlDB
}

// Create new local controller wrapper
func NewLocalControllerTestWrapper(
	t *testing.T,
) (*controller.Local, *logging.LoggerMock, *gorm.DB) {
	controllerTestWrapper.restartDB(t)
	loggerMock := controllerTestWrapper.logger.(*logging.LoggerMock)
	return controllerTestWrapper.testController.Local, loggerMock,
		controllerTestWrapper.astroCatPsqlDB
}

// Create new membership controller wrapper
func NewMembershipControllerTestWrapper(
	t *testing.T,
) (*controller.Membership, *logging.LoggerMock, *gorm.DB) {
	controllerTestWrapper.restartDB(t)
	loggerMock := controllerTestWrapper.logger.(*logging.LoggerMock)
	return controllerTestWrapper.testController.Membership, loggerMock,
		controllerTestWrapper.astroCatPsqlDB
}

// Create new onboarding controller wrapper
func NewOnboardingControllerTestWrapper(
	t *testing.T,
) (*controller.Onboarding, *logging.LoggerMock, *gorm.DB) {
	controllerTestWrapper.restartDB(t)
	loggerMock := controllerTestWrapper.logger.(*logging.LoggerMock)
	return controllerTestWrapper.testController.Onboarding, loggerMock,
		controllerTestWrapper.astroCatPsqlDB
}

// Create new community plan controller wrapper
func NewCommunityPlanControllerTestWrapper(
	t *testing.T,
) (*controller.CommunityPlan, *logging.LoggerMock, *gorm.DB) {
	controllerTestWrapper.restartDB(t)
	loggerMock := controllerTestWrapper.logger.(*logging.LoggerMock)
	return controllerTestWrapper.testController.CommunityPlan, loggerMock,
		controllerTestWrapper.astroCatPsqlDB
}

// Create new community service controller wrapper
func NewCommunityServiceControllerTestWrapper(
	t *testing.T,
) (*controller.CommunityService, *logging.LoggerMock, *gorm.DB) {
	controllerTestWrapper.restartDB(t)
	loggerMock := controllerTestWrapper.logger.(*logging.LoggerMock)
	return controllerTestWrapper.testController.CommunityService, loggerMock,
		controllerTestWrapper.astroCatPsqlDB
}

// Create new service local controller wrapper
func NewServiceLocalControllerTestWrapper(
	t *testing.T,
) (*controller.ServiceLocal, *logging.LoggerMock, *gorm.DB) {
	controllerTestWrapper.restartDB(t)
	loggerMock := controllerTestWrapper.logger.(*logging.LoggerMock)
	return controllerTestWrapper.testController.ServiceLocal, loggerMock,
		controllerTestWrapper.astroCatPsqlDB
}

// Create new service professional controller wrapper
func NewServiceProfessionalControllerTestWrapper(
	t *testing.T,
) (*controller.ServiceProfessional, *logging.LoggerMock, *gorm.DB) {
	controllerTestWrapper.restartDB(t)
	loggerMock := controllerTestWrapper.logger.(*logging.LoggerMock)
	return controllerTestWrapper.testController.ServiceProfessional, loggerMock,
		controllerTestWrapper.astroCatPsqlDB
}

// Create new audit log controller wrapper
func NewAuditLogControllerTestWrapper(
	t *testing.T,
) (*controller.AuditLog, *logging.LoggerMock, *gorm.DB) {
	controllerTestWrapper.restartDB(t)
	loggerMock := controllerTestWrapper.logger.(*logging.LoggerMock)
	return controllerTestWrapper.testController.AuditLog, loggerMock,
		controllerTestWrapper.astroCatPsqlDB
}

// Create new forgot password controller wrapper
func NewForgotPasswordControllerTestWrapper(
	t *testing.T,
) (*controller.ForgotPassword, *logging.LoggerMock, *gorm.DB) {
	controllerTestWrapper.restartDB(t)
	loggerMock := controllerTestWrapper.logger.(*logging.LoggerMock)
	return controllerTestWrapper.testController.ForgotPassword, loggerMock,
		controllerTestWrapper.astroCatPsqlDB
}

// Create new login controller wrapper
func NewLoginControllerTestWrapper(
	t *testing.T,
) (*controller.Login, *logging.LoggerMock, *gorm.DB) {
	controllerTestWrapper.restartDB(t)
	loggerMock := controllerTestWrapper.logger.(*logging.LoggerMock)
	return controllerTestWrapper.testController.Login, loggerMock,
		controllerTestWrapper.astroCatPsqlDB
}
