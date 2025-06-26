package factories

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type PlanModelF struct {
	Id               *uuid.UUID
	Fee              *float64
	Type             *model.PlanType
	ReservationLimit *int
}

// Create a new plan on DB
func NewPlanModel(db *gorm.DB, option ...PlanModelF) *model.Plan {
	reservationLimit := 10
	plan := &model.Plan{
		Id:               uuid.New(),
		Fee:              99.99,
		Type:             model.PlanTypeMonthly,
		ReservationLimit: &reservationLimit,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}

	if len(option) > 0 {
		parameters := option[0]
		if parameters.Id != nil {
			plan.Id = *parameters.Id
		}
		if parameters.Fee != nil {
			plan.Fee = *parameters.Fee
		}
		if parameters.Type != nil {
			plan.Type = *parameters.Type
		}
		if parameters.ReservationLimit != nil {
			plan.ReservationLimit = parameters.ReservationLimit
		}
	}

	result := db.Create(plan)
	if result.Error != nil {
		log.Fatalf("Error when trying to create plan: %v", result.Error)
	}

	return plan
}

// Create size number of new plans on DB
func NewPlanModelBatch(
	db *gorm.DB,
	size int,
	option ...PlanModelF,
) []*model.Plan {
	plans := []*model.Plan{}
	for i := 0; i < size; i++ {
		var plan *model.Plan
		if len(option) > 0 {
			plan = NewPlanModel(db, option[0])
		} else {
			plan = NewPlanModel(db)
		}
		plans = append(plans, plan)
	}
	return plans
}
