package schemas

import "github.com/google/uuid"

type CommunityService struct {
	CommunityId uuid.UUID `json:"community_id"`
	ServiceId   uuid.UUID `json:"service_id"`
}

type CreateCommunityServiceRequest struct {
	CommunityId uuid.UUID `json:"community_id" validate:"required"`
	ServiceId   uuid.UUID `json:"service_id"   validate:"required"`
}
