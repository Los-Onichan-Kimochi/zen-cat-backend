package schemas

import "github.com/google/uuid"

type CommunityPlan struct {
	Id          uuid.UUID `json:"id"`
	CommunityId uuid.UUID `json:"community_id"`
	PlanId      uuid.UUID `json:"plan_id"`
}

type CommunityPlans struct {
	CommunityPlans []*CommunityPlan `json:"community_plans"`
}

type CreateCommunityPlanRequest struct {
	CommunityId uuid.UUID `json:"community_id" validate:"required"`
	PlanId      uuid.UUID `json:"plan_id"      validate:"required"`
}

type BatchCreateCommunityPlanRequest struct {
	CommunityPlans []*CreateCommunityPlanRequest `json:"community_plans"`
}

type BulkDeleteCommunityPlanRequest struct {
	CommunityPlans []*CommunityPlan `json:"community_plans"`
}
