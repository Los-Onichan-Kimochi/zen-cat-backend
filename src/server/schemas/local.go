package schemas

import "github.com/google/uuid"

type Local struct {
	Id             uuid.UUID `json:"id"`
	LocalName      string    `json:"local_name"`
	StreetName     string    `json:"street_name"`
	BuildingNumber string    `json:"building_number"`
	District       string    `json:"district"`
	Province       string    `json:"province"`
	Region         string    `json:"region"`
	Reference      string    `json:"reference"`
	Capacity       int       `json:"capacity"`
	ImageUrl       string    `json:"image_url"`
}

type Locals struct {
	Locals []*Local `json:"locals"`
}

type CreateLocalRequest struct {
	LocalName      string `json:"local_name"`
	StreetName     string `json:"street_name"`
	BuildingNumber string `json:"building_number"`
	District       string `json:"district"`
	Province       string `json:"province"`
	Region         string `json:"region"`
	Reference      string `json:"reference"`
	Capacity       int    `json:"capacity"`
	ImageUrl       string `json:"image_url"`
}

type UpdateLocalRequest struct {
	LocalName      *string `json:"local_name"`
	StreetName     *string `json:"street_name"`
	BuildingNumber *string `json:"building_number"`
	District       *string `json:"district"`
	Province       *string `json:"province"`
	Region         *string `json:"region"`
	Reference      *string `json:"reference"`
	Capacity       *int    `json:"capacity"`
	ImageUrl       *string `json:"image_url"`
}

type BatchCreateLocalRequest struct {
	Locals []*CreateLocalRequest `json:"locals"`
}

type BulkDeleteLocalRequest struct {
	Locals []string `json:"locals"`
}
