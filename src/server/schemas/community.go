package schemas

import "github.com/google/uuid"

type Community struct {
	Id                  uuid.UUID `json:"id"`
	Name                string    `json:"name"`
	Purpose             string    `json:"purpose"`
	ImageUrl            string    `json:"image_url"`
	NumberSubscriptions int       `json:"number_subscriptions"`
}

type Communities struct {
	Communities []*Community `json:"communities"`
	// TODO: Add pagination if needed
}

type CreateCommunityRequest struct {
	Name       string  `json:"name"`
	Purpose    string  `json:"purpose"`
	ImageUrl   string  `json:"image_url"`
	ImageBytes *[]byte `json:"image_bytes"`
}

type UpdateCommunityRequest struct {
	Name       *string `json:"name"`
	Purpose    *string `json:"purpose"`
	ImageUrl   *string `json:"image_url"`
	ImageBytes *[]byte `json:"image_bytes"`
}

type BatchCreateCommunityRequest struct {
	Communities []*CreateCommunityRequest `json:"communities"`
}

type BulkDeleteCommunityRequest struct {
	Communities []uuid.UUID `json:"communities"`
}

type CommunityWithImage struct {
	Community
	ImageBytes *[]byte `json:"image_bytes"`
}
