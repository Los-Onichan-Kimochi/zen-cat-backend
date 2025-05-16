package adapter

import (
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"
	daoPsql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"

	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
)

type CommunityService struct {
	logger        logging.Logger
	DaoPostgresql *daoPsql.AstroCatPsqlCollection
}

// Create CommunityService adapter
func NewCommunityServiceAdapter(
	logger logging.Logger,
	daoPostgresql *controller.AstroCatPsqlCollection,
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
		CommunityId: communityServiceModel.CommunityId,
		ServiceId:   communityServiceModel.ServiceId,
	}, nil
}

// Gets a specific community-service association and adapts it.
func (cs *CommunityService) GetPostgresqlCommunityService(
	communityId uuid.UUID,
	serviceId uuid.UUID,
) (*schemas.CommunityService, *errors.Error) {
	associationModel, err := cs.DaoPostgresql.CommunityService.GetCommunityService(
		communityId,
		serviceId,
	)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.CommunityServiceNotFound
	}

	return &schemas.CommunityService{
		CommunityId: associationModel.CommunityId,
		ServiceId:   associationModel.ServiceId,
	}, nil
}

// Deletes a specific community-service association from postgresql DB.
func (cs *CommunityService) DeletePostgresqlCommunityService(
	communityId uuid.UUID,
	serviceId uuid.UUID,
) *errors.Error {
	err := cs.DaoPostgresql.CommunityService.DeleteCommunityService(communityId, serviceId)
	if err != nil {
		return &errors.BadRequestError.CommunityServiceNotDeleted
	}

	return nil
}
