package controller

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type Plan struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

// Create Plan postgresql controller
func NewPlanController(logger logging.Logger, postgresqlDB *gorm.DB) *Plan {
	return &Plan{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

// Gets a plan model given its ID.
func (p *Plan) GetPlan(planId uuid.UUID) (*model.Plan, error) {
	plan := &model.Plan{}

	result := p.PostgresqlDB.First(&plan, "id = ?", planId)
	if result.Error != nil {
		return nil, result.Error
	}

	return plan, nil
}

// Fetches all plans.
// TODO: Add filters and sorting.
func (p *Plan) FetchPlans() ([]*model.Plan, error) {
	plans := []*model.Plan{}

	result := p.PostgresqlDB.Find(&plans)
	if result.Error != nil {
		return nil, result.Error
	}

	return plans, nil
}

// Creates a plan given its model.
func (p *Plan) CreatePlan(plan *model.Plan) error {
	return p.PostgresqlDB.Create(plan).Error
}

// Updates a plan given fields to update.
func (p *Plan) UpdatePlan(
	id uuid.UUID,
	fee *float64,
	planType *model.PlanType,
	reservationLimit *int,
	updatedBy string,
) (*model.Plan, error) {
	updateFields := map[string]any{
		"updated_by": updatedBy,
	}

	if fee != nil {
		updateFields["fee"] = *fee
	}
	if planType != nil {
		updateFields["type"] = *planType
	}
	if reservationLimit != nil {
		updateFields["reservation_limit"] = *reservationLimit
	}

	var plan model.Plan
	// Check if there are any fields to update other than updated_by
	if len(updateFields) == 1 {
		if err := p.PostgresqlDB.First(&plan, "id = ?", id).Error; err != nil {
			return nil, err
		}
		return &plan, nil // No actual fields to update, return current record
	}

	result := p.PostgresqlDB.Model(&plan).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(updateFields)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound // Or handle as a non-error if appropriate
	}

	return &plan, nil
}

// Soft deletes a plan given its ID.
func (p *Plan) DeletePlan(planId uuid.UUID) error {
	result := p.PostgresqlDB.Delete(&model.Plan{}, "id = ?", planId)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
