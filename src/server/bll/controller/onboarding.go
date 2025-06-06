package controller

import (
	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
	bllAdapter "onichankimochi.com/astro_cat_backend/src/server/bll/adapter"
	errors "onichankimochi.com/astro_cat_backend/src/server/errors"
	schemas "onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type Onboarding struct {
	logger      logging.Logger
	Adapter     *bllAdapter.AdapterCollection
	EnvSettings *schemas.EnvSettings
}

func NewOnboardingController(
	logger logging.Logger,
	adapter *bllAdapter.AdapterCollection,
	envSettings *schemas.EnvSettings,
) *Onboarding {
	return &Onboarding{
		logger:      logger,
		Adapter:     adapter,
		EnvSettings: envSettings,
	}
}

func (o *Onboarding) GetOnboarding(onboardingId uuid.UUID) (*schemas.Onboarding, *errors.Error) {
	return o.Adapter.Onboarding.GetPostgresqlOnboarding(onboardingId)
}

func (o *Onboarding) GetOnboardingByUserId(userId uuid.UUID) (*schemas.Onboarding, *errors.Error) {
	return o.Adapter.Onboarding.GetPostgresqlOnboardingByUserId(userId)
}

func (o *Onboarding) FetchOnboardings() (*schemas.Onboardings, *errors.Error) {
	onboardings, err := o.Adapter.Onboarding.FetchPostgresqlOnboardings()
	if err != nil {
		return nil, err
	}

	return &schemas.Onboardings{Onboardings: onboardings}, nil
}

// CreateOnboardingForUser crea un onboarding ligado a un usuario específico
func (o *Onboarding) CreateOnboardingForUser(
	userId uuid.UUID,
	createOnboardingRequest schemas.CreateOnboardingRequest,
	updatedBy string,
) (*schemas.Onboarding, *errors.Error) {
	// Validar que el usuario existe antes de crear el onboarding
	_, userErr := o.Adapter.User.GetPostgresqlUser(userId)
	if userErr != nil {
		return nil, userErr
	}

	// Convertir la fecha de nacimiento a string si existe
	var birthDateStr *string
	if createOnboardingRequest.BirthDate != nil {
		dateStr := createOnboardingRequest.BirthDate.Format("2006-01-02")
		birthDateStr = &dateStr
	}

	return o.Adapter.Onboarding.CreatePostgresqlOnboarding(
		userId, // El onboarding siempre va ligado al usuario
		createOnboardingRequest.DocumentType,
		createOnboardingRequest.DocumentNumber,
		createOnboardingRequest.PhoneNumber,
		birthDateStr,
		createOnboardingRequest.Gender,
		createOnboardingRequest.City,
		createOnboardingRequest.PostalCode,
		createOnboardingRequest.District,
		createOnboardingRequest.Address,
		updatedBy,
	)
}

func (o *Onboarding) UpdateOnboarding(
	onboardingId uuid.UUID,
	updateOnboardingRequest schemas.UpdateOnboardingRequest,
	updatedBy string,
) (*schemas.Onboarding, *errors.Error) {
	// Convertir la fecha de nacimiento a string si existe
	var birthDateStr *string
	if updateOnboardingRequest.BirthDate != nil {
		dateStr := updateOnboardingRequest.BirthDate.Format("2006-01-02")
		birthDateStr = &dateStr
	}

	return o.Adapter.Onboarding.UpdatePostgresqlOnboarding(
		onboardingId,
		updateOnboardingRequest.DocumentType,
		updateOnboardingRequest.DocumentNumber,
		updateOnboardingRequest.PhoneNumber,
		birthDateStr,
		updateOnboardingRequest.Gender,
		updateOnboardingRequest.City,
		updateOnboardingRequest.PostalCode,
		updateOnboardingRequest.District,
		updateOnboardingRequest.Address,
		updatedBy,
	)
}

// UpdateOnboardingByUserId actualiza el onboarding usando el userId (más conveniente para el frontend)
func (o *Onboarding) UpdateOnboardingByUserId(
	userId uuid.UUID,
	updateOnboardingRequest schemas.UpdateOnboardingRequest,
	updatedBy string,
) (*schemas.Onboarding, *errors.Error) {
	// Primero obtener el onboarding por userId
	existingOnboarding, err := o.Adapter.Onboarding.GetPostgresqlOnboardingByUserId(userId)
	if err != nil {
		return nil, err
	}

	// Ahora actualizar usando el ID del onboarding
	return o.UpdateOnboarding(existingOnboarding.Id, updateOnboardingRequest, updatedBy)
}

func (o *Onboarding) DeleteOnboarding(onboardingId uuid.UUID) *errors.Error {
	return o.Adapter.Onboarding.DeletePostgresqlOnboarding(onboardingId)
}

func (o *Onboarding) DeleteOnboardingByUserId(userId uuid.UUID) *errors.Error {
	return o.Adapter.Onboarding.DeletePostgresqlOnboardingByUserId(userId)
}
