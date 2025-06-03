package controller

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/bll/adapter"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
)

type Auth struct {
	Logger logging.Logger
	Adapter *adapter.AdapterCollection
	EnvSettings *schemas.EnvSettings
}

func NewAuthController(
	logger logging.Logger,
	adapter *adapter.AdapterCollection,
	envSettings *schemas.EnvSettings,
) *Auth {
	return &Auth{
		Logger: logger,
		Adapter: adapter,
		EnvSettings: envSettings,
	}
}

// Generates a new token for a player
func (a *Auth) GenerateToken(
	userId uuid.UUID,
	userEmail string,
	userPassword string,
	userRoles []string,
	expirationDelta time.Duration,
) schemas.TokenResponse {
	// Set custom claim
	expirationTime := time.Now().Add(expirationDelta)

	claims := &schemas.CustomClaims{
		UserId: userId,
		UserEmail: userEmail,
		UserPassword: userPassword,
		UserRoles: userRoles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, _ := token.SignedString(a.EnvSettings.TokenSignatureKey)


	return schemas.TokenResponse{
		AccessToken: accessToken,
		RefreshToken: accessToken,
		ExpiresIn: expirationDelta,
	}
}

func (a *Auth) AccessTokenValidation(
	c echo.Context,
) (*jwt.Token, *schemas.Credentials, *errors.Error) {

	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return nil, nil, &errors.AuthenticationError.UnauthorizedUser
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return nil, nil, &errors.AuthenticationError.UnauthorizedUser
	}

	accessToken, tokenErr := jwt.ParseWithClaims(tokenString, &schemas.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(a.EnvSettings.TokenSignatureKey), nil
	})

	if tokenErr != nil || !accessToken.Valid {
		return nil, nil, &errors.AuthenticationError.UnauthorizedUser
	}

	credentials, signErr := a.GetCredentialByAccessToken(accessToken)
	if signErr != nil {
		return nil, nil, &errors.AuthenticationError.UnauthorizedUser
	}

	return accessToken, credentials, nil
}

func (a *Auth) GetCredentialByAccessToken(
	accessToken *jwt.Token,
) (*schemas.Credentials, *errors.Error) {
	claims := accessToken.Claims.(*schemas.CustomClaims)
	userId := claims.UserId
	userEmail := claims.UserEmail
	userPassword := claims.UserPassword
	userRoles := claims.UserRoles

	user, err := a.Adapter.User.GetPostgresqlUser(userId)
	if err != nil {
		return nil, err
	}

	if user.Email != userEmail {
		return nil, &errors.AuthenticationError.UnauthorizedUser
	}

	if user.Password != userPassword {
		return nil, &errors.AuthenticationError.UnauthorizedUser
	}

	return &schemas.Credentials{
		UserId: userId,
		UserEmail: userEmail,
		UserPassword: userPassword,
		UserRoles: userRoles,
	}, nil
}

func (a *Auth) RefreshToken(accessToken *jwt.Token) (*schemas.TokenResponse, *errors.Error) {
	credentials, err := a.GetCredentialByAccessToken(accessToken)
	if err != nil {
		return nil, &errors.AuthenticationError.UnauthorizedUser
	}
	newToken := a.GenerateToken(
		credentials.UserId,
		credentials.UserEmail,
		credentials.UserPassword,
		credentials.UserRoles,
		time.Hour * 2,
	)

	return &newToken, nil
}
