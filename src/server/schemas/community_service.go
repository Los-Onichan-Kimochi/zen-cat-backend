package schemas

import "github.com/google/uuid"

type CommunityService struct {
	CommunityId uuid.UUID `json:"community_id"`
	ServiceId   uuid.UUID `json:"service_id"`
}

type CommunityServices struct {
	CommunityServices []*CommunityService `json:"community_services"`
}

type CreateCommunityServiceRequest struct {
	CommunityId uuid.UUID `json:"community_id" validate:"required"`
	ServiceId   uuid.UUID `json:"service_id"   validate:"required"`
}

type BatchCreateCommunityServiceRequest struct {
	CommunityServices []*CreateCommunityServiceRequest `json:"community_services"`
}
