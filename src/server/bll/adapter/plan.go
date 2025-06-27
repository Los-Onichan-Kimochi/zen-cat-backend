package adapter

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
	daoPsql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
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

// Gets a plan from a Postgresql DB given its ID and adapts it to a plan schema.
func (p *Plan) GetPostgresqlPlan(id uuid.UUID) (*schemas.Plan, *errors.Error) {
	planModel, err := p.DaoPostgresql.Plan.GetPlan(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ObjectNotFoundError.PlanNotFound
		}
		return nil, &errors.BadRequestError.PlanNotCreated
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

// Bulk creates plans into postgresql DB.
func (p *Plan) BulkCreatePostgresqlPlans(
	plansData []*schemas.CreatePlanRequest,
	updatedBy string,
) ([]*schemas.Plan, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	plansModel := make([]*model.Plan, len(plansData))
	for i, planData := range plansData {
		plansModel[i] = &model.Plan{
			Id:               uuid.New(),
			Fee:              planData.Fee,
			Type:             planData.Type,
			ReservationLimit: planData.ReservationLimit,
			AuditFields: model.AuditFields{
				UpdatedBy: updatedBy,
			},
		}
	}
	if err := p.DaoPostgresql.Plan.BulkCreatePlans(plansModel); err != nil {
		return nil, &errors.BadRequestError.PlanNotCreated
	}

	plans := make([]*schemas.Plan, len(plansModel))
	for i, planModel := range plansModel {
		plans[i] = &schemas.Plan{
			Id:               planModel.Id,
			Fee:              planModel.Fee,
			Type:             planModel.Type,
			ReservationLimit: planModel.ReservationLimit,
		}
	}

	return plans, nil
}

// Updates a plan from a Postgresql DB given its ID and adapts it to a plan schema.
func (p *Plan) UpdatePostgresqlPlan(
	id uuid.UUID,
	fee *float64,
	planType *model.PlanType,
	reservationLimit *int,
	updatedBy string,
) (*schemas.Plan, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	planModel, err := p.DaoPostgresql.Plan.UpdatePlan(id, fee, planType, reservationLimit, updatedBy)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ObjectNotFoundError.PlanNotFound
		}
		return nil, &errors.BadRequestError.PlanNotUpdated
	}

	return &schemas.Plan{
		Id:               planModel.Id,
		Fee:              planModel.Fee,
		Type:             planModel.Type,
		ReservationLimit: planModel.ReservationLimit,
	}, nil
}

// Soft deletes a plan from a Postgresql DB given its ID.
func (p *Plan) DeletePostgresqlPlan(id uuid.UUID) *errors.Error {
	err := p.DaoPostgresql.Plan.DeletePlan(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &errors.ObjectNotFoundError.PlanNotFound
		}
		return &errors.BadRequestError.PlanNotSoftDeleted
	}

	return nil
}

// Bulk deletes plans from postgresql DB.
func (p *Plan) BulkDeletePostgresqlPlans(
	planIds []uuid.UUID,
) *errors.Error {
	if err := p.DaoPostgresql.Plan.BulkDeletePlans(planIds); err != nil {
		return &errors.BadRequestError.PlanNotSoftDeleted
	}

	return nil
}
