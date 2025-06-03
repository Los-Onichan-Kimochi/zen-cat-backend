package adapter

import (
	"strings"

	daoPsql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"

	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
)

type ServiceLocal struct {
	logger        logging.Logger
	DaoPostgresql *daoPsql.AstroCatPsqlCollection
}

// Create ServiceLocal adapter
func NewServiceLocalAdapter(
	logger logging.Logger,
	daoPostgresql *daoPsql.AstroCatPsqlCollection,
) *ServiceLocal {
	return &ServiceLocal{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

// Creates a service-local association into postgresql DB.
func (cp *ServiceLocal) CreatePostgresqlServiceLocal(
	serviceId uuid.UUID,
	localId uuid.UUID,
	updatedBy string,
) (*schemas.ServiceLocal, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	serviceLocalModel := &model.ServiceLocal{
		Id:          uuid.New(),
		ServiceId: serviceId,
		LocalId:      localId,
		AuditFields: model.AuditFields{
			UpdatedBy: updatedBy,
		},
	}

	err := cp.DaoPostgresql.ServiceLocal.CreateServiceLocal(serviceLocalModel)
	if err != nil {
		return nil, &errors.BadRequestError.ServiceLocalNotCreated
	}

	return &schemas.ServiceLocal{
		Id:          serviceLocalModel.Id,
		ServiceId: serviceLocalModel.ServiceId,
		LocalId:      serviceLocalModel.LocalId,
	}, nil
}

// Gets a specific service-local association and adapts it.
func (cp *ServiceLocal) GetPostgresqlServiceLocal(
	serviceId uuid.UUID,
	localId uuid.UUID,
) (*schemas.ServiceLocal, *errors.Error) {
	associationModel, err := cp.DaoPostgresql.ServiceLocal.GetServiceLocal(serviceId, localId)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.ServiceLocalNotFound
	}

	return &schemas.ServiceLocal{
		Id:          associationModel.Id,
		ServiceId: associationModel.ServiceId,
		LocalId:      associationModel.LocalId,
	}, nil
}

// Deletes a specific service-local association from postgresql DB.
func (cp *ServiceLocal) DeletePostgresqlServiceLocal(
	serviceId uuid.UUID,
	localId uuid.UUID,
) *errors.Error {
	err := cp.DaoPostgresql.ServiceLocal.DeleteServiceLocal(serviceId, localId)
	if err != nil {
		return &errors.BadRequestError.ServiceLocalNotDeleted
	}

	return nil
}

// Creates multiple service-local associations.
func (cp *ServiceLocal) BulkCreatePostgresqlServiceLocals(
	serviceLocals []*schemas.CreateServiceLocalRequest,
	updatedBy string,
) ([]*schemas.ServiceLocal, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	serviceLocalModels := make([]*model.ServiceLocal, len(serviceLocals))
	for i, serviceLocal := range serviceLocals {
		serviceLocalModels[i] = &model.ServiceLocal{
			Id:          uuid.New(),
			ServiceId: serviceLocal.ServiceId,
			LocalId:      serviceLocal.LocalId,
			AuditFields: model.AuditFields{
				UpdatedBy: updatedBy,
			},
		}
	}
	err := cp.DaoPostgresql.ServiceLocal.BulkCreateServiceLocals(serviceLocalModels)
	if err != nil {
		if strings.Contains(err.Error(), "already exist") {
			return nil, &errors.ConflictError.ServiceLocalAlreadyExists
		}
		return nil, &errors.BadRequestError.ServiceLocalNotCreated
	}

	serviceLocalsResponse := make([]*schemas.ServiceLocal, len(serviceLocals))
	for i, serviceLocal := range serviceLocalModels {
		serviceLocalsResponse[i] = &schemas.ServiceLocal{
			Id:          serviceLocal.Id,
			ServiceId: serviceLocal.ServiceId,
			LocalId:      serviceLocal.LocalId,
		}
	}

	return serviceLocalsResponse, nil
}

// Fetch all service-local associations from postgresql DB and adapts them to a ServiceLocal schema.
func (cp *ServiceLocal) FetchPostgresqlServiceLocals(
	serviceId *uuid.UUID,
	localId *uuid.UUID,
) ([]*schemas.ServiceLocal, *errors.Error) {
	serviceLocalModels, err := cp.DaoPostgresql.ServiceLocal.FetchServiceLocals(
		serviceId,
		localId,
	)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.ServiceLocalNotFound
	}

	serviceLocals := make([]*schemas.ServiceLocal, len(serviceLocalModels))
	for i, serviceLocal := range serviceLocalModels {
		serviceLocals[i] = &schemas.ServiceLocal{
			Id:          serviceLocal.Id,
			ServiceId: serviceLocal.ServiceId,
			LocalId:      serviceLocal.LocalId,
		}
	}

	return serviceLocals, nil
}

// Bulk deletes service-local associations from postgresql DB.
func (cp *ServiceLocal) BulkDeletePostgresqlServiceLocals(
	serviceLocals []*schemas.DeleteServiceLocalRequest,
) *errors.Error {
	if len(serviceLocals) == 0 {
		return nil
	}

	// Validate that all service-local ids to delete are valid
	serviceLocalModels := make([]*model.ServiceLocal, len(serviceLocals))
	for i, serviceLocal := range serviceLocals {
		if serviceLocal.ServiceId == uuid.Nil || serviceLocal.LocalId == uuid.Nil {
			return &errors.UnprocessableEntityError.InvalidServiceLocalId
		}

		serviceLocalModels[i] = &model.ServiceLocal{
			ServiceId: serviceLocal.ServiceId,
			LocalId:      serviceLocal.LocalId,
		}
	}

	if err := cp.DaoPostgresql.ServiceLocal.BulkDeleteServiceLocals(serviceLocalModels); err != nil {
		return &errors.BadRequestError.ServiceLocalNotDeleted
	}

	return nil
}

