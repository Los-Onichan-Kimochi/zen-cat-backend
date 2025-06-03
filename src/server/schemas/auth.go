package schemas

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenResponse struct {
	AccessToken string
	RefreshToken string
	ExpiresIn time.Duration
}

type CustomClaims struct {
	UserId uuid.UUID
	UserEmail string
	UserPassword string
	UserRoles []string
	jwt.RegisteredClaims
}

type Credentials struct {
	UserId uuid.UUID
	UserEmail string
	UserPassword string
	UserRoles []string
}
