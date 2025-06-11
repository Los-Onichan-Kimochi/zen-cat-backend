package schemas

import "github.com/google/uuid"

type UserRol string

const (
	UserRolAdmin  UserRol = "ADMINISTRATOR"
	UserRolClient UserRol = "CLIENT"
)

type User struct {
	Id             uuid.UUID     `json:"id"`
	Name           string        `json:"name"`
	FirstLastName  string        `json:"first_last_name"`
	SecondLastName *string       `json:"second_last_name"`
	Password       string        `json:"password"`
	Email          string        `json:"email"`
	Rol            UserRol       `json:"rol"`
	ImageUrl       string        `json:"image_url"`
	Memberships    []*Membership `json:"memberships,omitempty"`
	Onboarding     *Onboarding   `json:"onboarding,omitempty"`
}

type Users struct {
	Users []*User `json:"users"`
}

type CreateUserRequest struct {
	Name           string        `json:"name"`
	FirstLastName  string        `json:"first_last_name"`
	SecondLastName string        `json:"second_last_name"`
	Password       string        `json:"password"`
	Email          string        `json:"email"`
	Rol            string        `json:"rol"`
	ImageUrl       string        `json:"image_url"`
	Onboarding     *Onboarding   `json:"onboarding,omitempty"`
	Memberships    []*Membership `json:"memberships,omitempty"`
}

type UpdateUserRequest struct {
	Name           *string       `json:"name"`
	FirstLastName  *string       `json:"first_last_name"`
	SecondLastName *string       `json:"second_last_name"`
	Password       *string       `json:"password"`
	Email          *string       `json:"email"`
	Rol            *string       `json:"rol"`
	ImageUrl       *string       `json:"image_url"`
	Onboarding     *Onboarding   `json:"onboarding,omitempty"`
	Memberships    []*Membership `json:"memberships,omitempty"`
}

type BulkDeleteUserRequest struct {
	Users []uuid.UUID `json:"users"`
}

type BulkCreateUserRequest struct {
	Users []*CreateUserRequest `json:"users"`
}

type ChangePasswordInput struct {
	NewPassword string `json:"new_password" validate:"required"`
}
