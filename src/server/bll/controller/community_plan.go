package controller

import (

	// "gorm.io/gorm" // No longer directly needed here for these checks
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
