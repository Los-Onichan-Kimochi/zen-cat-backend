package model

import "github.com/google/uuid"

type PlanType string

const (
	PlanTypeMonthly PlanType = "MONTHLY"
	PlanTypeAnual   PlanType = "ANUAL"
)

type Plan struct {
	Id               uuid.UUID `gorm:"type:uuid;primaryKey"`
	Fee              float64
	Type             PlanType
	ReservationLimit *int
	AuditFields
}

func (Plan) TableName() string {
	return "astro_cat_plan"
}
