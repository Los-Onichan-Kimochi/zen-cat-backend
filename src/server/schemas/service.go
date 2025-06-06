package schemas

import "github.com/google/uuid"

type Service struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ImageUrl    string    `json:"image_url"`
	IsVirtual   bool      `json:"is_virtual"`
}

type Services struct {
	Services []*Service `json:"services"`
	// TODO: Add pagination if needed
}

type CreateServiceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url"`
	IsVirtual   bool   `json:"is_virtual"`
}

type UpdateServiceRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	ImageUrl    *string `json:"image_url"`
	IsVirtual   *bool   `json:"is_virtual"`
}

type BulkDeleteServiceRequest struct {
	Services []string `json:"services"`
}
