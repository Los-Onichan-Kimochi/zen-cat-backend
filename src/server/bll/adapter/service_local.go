package adapter

import (
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
