package adapter

import (
	daoPsql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"

	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
)

type ServiceProfessional struct {
	logger        logging.Logger
	DaoPostgresql *daoPsql.AstroCatPsqlCollection
}

// Create ServiceProfessional adapter
func NewServiceProfessionalAdapter(
	logger logging.Logger,
	daoPostgresql *daoPsql.AstroCatPsqlCollection,
) *ServiceProfessional {
	return &ServiceProfessional{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

// Creates a service-professional association into postgresql DB.
func (cp *ServiceProfessional) CreatePostgresqlServiceProfessional(
	serviceId uuid.UUID,
	professionalId uuid.UUID,
	updatedBy string,
) (*schemas.ServiceProfessional, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	serviceProfessionalModel := &model.ServiceProfessional{
		ServiceId: serviceId,
		ProfessionalId:      professionalId,
		AuditFields: model.AuditFields{
			UpdatedBy: updatedBy,
		},
	}

	err := cp.DaoPostgresql.ServiceProfessional.CreateServiceProfessional(serviceProfessionalModel)
	if err != nil {
		return nil, &errors.BadRequestError.ServiceProfessionalNotCreated
	}

	return &schemas.ServiceProfessional{
		ServiceId: serviceProfessionalModel.ServiceId,
		ProfessionalId:      serviceProfessionalModel.ProfessionalId,
	}, nil
}

// Gets a specific service-professional association and adapts it.
func (cp *ServiceProfessional) GetPostgresqlServiceProfessional(
	serviceId uuid.UUID,
	professionalId uuid.UUID,
) (*schemas.ServiceProfessional, *errors.Error) {
	associationModel, err := cp.DaoPostgresql.ServiceProfessional.GetServiceProfessional(serviceId, professionalId)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.ServiceProfessionalNotFound
	}

	return &schemas.ServiceProfessional{
		ServiceId: associationModel.ServiceId,
		ProfessionalId:      associationModel.ProfessionalId,
	}, nil
}

// Deletes a specific service-professional association from postgresql DB.
func (cp *ServiceProfessional) DeletePostgresqlServiceProfessional(
	serviceId uuid.UUID,
	professionalId uuid.UUID,
) *errors.Error {
	err := cp.DaoPostgresql.ServiceProfessional.DeleteServiceProfessional(serviceId, professionalId)
	if err != nil {
		return &errors.BadRequestError.ServiceProfessionalNotDeleted
	}

	return nil
}
