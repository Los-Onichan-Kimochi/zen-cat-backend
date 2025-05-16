package schemas

import "github.com/google/uuid"

type CommunityPlan struct {
	CommunityId uuid.UUID `json:"community_id"`
	PlanId      uuid.UUID `json:"plan_id"`
}

type CreateCommunityPlanRequest struct {
	CommunityId uuid.UUID `json:"community_id" validate:"required"`
	PlanId      uuid.UUID `json:"plan_id"      validate:"required"`
}
