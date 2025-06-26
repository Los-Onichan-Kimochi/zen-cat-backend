package adapter

import (
	"strings"

	daoPsql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
)

type CommunityService struct {
	logger        logging.Logger
	DaoPostgresql *daoPsql.AstroCatPsqlCollection
}

// Create CommunityService adapter
func NewCommunityServiceAdapter(
	logger logging.Logger,
	daoPostgresql *daoPsql.AstroCatPsqlCollection,
) *CommunityService {
	return &CommunityService{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

// Creates a community-service association in the PostgreSQL DB.
func (cs *CommunityService) CreatePostgresqlCommunityService(
	communityId uuid.UUID,
	serviceId uuid.UUID,
	updatedBy string,
) (*schemas.CommunityService, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	communityServiceModel := &model.CommunityService{
		Id:          uuid.New(),
		CommunityId: communityId,
		ServiceId:   serviceId,
		AuditFields: model.AuditFields{
			UpdatedBy: updatedBy,
		},
	}

	err := cs.DaoPostgresql.CommunityService.CreateCommunityService(communityServiceModel)
	if err != nil {
		return nil, &errors.BadRequestError.CommunityServiceNotCreated
	}

	return &schemas.CommunityService{
		Id:          communityServiceModel.Id,
		CommunityId: communityServiceModel.CommunityId,
		ServiceId:   communityServiceModel.ServiceId,
	}, nil
}

// Gets a specific community-service association and adapts it.
func (cs *CommunityService) GetPostgresqlCommunityService(
	communityId uuid.UUID,
	serviceId uuid.UUID,
) (*schemas.CommunityService, *errors.Error) {
	communityServiceModel, err := cs.DaoPostgresql.CommunityService.GetCommunityService(
		communityId,
		serviceId,
	)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.CommunityServiceNotFound
	}

	return &schemas.CommunityService{
		Id:          communityServiceModel.Id,
		CommunityId: communityServiceModel.CommunityId,
		ServiceId:   communityServiceModel.ServiceId,
	}, nil
}

// Todo: add comment
func (cs *CommunityService) GetPostgresqlServicesByCommunityId(
	communityId uuid.UUID,
) ([]*schemas.Service, *errors.Error) {
	servicesModel, err := cs.DaoPostgresql.CommunityService.GetServicesByCommunityId(
		communityId,
	)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.CommunityServiceNotFound
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

// Deletes a specific community-service association from postgresql DB.
func (cs *CommunityService) DeletePostgresqlCommunityService(
	communityId uuid.UUID,
	serviceId uuid.UUID,
) *errors.Error {
	err := cs.DaoPostgresql.CommunityService.DeleteCommunityService(communityId, serviceId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &errors.ObjectNotFoundError.CommunityServiceNotFound
		}
		return &errors.BadRequestError.CommunityServiceNotDeleted
	}

	return nil
}

// Creates multiple community-service associations.
func (cs *CommunityService) BulkCreatePostgresqlCommunityServices(
	communityServices []*schemas.CreateCommunityServiceRequest,
	updatedBy string,
) ([]*schemas.CommunityService, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	communityServiceModels := make([]*model.CommunityService, len(communityServices))
	for i, communityService := range communityServices {
		communityServiceModels[i] = &model.CommunityService{
			Id:          uuid.New(),
			CommunityId: communityService.CommunityId,
			ServiceId:   communityService.ServiceId,
			AuditFields: model.AuditFields{
				UpdatedBy: updatedBy,
			},
		}
	}
	err := cs.DaoPostgresql.CommunityService.BulkCreateCommunityServices(communityServiceModels)
	if err != nil {
		if strings.Contains(err.Error(), "already exist") {
			return nil, &errors.ConflictError.CommunityServiceAlreadyExists
		}
		return nil, &errors.BadRequestError.CommunityServiceNotCreated
	}

	communityServicesResponse := make([]*schemas.CommunityService, len(communityServices))
	for i, communityService := range communityServiceModels {
		communityServicesResponse[i] = &schemas.CommunityService{
			Id:          communityService.Id,
			CommunityId: communityService.CommunityId,
			ServiceId:   communityService.ServiceId,
		}
	}

	return communityServicesResponse, nil
}

// Bulk deletes community-service associations from postgresql DB.
func (cs *CommunityService) BulkDeletePostgresqlCommunityServices(
	communityServices []*schemas.DeleteCommunityServiceRequest,
) *errors.Error {
	if len(communityServices) == 0 {
		return nil
	}

	// Validate that all community-service ids to delete are valid
	communityServiceModels := make([]*model.CommunityService, len(communityServices))
	for i, communityService := range communityServices {
		if communityService.CommunityId == uuid.Nil || communityService.ServiceId == uuid.Nil {
			return &errors.UnprocessableEntityError.InvalidCommunityServiceId
		}

		communityServiceModels[i] = &model.CommunityService{
			CommunityId: communityService.CommunityId,
			ServiceId:   communityService.ServiceId,
		}
	}

	if err := cs.DaoPostgresql.CommunityService.BulkDeleteCommunityServices(communityServiceModels); err != nil {
		return &errors.BadRequestError.CommunityServiceNotDeleted
	}

	return nil
}

// Fetch all community-service associations from postgresql DB and adapts them to a CommunityService schema.
func (cs *CommunityService) FetchPostgresqlCommunityServices(
	communityId *uuid.UUID,
	serviceId *uuid.UUID,
) ([]*schemas.CommunityService, *errors.Error) {
	communityServiceModels, err := cs.DaoPostgresql.CommunityService.FetchCommunityServices(
		communityId,
		serviceId,
	)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.CommunityServiceNotFound
	}

	communityServices := make([]*schemas.CommunityService, len(communityServiceModels))
	for i, communityService := range communityServiceModels {
		communityServices[i] = &schemas.CommunityService{
			Id:          communityService.Id,
			CommunityId: communityService.CommunityId,
			ServiceId:   communityService.ServiceId,
		}
	}

	return communityServices, nil
}
