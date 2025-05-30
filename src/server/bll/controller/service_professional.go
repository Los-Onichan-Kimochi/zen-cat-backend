package controller

import (

	// "gorm.io/gorm" // No longer directly needed here for these checks
	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
	bllAdapter "onichankimochi.com/astro_cat_backend/src/server/bll/adapter"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type ServiceProfessional struct {
	logger      logging.Logger
	Adapter     *bllAdapter.AdapterCollection
	EnvSettings *schemas.EnvSettings
}

// Create ServiceProfessional controller
func NewServiceProfessionalController(
	logger logging.Logger,
	adapter *bllAdapter.AdapterCollection,
	envSettings *schemas.EnvSettings,
) *ServiceProfessional {
	return &ServiceProfessional{
		logger:      logger,
		Adapter:     adapter,
		EnvSettings: envSettings,
	}
}

// Creates a service-professional association.
func (cp *ServiceProfessional) CreateServiceProfessional(
	req schemas.CreateServiceProfessionalRequest,
	updatedBy string,
) (*schemas.ServiceProfessional, *errors.Error) {
	serviceId := req.ServiceId
	professionalId := req.ProfessionalId

	_, err := cp.Adapter.Service.GetPostgresqlService(serviceId)
	if err != nil {
		return nil, err
	}

	_, err = cp.Adapter.Professional.GetPostgresqlProfessional(professionalId)
	if err != nil {
		return nil, err
	}

	_, err = cp.Adapter.ServiceProfessional.GetPostgresqlServiceProfessional(serviceId, professionalId)
	if err == nil {
		return nil, &errors.BadRequestError.ServiceProfessionalAlreadyExists
	} else if err.Code != errors.ObjectNotFoundError.ServiceProfessionalNotFound.Code {
		return nil, &errors.InternalServerError.Default
	}

	return cp.Adapter.ServiceProfessional.CreatePostgresqlServiceProfessional(serviceId, professionalId, updatedBy)
}

// Gets a specific service-professional association.
func (cp *ServiceProfessional) GetServiceProfessional(
	serviceIdString string,
	professionalIdString string,
) (*schemas.ServiceProfessional, *errors.Error) {
	serviceId, err := uuid.Parse(serviceIdString)
	if err != nil {
		return nil, &errors.UnprocessableEntityError.InvalidServiceId
	}

	professionalId, err := uuid.Parse(professionalIdString)
	if err != nil {
		return nil, &errors.UnprocessableEntityError.InvalidProfessionalId
	}

	return cp.Adapter.ServiceProfessional.GetPostgresqlServiceProfessional(serviceId, professionalId)
}

// Deletes a specific service-professional association.
func (cp *ServiceProfessional) DeleteServiceProfessional(
	serviceIdString string,
	professionalIdString string,
) *errors.Error {
	serviceId, parseErr := uuid.Parse(serviceIdString)
	if parseErr != nil {
		return &errors.UnprocessableEntityError.InvalidServiceId
	}

	professionalId, parseErr := uuid.Parse(professionalIdString)
	if parseErr != nil {
		return &errors.UnprocessableEntityError.InvalidProfessionalId
	}

	_, err := cp.Adapter.ServiceProfessional.GetPostgresqlServiceProfessional(serviceId, professionalId)
	if err != nil {
		return err
	}

	return cp.Adapter.ServiceProfessional.DeletePostgresqlServiceProfessional(serviceId, professionalId)
}
