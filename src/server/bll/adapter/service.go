package adapter

import (
	"github.com/google/uuid"
	daoPostgresql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"

	"gorm.io/gorm"
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

// Gets a service from a Postgresql DB given its ID and adapts it to a service schema.
func (s *Service) GetPostgresqlService(id uuid.UUID) (*schemas.Service, *errors.Error) {
	serviceModel, err := s.DaoPostgresql.Service.GetService(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ObjectNotFoundError.ServiceNotFound
		}
		return nil, &errors.BadRequestError.ServiceNotCreated
	}

	return &schemas.Service{
		Id:          serviceModel.Id,
		Name:        serviceModel.Name,
		Description: serviceModel.Description,
		ImageUrl:    serviceModel.ImageUrl,
	}, nil
}

// Fetch services from postgresql DB and adapts them to a Service schema.
func (c *Service) FetchPostgresqlServices(ids []uuid.UUID) ([]*schemas.Service, *errors.Error) {
	servicesModel, err := c.DaoPostgresql.Service.FetchServices(ids)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.ServiceNotFound
	}

	services := make([]*schemas.Service, len(servicesModel))
	for i, serviceModel := range servicesModel {
		services[i] = &schemas.Service{
			Id:          serviceModel.Id,
			Name:        serviceModel.Name,
			Description: serviceModel.Description,
			ImageUrl:    serviceModel.ImageUrl,
			IsVirtual:   serviceModel.IsVirtual,
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

	// Validate name is not empty
	if name == "" {
		return nil, &errors.BadRequestError.InvalidServiceName
	}

	serviceModel := &model.Service{
		Id:          uuid.New(),
		Name:        name,
		Description: description,
		ImageUrl:    imageUrl,
		IsVirtual:   isVirtual,
		AuditFields: model.AuditFields{
			UpdatedBy: updatedBy,
		},
	}

	if err := c.DaoPostgresql.Service.CreateService(serviceModel); err != nil {
		return nil, &errors.BadRequestError.ServiceNotCreated
	}

	return &schemas.Service{
		Id:          serviceModel.Id,
		Name:        serviceModel.Name,
		Description: serviceModel.Description,
		ImageUrl:    serviceModel.ImageUrl,
		IsVirtual:   serviceModel.IsVirtual,
	}, nil
}

// Updates a service from a Postgresql DB given its ID and adapts it to a service schema.
func (s *Service) UpdatePostgresqlService(
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

	serviceModel, err := s.DaoPostgresql.Service.UpdateService(id, name, description, imageUrl, isVirtual, updatedBy)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ObjectNotFoundError.ServiceNotFound
		}
		return nil, &errors.BadRequestError.ServiceNotUpdated
	}

	return &schemas.Service{
		Id:          serviceModel.Id,
		Name:        serviceModel.Name,
		Description: serviceModel.Description,
		ImageUrl:    serviceModel.ImageUrl,
		IsVirtual:   serviceModel.IsVirtual,
	}, nil
}

// Soft deletes a service from a Postgresql DB given its ID.
func (s *Service) DeletePostgresqlService(id uuid.UUID) *errors.Error {
	err := s.DaoPostgresql.Service.DeleteService(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &errors.ObjectNotFoundError.ServiceNotFound
		}
		return &errors.BadRequestError.ServiceNotSoftDeleted
	}

	return nil
}

// Bulk deletes services from postgresql DB
func (s *Service) BulkDeletePostgresqlServices(serviceIds []string) *errors.Error {
	// Convert string IDs to UUIDs
	uuidIds := make([]uuid.UUID, len(serviceIds))
	for i, id := range serviceIds {
		parsedId, err := uuid.Parse(id)
		if err != nil {
			return &errors.UnprocessableEntityError.InvalidServiceId
		}
		uuidIds[i] = parsedId
	}

	if err := s.DaoPostgresql.Service.BulkDeleteServices(uuidIds); err != nil {
		return &errors.BadRequestError.ServiceNotSoftDeleted
	}

	return nil
}
