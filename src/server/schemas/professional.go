package schemas

import "github.com/google/uuid"

type Professional struct {
	Id             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	FirstLastName  string    `json:"first_last_name"`
	SecondLastName string    `json:"second_last_name"`
	Specialty      string    `json:"specialty"`
	Email          string    `json:"email"`
	PhoneNumber    string    `json:"phone_number"`
	Type           string    `json:"type"`
	ImageUrl       string    `json:"image_url"`
}

type Professionals struct {
	Professionals []*Professional `json:"professionals"`
}

type CreateProfessionalRequest struct {
	Name           string `json:"name"`
	FirstLastName  string `json:"first_last_name"`
	SecondLastName string `json:"second_last_name"`
	Specialty      string `json:"specialty"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	Type           string `json:"type"`
	ImageUrl       string `json:"image_url"`
}

type UpdateProfessionalRequest struct {
	Name           *string `json:"name"`
	FirstLastName  *string `json:"first_last_name"`
	SecondLastName *string `json:"second_last_name"`
	Specialty      *string `json:"specialty"`
	Email          *string `json:"email"`
	PhoneNumber    *string `json:"phone_number"`
	Type           *string `json:"type"`
	ImageUrl       *string `json:"image_url"`
}
