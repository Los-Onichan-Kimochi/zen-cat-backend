package controller

import (
	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
	bllAdapter "onichankimochi.com/astro_cat_backend/src/server/bll/adapter"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type CommunityService struct {
	logger      logging.Logger
	Adapter     *bllAdapter.AdapterCollection
	EnvSettings *schemas.EnvSettings
}

// Create CommunityService controller
func NewCommunityServiceController(
	logger logging.Logger,
	adapter *bllAdapter.AdapterCollection,
	envSettings *schemas.EnvSettings,
) *CommunityService {
	return &CommunityService{
		logger:      logger,
		Adapter:     adapter,
		EnvSettings: envSettings,
	}
}

// Creates a community-service association.
func (cs *CommunityService) CreateCommunityService(
	req schemas.CreateCommunityServiceRequest,
	updatedBy string,
) (*schemas.CommunityService, *errors.Error) {
	communityId := req.CommunityId
	serviceId := req.ServiceId

	_, err := cs.Adapter.Community.GetPostgresqlCommunity(communityId)
	if err != nil {
		return nil, err
	}

	_, err = cs.Adapter.Service.GetPostgresqlService(serviceId)
	if err != nil {
		return nil, err
	}

	_, err = cs.Adapter.CommunityService.GetPostgresqlCommunityService(communityId, serviceId)
	if err == nil {
		return nil, &errors.ConflictError.CommunityServiceAlreadyExists
	} else if err.Code != errors.ObjectNotFoundError.CommunityServiceNotFound.Code {
		return nil, &errors.InternalServerError.Default
	}

	return cs.Adapter.CommunityService.CreatePostgresqlCommunityService(
		communityId,
		serviceId,
		updatedBy,
	)
}

// Gets a specific community-service association.
func (cs *CommunityService) GetCommunityService(
	communityIdString string,
	serviceIdString string,
) (*schemas.CommunityService, *errors.Error) {
	communityId, err := uuid.Parse(communityIdString)
	if err != nil {
		return nil, &errors.UnprocessableEntityError.InvalidCommunityId
	}

	serviceId, err := uuid.Parse(serviceIdString)
	if err != nil {
		return nil, &errors.UnprocessableEntityError.InvalidServiceId
	}

	return cs.Adapter.CommunityService.GetPostgresqlCommunityService(communityId, serviceId)
}

// Deletes a specific community-service association.
func (cs *CommunityService) DeleteCommunityService(
	communityIdString string,
	serviceIdString string,
) *errors.Error {
	communityId, parseErr := uuid.Parse(communityIdString)
	if parseErr != nil {
		return &errors.UnprocessableEntityError.InvalidCommunityId
	}

	serviceId, parseErrS := uuid.Parse(serviceIdString)
	if parseErrS != nil {
		return &errors.UnprocessableEntityError.InvalidServiceId
	}

	_, err := cs.Adapter.CommunityService.GetPostgresqlCommunityService(communityId, serviceId)
	if err != nil {
		return err
	}

	return cs.Adapter.CommunityService.DeletePostgresqlCommunityService(communityId, serviceId)
}
