package adapter

import (
	daoPsql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"

	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
)

type CommunityPlan struct {
	logger        logging.Logger
	DaoPostgresql *daoPsql.AstroCatPsqlCollection
}

// Create CommunityPlan adapter
func NewCommunityPlanAdapter(
	logger logging.Logger,
	daoPostgresql *daoPsql.AstroCatPsqlCollection,
) *CommunityPlan {
	return &CommunityPlan{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

// Creates a community-plan association into postgresql DB.
func (cp *CommunityPlan) CreatePostgresqlCommunityPlan(
	communityId uuid.UUID,
	planId uuid.UUID,
	updatedBy string,
) (*schemas.CommunityPlan, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	communityPlanModel := &model.CommunityPlan{
		CommunityId: communityId,
		PlanId:      planId,
		AuditFields: model.AuditFields{
			UpdatedBy: updatedBy,
		},
	}

	err := cp.DaoPostgresql.CommunityPlan.CreateCommunityPlan(communityPlanModel)
	if err != nil {
		return nil, &errors.BadRequestError.CommunityPlanNotCreated
	}

	return &schemas.CommunityPlan{
		CommunityId: communityPlanModel.CommunityId,
		PlanId:      communityPlanModel.PlanId,
	}, nil
}

// Gets a specific community-plan association and adapts it.
func (cp *CommunityPlan) GetPostgresqlCommunityPlan(
	communityId uuid.UUID,
	planId uuid.UUID,
) (*schemas.CommunityPlan, *errors.Error) {
	associationModel, err := cp.DaoPostgresql.CommunityPlan.GetCommunityPlan(communityId, planId)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.CommunityPlanNotFound
	}

	return &schemas.CommunityPlan{
		CommunityId: associationModel.CommunityId,
		PlanId:      associationModel.PlanId,
	}, nil
}

// Deletes a specific community-plan association from postgresql DB.
func (cp *CommunityPlan) DeletePostgresqlCommunityPlan(
	communityId uuid.UUID,
	planId uuid.UUID,
) *errors.Error {
	err := cp.DaoPostgresql.CommunityPlan.DeleteCommunityPlan(communityId, planId)
	if err != nil {
		return &errors.BadRequestError.CommunityPlanNotDeleted
	}

	return nil
}
