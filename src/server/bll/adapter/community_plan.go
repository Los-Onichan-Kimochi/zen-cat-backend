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
		Id:          uuid.New(),
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
		Id:          communityPlanModel.Id,
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
		Id:          associationModel.Id,
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

// Creates multiple community-plan associations.
func (cp *CommunityPlan) BulkCreatePostgresqlCommunityPlans(
	communityPlans []*schemas.CreateCommunityPlanRequest,
	updatedBy string,
) ([]*schemas.CommunityPlan, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	communityPlanModels := make([]*model.CommunityPlan, len(communityPlans))
	for i, communityPlan := range communityPlans {
		communityPlanModels[i] = &model.CommunityPlan{
			Id:          uuid.New(),
			CommunityId: communityPlan.CommunityId,
			PlanId:      communityPlan.PlanId,
			AuditFields: model.AuditFields{
				UpdatedBy: updatedBy,
			},
		}
	}
	err := cp.DaoPostgresql.CommunityPlan.BulkCreateCommunityPlans(communityPlanModels)
	if err != nil {
		if strings.Contains(err.Error(), "already exist") {
			return nil, &errors.ConflictError.CommunityPlanAlreadyExists
		}
		return nil, &errors.BadRequestError.CommunityPlanNotCreated
	}

	communityPlansResponse := make([]*schemas.CommunityPlan, len(communityPlans))
	for i, communityPlan := range communityPlanModels {
		communityPlansResponse[i] = &schemas.CommunityPlan{
			Id:          communityPlan.Id,
			CommunityId: communityPlan.CommunityId,
			PlanId:      communityPlan.PlanId,
		}
	}

	return communityPlansResponse, nil
}

// Fetch all community-plan associations from postgresql DB and adapts them to a CommunityPlan schema.
func (cp *CommunityPlan) FetchPostgresqlCommunityPlans(
	communityId *uuid.UUID,
	planId *uuid.UUID,
) ([]*schemas.CommunityPlan, *errors.Error) {
	communityPlanModels, err := cp.DaoPostgresql.CommunityPlan.FetchCommunityPlans(
		communityId,
		planId,
	)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.CommunityPlanNotFound
	}

	communityPlans := make([]*schemas.CommunityPlan, len(communityPlanModels))
	for i, communityPlan := range communityPlanModels {
		communityPlans[i] = &schemas.CommunityPlan{
			Id:          communityPlan.Id,
			CommunityId: communityPlan.CommunityId,
			PlanId:      communityPlan.PlanId,
		}
	}

	return communityPlans, nil
}

// Bulk deletes community-plan associations from postgresql DB.
func (cp *CommunityPlan) BulkDeletePostgresqlCommunityPlans(
	communityPlans []*schemas.CommunityPlan,
) *errors.Error {
	if len(communityPlans) == 0 {
		return nil
	}

	// Validate that all community-plan ids to delete are valid
	communityPlanModels := make([]*model.CommunityPlan, len(communityPlans))
	for i, communityPlan := range communityPlans {
		if communityPlan.CommunityId == uuid.Nil || communityPlan.PlanId == uuid.Nil {
			return &errors.UnprocessableEntityError.InvalidCommunityPlanId
		}

		communityPlanModels[i] = &model.CommunityPlan{
			Id:          communityPlan.Id,
			CommunityId: communityPlan.CommunityId,
			PlanId:      communityPlan.PlanId,
		}
	}

	if err := cp.DaoPostgresql.CommunityPlan.BulkDeleteCommunityPlans(communityPlanModels); err != nil {
		return &errors.BadRequestError.CommunityPlanNotDeleted
	}

	return nil
}
