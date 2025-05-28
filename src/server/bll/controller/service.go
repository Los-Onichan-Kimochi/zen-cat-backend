package controller

import (
	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
	bllAdapter "onichankimochi.com/astro_cat_backend/src/server/bll/adapter"
	errors "onichankimochi.com/astro_cat_backend/src/server/errors"
	schemas "onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type Service struct {
	logger      logging.Logger
	Adapter     *bllAdapter.AdapterCollection
	EnvSettings *schemas.EnvSettings
}

// Create Service controller
func NewServiceController(
	logger logging.Logger,
	adapter *bllAdapter.AdapterCollection,
	envSettings *schemas.EnvSettings,
) *Service {
	return &Service{
		logger:      logger,
		Adapter:     adapter,
		EnvSettings: envSettings,
	}
}

// Gets a service.
func (c *Service) GetService(serviceId uuid.UUID) (*schemas.Service, *errors.Error) {
	return c.Adapter.Service.GetPostgresqlService(serviceId)
}

// Fetch all services.
func (c *Service) FetchServices() (*schemas.Services, *errors.Error) {
	services, err := c.Adapter.Service.FetchPostgresqlServices()
	if err != nil {
		return nil, err
	}

	return &schemas.Services{Services: services}, nil
}

// Creates a service.
func (c *Service) CreateService(
	createServiceData schemas.CreateServiceRequest,
	updatedBy string,
) (*schemas.Service, *errors.Error) {
	return c.Adapter.Service.CreatePostgresqlService(
		createServiceData.Name,
		createServiceData.Description,
		createServiceData.ImageUrl,
		createServiceData.IsVirtual,
		updatedBy,
	)
}

// Updates a service.
func (c *Service) UpdateService(
	serviceId uuid.UUID,
	updateServiceData schemas.UpdateServiceRequest,
	updatedBy string,
) (*schemas.Service, *errors.Error) {
	return c.Adapter.Service.UpdatePostgresqlService(
		serviceId,
		updateServiceData.Name,
		updateServiceData.Description,
		updateServiceData.ImageUrl,
		updateServiceData.IsVirtual,
		updatedBy,
	)
}

// Deletes a service.
func (l *Service) DeleteService(serviceId uuid.UUID) *errors.Error {
	return l.Adapter.Service.DeletePostgresqlService(serviceId)
}

// TODO: Add BulkCreateCommunities (Batch)
