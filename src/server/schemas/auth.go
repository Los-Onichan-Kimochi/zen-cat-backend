package schemas

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Name           string  `json:"name" validate:"required"`
	FirstLastName  string  `json:"first_last_name" validate:"required"`
	SecondLastName *string `json:"second_last_name"`
	Email          string  `json:"email" validate:"required,email"`
	Password       string  `json:"password" validate:"required,min=6"`
	ImageUrl       string  `json:"image_url"`
}

type LoginResponse struct {
	User   UserProfile   `json:"user"`
	Tokens TokenResponse `json:"tokens"`
}

type UserProfile struct {
	Id             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	FirstLastName  string    `json:"first_last_name"`
	SecondLastName *string   `json:"second_last_name"`
	Email          string    `json:"email"`
	Rol            UserRol   `json:"rol"`
	ImageUrl       string    `json:"image_url"`
}

type TokenResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    time.Duration `json:"expires_in"`
}

type CustomClaims struct {
	UserId        uuid.UUID `json:"user_id"`
	UserEmail     string    `json:"user_email"`
	UserPassword  string    `json:"user_password"`
	UserRoles     []string  `json:"user_roles"`
	UserName      string    `json:"user_name"`
	UserFirstName string    `json:"user_first_name"`
	UserLastName  *string   `json:"user_last_name"`
	UserImageUrl  string    `json:"user_image_url"`
	jwt.RegisteredClaims
}

type Credentials struct {
	UserId        uuid.UUID `json:"user_id"`
	UserEmail     string    `json:"user_email"`
	UserPassword  string    `json:"user_password"`
	UserRoles     []string  `json:"user_roles"`
	UserName      string    `json:"user_name"`
	UserFirstName string    `json:"user_first_name"`
	UserLastName  *string   `json:"user_last_name"`
	UserImageUrl  string    `json:"user_image_url"`
}

type GoogleLoginRequest struct {
	Token string `json:"token" example:"eyJhbGciOi..."`
}

type GoogleLoginResponse struct {
	User   User          `json:"user"`
	Tokens TokenResponse `json:"tokens"`
}
