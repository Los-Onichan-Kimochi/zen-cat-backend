package controller

import (
	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
	bllAdapter "onichankimochi.com/astro_cat_backend/src/server/bll/adapter"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type Plan struct {
	logger      logging.Logger
	Adapter     *bllAdapter.AdapterCollection
	EnvSettings *schemas.EnvSettings
}

// Create Plan controller
func NewPlanController(
	logger logging.Logger,
	adapter *bllAdapter.AdapterCollection,
	envSettings *schemas.EnvSettings,
) *Plan {
	return &Plan{
		logger:      logger,
		Adapter:     adapter,
		EnvSettings: envSettings,
	}
}

// Gets a plan.
func (p *Plan) GetPlan(planId uuid.UUID) (*schemas.Plan, *errors.Error) {
	return p.Adapter.Plan.GetPostgresqlPlan(planId)
}

// Fetch all plans.
func (p *Plan) FetchPlans() (*schemas.Plans, *errors.Error) {
	plans, err := p.Adapter.Plan.FetchPostgresqlPlans()
	if err != nil {
		return nil, err
	}

	return &schemas.Plans{Plans: plans}, nil
}

// Creates a plan.
func (p *Plan) CreatePlan(
	createPlanData schemas.CreatePlanRequest,
	updatedBy string,
) (*schemas.Plan, *errors.Error) {
	if createPlanData.Type == "" {
		return nil, &errors.BadRequestError.InvalidPlanType
	}

	return p.Adapter.Plan.CreatePostgresqlPlan(
		createPlanData.Fee,
		createPlanData.Type,
		createPlanData.ReservationLimit,
		updatedBy,
	)
}

// Updates a plan.
func (p *Plan) UpdatePlan(
	planId uuid.UUID,
	updatePlanData schemas.UpdatePlanRequest,
	updatedBy string,
) (*schemas.Plan, *errors.Error) {
	return p.Adapter.Plan.UpdatePostgresqlPlan(
		planId,
		updatePlanData.Fee,
		updatePlanData.Type,
		updatePlanData.ReservationLimit,
		updatedBy,
	)
}

// Deletes a plan.
func (p *Plan) DeletePlan(planId uuid.UUID) *errors.Error {
	return p.Adapter.Plan.DeletePostgresqlPlan(planId)
}

// TODO: Add BulkCreatePlans (Batch)
