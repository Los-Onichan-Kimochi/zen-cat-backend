package schemas

import "github.com/google/uuid"

type UserRol string

const (
	UserRolAdmin  UserRol = "ADMINISTRATOR"
	UserRolClient UserRol = "CLIENT"
	UserRolGuest  UserRol = "GUEST"
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
	ImageBytes     *[]byte       `json:"image_bytes"`
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
	ImageBytes     *[]byte       `json:"image_bytes"`
	Onboarding     *Onboarding   `json:"onboarding,omitempty"`
	Memberships    []*Membership `json:"memberships,omitempty"`
}

type BulkDeleteUserRequest struct {
	Users []uuid.UUID `json:"users"`
}

type BulkCreateUserRequest struct {
	Users []*CreateUserRequest `json:"users"`
}

type UserWithImage struct {
	User
	ImageBytes *[]byte `json:"image_bytes"`
}

type ChangePasswordInput struct {
	Email       string `json:"email"        validate:"required,email"`
	NewPassword string `json:"new_password" validate:"required"`
}

type CheckUserExistsResponse struct {
	Email  string `json:"email"`
	Exists bool   `json:"exists"`
}

// Schema for role management
type ChangeUserRoleRequest struct {
	Rol UserRol `json:"rol"`
}

// Schema for user statistics
type UserStats struct {
	TotalUsers        int64                  `json:"total_users"`
	AdminCount        int64                  `json:"admin_count"`
	ClientCount       int64                  `json:"client_count"`
	GuestCount        int64                  `json:"guest_count"`
	RoleDistribution  []UserRoleDistribution `json:"role_distribution"`
	RecentConnections []UserConnection       `json:"recent_connections"`
}

type UserRoleDistribution struct {
	Role  UserRol `json:"role"`
	Count int64   `json:"count"`
}

type UserConnection struct {
	UserId       uuid.UUID `json:"user_id"`
	UserEmail    string    `json:"user_email"`
	UserName     string    `json:"user_name"`
	Role         UserRol   `json:"role"`
	LastLogin    *string   `json:"last_login"`
	ConnectionIP *string   `json:"connection_ip"`
}
