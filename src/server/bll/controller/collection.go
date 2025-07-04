package controller

import (
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/bll/adapter"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type ControllerCollection struct {
	Logger              logging.Logger
	EnvSettings         *schemas.EnvSettings
	Auth                *Auth
	Login               *Login
	Community           *Community
	Professional        *Professional
	Local               *Local
	User                *User
	Onboarding          *Onboarding
	Membership          *Membership
	Service             *Service
	Plan                *Plan
	CommunityPlan       *CommunityPlan
	CommunityService    *CommunityService
	ServiceLocal        *ServiceLocal
	ServiceProfessional *ServiceProfessional
	Session             *Session
	Reservation         *Reservation
	ForgotPassword      *ForgotPassword
	Contact             *Contact
	AuditLog            *AuditLog
}

// Create bll controller collection
func NewControllerCollection(
	logger logging.Logger,
	envSettings *schemas.EnvSettings,
) (*ControllerCollection, *gorm.DB) {
	bllAdapter, astroCatPsqlDB := adapter.NewAdapterCollection(
		logger,
		envSettings,
	)
	auth := NewAuthController(logger, bllAdapter, envSettings)
	login := NewLoginController(logger, bllAdapter, envSettings, auth)
	community := NewCommunityController(logger, bllAdapter, envSettings)
	professional := NewProfessionalController(logger, bllAdapter, envSettings)
	local := NewLocalController(logger, bllAdapter, envSettings)
	user := NewUserController(logger, bllAdapter, envSettings)
	onboarding := NewOnboardingController(logger, bllAdapter, envSettings)
	membership := NewMembershipController(logger, bllAdapter, envSettings)
	service := NewServiceController(logger, bllAdapter, envSettings)
	plan := NewPlanController(logger, bllAdapter, envSettings)
	communityPlan := NewCommunityPlanController(logger, bllAdapter, envSettings)
	communityService := NewCommunityServiceController(logger, bllAdapter, envSettings)
	serviceLocal := NewServiceLocalController(logger, bllAdapter, envSettings)
	serviceProfessional := NewServiceProfessionalController(logger, bllAdapter, envSettings)
	session := NewSessionController(logger, bllAdapter, envSettings)
	reservation := NewReservationController(logger, bllAdapter, envSettings)
	forgotPassword := NewForgotPasswordController(logger, bllAdapter, envSettings)
	contact := NewContactController(logger, bllAdapter, envSettings)
	auditLog := NewAuditLogController(logger, bllAdapter, envSettings)

	return &ControllerCollection{
		Logger:              logger,
		EnvSettings:         envSettings,
		Auth:                auth,
		Login:               login,
		Community:           community,
		Professional:        professional,
		Local:               local,
		User:                user,
		Onboarding:          onboarding,
		Membership:          membership,
		Service:             service,
		Plan:                plan,
		CommunityPlan:       communityPlan,
		CommunityService:    communityService,
		ServiceLocal:        serviceLocal,
		ServiceProfessional: serviceProfessional,
		Session:             session,
		Reservation:         reservation,
		ForgotPassword:      forgotPassword,
		Contact:             contact,
		AuditLog:            auditLog,
	}, astroCatPsqlDB
}
