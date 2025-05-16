package controller

import (

	// "gorm.io/gorm" // No longer directly needed here for these checks
	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
	bllAdapter "onichankimochi.com/astro_cat_backend/src/server/bll/adapter"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type CommunityPlan struct {
	logger      logging.Logger
	Adapter     *bllAdapter.AdapterCollection
	EnvSettings *schemas.EnvSettings
}

// Create CommunityPlan controller
func NewCommunityPlanController(
	logger logging.Logger,
	adapter *bllAdapter.AdapterCollection,
	envSettings *schemas.EnvSettings,
) *CommunityPlan {
	return &CommunityPlan{
		logger:      logger,
		Adapter:     adapter,
		EnvSettings: envSettings,
	}
}

// Creates a community-plan association.
func (cp *CommunityPlan) CreateCommunityPlan(
	req schemas.CreateCommunityPlanRequest,
	updatedBy string,
) (*schemas.CommunityPlan, *errors.Error) {
	communityId := req.CommunityId
	planId := req.PlanId

	_, err := cp.Adapter.Community.GetPostgresqlCommunity(communityId)
	if err != nil {
		return nil, err
	}

	_, err = cp.Adapter.Plan.GetPostgresqlPlan(planId)
	if err != nil {
		return nil, err
	}

	_, err = cp.Adapter.CommunityPlan.GetPostgresqlCommunityPlan(communityId, planId)
	if err == nil {
		return nil, &errors.BadRequestError.CommunityPlanAlreadyExists
	} else if err.Code != errors.ObjectNotFoundError.CommunityPlanNotFound.Code {
		return nil, &errors.InternalServerError.Default
	}

	return cp.Adapter.CommunityPlan.CreatePostgresqlCommunityPlan(communityId, planId, updatedBy)
}

// Gets a specific community-plan association.
func (cp *CommunityPlan) GetCommunityPlan(
	communityIdString string,
	planIdString string,
) (*schemas.CommunityPlan, *errors.Error) {
	communityId, err := uuid.Parse(communityIdString)
	if err != nil {
		return nil, &errors.UnprocessableEntityError.InvalidCommunityId
	}

	planId, err := uuid.Parse(planIdString)
	if err != nil {
		return nil, &errors.UnprocessableEntityError.InvalidPlanId
	}

	return cp.Adapter.CommunityPlan.GetPostgresqlCommunityPlan(communityId, planId)
}

// Deletes a specific community-plan association.
func (cp *CommunityPlan) DeleteCommunityPlan(
	communityIdString string,
	planIdString string,
) *errors.Error {
	communityId, parseErr := uuid.Parse(communityIdString)
	if parseErr != nil {
		return &errors.UnprocessableEntityError.InvalidCommunityId
	}

	planId, parseErr := uuid.Parse(planIdString)
	if parseErr != nil {
		return &errors.UnprocessableEntityError.InvalidPlanId
	}

	_, err := cp.Adapter.CommunityPlan.GetPostgresqlCommunityPlan(communityId, planId)
	if err != nil {
		return err
	}

	return cp.Adapter.CommunityPlan.DeletePostgresqlCommunityPlan(communityId, planId)
}
