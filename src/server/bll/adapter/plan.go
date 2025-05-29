package adapter

import (
	"github.com/google/uuid"
	daoPsql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"

	"onichankimochi.com/astro_cat_backend/src/logging"
)

type Plan struct {
	logger        logging.Logger
	DaoPostgresql *daoPsql.AstroCatPsqlCollection
}

// Creates Plan adapter
func NewPlanAdapter(
	logger logging.Logger,
	daoPostgresql *daoPsql.AstroCatPsqlCollection,
) *Plan {
	return &Plan{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

// Gets a plan from postgresql DB and adapts it to a Plan schema.
func (p *Plan) GetPostgresqlPlan(planId uuid.UUID) (*schemas.Plan, *errors.Error) {
	planModel, err := p.DaoPostgresql.Plan.GetPlan(planId)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.PlanNotFound
	}

	return &schemas.Plan{
		Id:               planModel.Id,
		Fee:              planModel.Fee,
		Type:             planModel.Type,
		ReservationLimit: planModel.ReservationLimit,
	}, nil
}

// Fetches plans from postgresql DB and adapts them to Plan schemas.
func (p *Plan) FetchPostgresqlPlans(ids []uuid.UUID) ([]*schemas.Plan, *errors.Error) {
	planModels, err := p.DaoPostgresql.Plan.FetchPlans(ids)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.PlanNotFound
	}

	plans := make([]*schemas.Plan, len(planModels))
	for i, planModel := range planModels {
		plans[i] = &schemas.Plan{
			Id:               planModel.Id,
			Fee:              planModel.Fee,
			Type:             planModel.Type,
			ReservationLimit: planModel.ReservationLimit,
		}
	}

	return plans, nil
}

// Creates a plan in postgresql DB and returns it as a Plan schema.
func (p *Plan) CreatePostgresqlPlan(
	fee float64,
	planType model.PlanType,
	reservationLimit *int,
	updatedBy string,
) (*schemas.Plan, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	planModel := &model.Plan{
		Id:               uuid.New(),
		Fee:              fee,
		Type:             planType,
		ReservationLimit: reservationLimit,
		AuditFields: model.AuditFields{
			UpdatedBy: updatedBy,
		},
	}

	if err := p.DaoPostgresql.Plan.CreatePlan(planModel); err != nil {
		return nil, &errors.BadRequestError.PlanNotCreated
	}

	return &schemas.Plan{
		Id:               planModel.Id,
		Fee:              planModel.Fee,
		Type:             planModel.Type,
		ReservationLimit: planModel.ReservationLimit,
	}, nil
}

// Updates a plan in postgresql DB and returns it as a Plan schema.
func (p *Plan) UpdatePostgresqlPlan(
	planId uuid.UUID,
	fee *float64,
	planType *model.PlanType,
	reservationLimit *int,
	updatedBy string,
) (*schemas.Plan, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	planModel, err := p.DaoPostgresql.Plan.UpdatePlan(
		planId,
		fee,
		planType,
		reservationLimit,
		updatedBy,
	)
	if err != nil {
		return nil, &errors.BadRequestError.PlanNotUpdated
	}

	return &schemas.Plan{
		Id:               planModel.Id,
		Fee:              planModel.Fee,
		Type:             planModel.Type,
		ReservationLimit: planModel.ReservationLimit,
	}, nil
}

// Deletes a plan from postgresql DB.
func (p *Plan) DeletePostgresqlPlan(planId uuid.UUID) *errors.Error {
	err := p.DaoPostgresql.Plan.DeletePlan(planId)
	if err != nil {
		return &errors.BadRequestError.PlanNotSoftDeleted
	}

	return nil
}
