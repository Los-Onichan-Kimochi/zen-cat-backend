package adapter

import (
	"github.com/google/uuid"
	daoPostgresql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"

	"onichankimochi.com/astro_cat_backend/src/logging"
)

type Service struct {
	logger        logging.Logger
	DaoPostgresql *daoPostgresql.AstroCatPsqlCollection
}

// Creates Service adapter
func NewServiceAdapter(
	logger logging.Logger,
	daoPostgresql *daoPostgresql.AstroCatPsqlCollection,
) *Service {
	return &Service{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

// Gets a service from postgresql DB and adapts it to a Service schema.
func (c *Service) GetPostgresqlService(
	serviceId uuid.UUID,
) (*schemas.Service, *errors.Error) {
	serviceModel, err := c.DaoPostgresql.Service.GetService(serviceId)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.ServiceNotFound
	}

	return &schemas.Service{
		Id:                  serviceModel.Id,
		Name:                serviceModel.Name,
		Description: 	   	 serviceModel.Description,
		ImageUrl:            serviceModel.ImageUrl,
		IsVirtual:           serviceModel.IsVirtual,
	}, nil
}

// Fetch services from postgresql DB and adapts them to a Service schema.
func (c *Service) FetchPostgresqlServices() ([]*schemas.Service, *errors.Error) {
	servicesModel, err := c.DaoPostgresql.Service.FetchServices()
	if err != nil {
		return nil, &errors.ObjectNotFoundError.ServiceNotFound
	}

	services := make([]*schemas.Service, len(servicesModel))
	for i, serviceModel := range servicesModel {
		services[i] = &schemas.Service{
			Id:                  serviceModel.Id,
			Name:                serviceModel.Name,
			Description: 	   	 serviceModel.Description,
			ImageUrl:            serviceModel.ImageUrl,
			IsVirtual:           serviceModel.IsVirtual,
		}
	}

	return services, nil
}

// Creates a service into postgresql DB and returns it.
func (c *Service) CreatePostgresqlService(
	name string,
	description string,
	imageUrl string,
	isVirtual bool,
	updatedBy string,
) (*schemas.Service, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	serviceModel := &model.Service{
		Id:                  uuid.New(),
		Name:                name,
		Description:         description,
		ImageUrl:            imageUrl,
		IsVirtual:           isVirtual,
		AuditFields: model.AuditFields{
			UpdatedBy: updatedBy,
		},
	}

	if err := c.DaoPostgresql.Service.CreateService(serviceModel); err != nil {
		return nil, &errors.BadRequestError.ServiceNotCreated
	}

	return &schemas.Service{
		Id:                  serviceModel.Id,
		Name:                serviceModel.Name,
		Description:         serviceModel.Description,
		ImageUrl:            serviceModel.ImageUrl,
		IsVirtual: 			 serviceModel.IsVirtual,
	}, nil
}

// Updates a service given fields in postgresql DB and returns it.
func (c *Service) UpdatePostgresqlService(
	id uuid.UUID,
	name *string,
	description *string,
	imageUrl *string,
	isVirtual *bool,
	updatedBy string,
) (*schemas.Service, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	serviceModel, err := c.DaoPostgresql.Service.UpdateService(
		id,
		name,
		description,
		imageUrl,
		isVirtual,
		updatedBy,
	)
	if err != nil {
		return nil, &errors.BadRequestError.ServiceNotUpdated
	}

	return &schemas.Service{
		Id:                 serviceModel.Id,
		Name:               serviceModel.Name,
		Description:        serviceModel.Description,
		ImageUrl:           serviceModel.ImageUrl,
		IsVirtual: 			serviceModel.IsVirtual,
	}, nil
}
