package controller

import (

	// "gorm.io/gorm" // No longer directly needed here for these checks
	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
	bllAdapter "onichankimochi.com/astro_cat_backend/src/server/bll/adapter"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type ServiceLocal struct {
	logger      logging.Logger
	Adapter     *bllAdapter.AdapterCollection
	EnvSettings *schemas.EnvSettings
}

// Create ServiceLocal controller
func NewServiceLocalController(
	logger logging.Logger,
	adapter *bllAdapter.AdapterCollection,
	envSettings *schemas.EnvSettings,
) *ServiceLocal {
	return &ServiceLocal{
		logger:      logger,
		Adapter:     adapter,
		EnvSettings: envSettings,
	}
}

// Creates a service-local association.
func (cp *ServiceLocal) CreateServiceLocal(
	req schemas.CreateServiceLocalRequest,
	updatedBy string,
) (*schemas.ServiceLocal, *errors.Error) {
	serviceId := req.ServiceId
	localId := req.LocalId

	_, err := cp.Adapter.Service.GetPostgresqlService(serviceId)
	if err != nil {
		return nil, err
	}

	_, err = cp.Adapter.Local.GetPostgresqlLocal(localId)
	if err != nil {
		return nil, err
	}

	_, err = cp.Adapter.ServiceLocal.GetPostgresqlServiceLocal(serviceId, localId)
	if err == nil {
		return nil, &errors.BadRequestError.ServiceLocalAlreadyExists
	} else if err.Code != errors.ObjectNotFoundError.ServiceLocalNotFound.Code {
		return nil, &errors.InternalServerError.Default
	}

	return cp.Adapter.ServiceLocal.CreatePostgresqlServiceLocal(serviceId, localId, updatedBy)
}

// Gets a specific service-local association.
func (cp *ServiceLocal) GetServiceLocal(
	serviceIdString string,
	localIdString string,
) (*schemas.ServiceLocal, *errors.Error) {
	serviceId, err := uuid.Parse(serviceIdString)
	if err != nil {
		return nil, &errors.UnprocessableEntityError.InvalidServiceId
	}

	localId, err := uuid.Parse(localIdString)
	if err != nil {
		return nil, &errors.UnprocessableEntityError.InvalidLocalId
	}

	return cp.Adapter.ServiceLocal.GetPostgresqlServiceLocal(serviceId, localId)
}

// Deletes a specific service-local association.
func (cp *ServiceLocal) DeleteServiceLocal(
	serviceIdString string,
	localIdString string,
) *errors.Error {
	serviceId, parseErr := uuid.Parse(serviceIdString)
	if parseErr != nil {
		return &errors.UnprocessableEntityError.InvalidServiceId
	}

	localId, parseErr := uuid.Parse(localIdString)
	if parseErr != nil {
		return &errors.UnprocessableEntityError.InvalidLocalId
	}

	_, err := cp.Adapter.ServiceLocal.GetPostgresqlServiceLocal(serviceId, localId)
	if err != nil {
		return err
	}

	return cp.Adapter.ServiceLocal.DeletePostgresqlServiceLocal(serviceId, localId)
}
