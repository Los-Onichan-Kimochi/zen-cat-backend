package schemas

import (
	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type Plan struct {
	Id               uuid.UUID      `json:"id"`
	Fee              float64        `json:"fee"`
	Type             model.PlanType `json:"type"`
	ReservationLimit *int           `json:"reservation_limit"`
}

type Plans struct {
	Plans []*Plan `json:"plans"`
	// TODO: Add pagination if needed
}

type CreatePlanRequest struct {
	Fee              float64        `json:"fee"`
	Type             model.PlanType `json:"type"`
	ReservationLimit *int           `json:"reservation_limit"`
}

type UpdatePlanRequest struct {
	Fee              *float64        `json:"fee"`
	Type             *model.PlanType `json:"type"`
	ReservationLimit *int            `json:"reservation_limit"`
}

type BulkCreatePlanRequest struct {
	Plans []*CreatePlanRequest `json:"plans"`
}

type BulkDeletePlanRequest struct {
	Plans []uuid.UUID `json:"plans"`
}
