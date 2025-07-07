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
		Id:             uuid.New(),
		ServiceId:      serviceId,
		ProfessionalId: professionalId,
		AuditFields: model.AuditFields{
			UpdatedBy: updatedBy,
		},
	}

	err := cp.DaoPostgresql.ServiceProfessional.CreateServiceProfessional(serviceProfessionalModel)
	if err != nil {
		return nil, &errors.BadRequestError.ServiceProfessionalNotCreated
	}

	return &schemas.ServiceProfessional{
		Id:             serviceProfessionalModel.Id,
		ServiceId:      serviceProfessionalModel.ServiceId,
		ProfessionalId: serviceProfessionalModel.ProfessionalId,
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
		Id:             associationModel.Id,
		ServiceId:      associationModel.ServiceId,
		ProfessionalId: associationModel.ProfessionalId,
	}, nil
}

// Todo: add comment
func (cp *ServiceProfessional) GetPostgresqlProfessionalsByServiceId(
	serviceId uuid.UUID,
) ([]*schemas.Professional, *errors.Error) {
	professionalsModel, err := cp.DaoPostgresql.ServiceProfessional.GetProfessionalsByServiceId(
		serviceId,
	)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.ServiceProfessionalNotFound
	}

	professionals := make([]*schemas.Professional, len(professionalsModel))
	for i, professionalModel := range professionalsModel {
		professionals[i] = &schemas.Professional{
			Id:             professionalModel.Id,
			Name:           professionalModel.Name,
			FirstLastName:  professionalModel.FirstLastName,
			SecondLastName: professionalModel.SecondLastName,
			Specialty:      professionalModel.Specialty,
			Email:          professionalModel.Email,
			PhoneNumber:    professionalModel.PhoneNumber,
			Type:           string(professionalModel.Type),
			ImageUrl:       professionalModel.ImageUrl,
		}
	}

	return professionals, nil
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

// Creates multiple service-professional associations.
func (cp *ServiceProfessional) BulkCreatePostgresqlServiceProfessionals(
	serviceProfessionals []*schemas.CreateServiceProfessionalRequest,
	updatedBy string,
) ([]*schemas.ServiceProfessional, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	serviceProfessionalModels := make([]*model.ServiceProfessional, len(serviceProfessionals))
	for i, serviceProfessional := range serviceProfessionals {
		serviceProfessionalModels[i] = &model.ServiceProfessional{
			Id:             uuid.New(),
			ServiceId:      serviceProfessional.ServiceId,
			ProfessionalId: serviceProfessional.ProfessionalId,
			AuditFields: model.AuditFields{
				UpdatedBy: updatedBy,
			},
		}
	}
	err := cp.DaoPostgresql.ServiceProfessional.BulkCreateServiceProfessionals(serviceProfessionalModels)
	if err != nil {
		if strings.Contains(err.Error(), "already exist") {
			return nil, &errors.ConflictError.ServiceProfessionalAlreadyExists
		}
		return nil, &errors.BadRequestError.ServiceProfessionalNotCreated
	}

	serviceProfessionalsResponse := make([]*schemas.ServiceProfessional, len(serviceProfessionals))
	for i, serviceProfessional := range serviceProfessionalModels {
		serviceProfessionalsResponse[i] = &schemas.ServiceProfessional{
			Id:             serviceProfessional.Id,
			ServiceId:      serviceProfessional.ServiceId,
			ProfessionalId: serviceProfessional.ProfessionalId,
		}
	}

	return serviceProfessionalsResponse, nil
}

// Fetch all service-professional associations from postgresql DB and adapts them to a ServiceProfessional schema.
func (cp *ServiceProfessional) FetchPostgresqlServiceProfessionals(
	serviceId *uuid.UUID,
	professionalId *uuid.UUID,
) ([]*schemas.ServiceProfessional, *errors.Error) {
	serviceProfessionalModels, err := cp.DaoPostgresql.ServiceProfessional.FetchServiceProfessionals(
		serviceId,
		professionalId,
	)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.ServiceProfessionalNotFound
	}

	serviceProfessionals := make([]*schemas.ServiceProfessional, len(serviceProfessionalModels))
	for i, serviceProfessional := range serviceProfessionalModels {
		serviceProfessionals[i] = &schemas.ServiceProfessional{
			Id:             serviceProfessional.Id,
			ServiceId:      serviceProfessional.ServiceId,
			ProfessionalId: serviceProfessional.ProfessionalId,
		}
	}

	return serviceProfessionals, nil
}

// Bulk deletes service-professional associations from postgresql DB.
func (cp *ServiceProfessional) BulkDeletePostgresqlServiceProfessionals(
	serviceProfessionals []*schemas.DeleteServiceProfessionalRequest,
) *errors.Error {
	if len(serviceProfessionals) == 0 {
		return nil
	}

	// Validate that all service-professional ids to delete are valid
	serviceProfessionalModels := make([]*model.ServiceProfessional, len(serviceProfessionals))
	for i, serviceProfessional := range serviceProfessionals {
		if serviceProfessional.ServiceId == uuid.Nil || serviceProfessional.ProfessionalId == uuid.Nil {
			return &errors.UnprocessableEntityError.InvalidServiceProfessionalId
		}

		serviceProfessionalModels[i] = &model.ServiceProfessional{
			ServiceId:      serviceProfessional.ServiceId,
			ProfessionalId: serviceProfessional.ProfessionalId,
		}
	}

	if err := cp.DaoPostgresql.ServiceProfessional.BulkDeleteServiceProfessionals(serviceProfessionalModels); err != nil {
		return &errors.BadRequestError.ServiceProfessionalNotDeleted
	}

	return nil
}
