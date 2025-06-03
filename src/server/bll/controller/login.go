package controller

import (
	"time"

	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/bll/adapter"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type Login struct {
	Logger      logging.Logger
	Adapter     *adapter.AdapterCollection
	EnvSettings *schemas.EnvSettings
	Auth        *Auth // Referencia al controlador de Auth para reutilizar funciones
}

func NewLoginController(
	logger logging.Logger,
	adapter *adapter.AdapterCollection,
	envSettings *schemas.EnvSettings,
	auth *Auth,
) *Login {
	return &Login{
		Logger:      logger,
		Adapter:     adapter,
		EnvSettings: envSettings,
		Auth:        auth,
	}
}

// Authenticates a user with email and password and returns extended user info
func (l *Login) Login(
	email string,
	password string,
) (*schemas.LoginResponse, *errors.Error) {
	// Get user by email
	user, err := l.Adapter.User.GetPostgresqlUserByEmail(email)
	if err != nil {
		return nil, &errors.AuthenticationError.UnauthorizedUser
	}

	// Verify password (in production, you should use bcrypt or similar)
	if user.Password != password {
		return nil, &errors.AuthenticationError.UnauthorizedUser
	}

	// Generate user roles array
	userRoles := []string{string(user.Rol)}

	// Generate token with extended user info
	tokenResponse := l.Auth.GenerateToken(
		user.Id,
		user.Email,
		user.Password,
		userRoles,
		time.Hour*2, // 2 hours expiration
	)

	// Return login response with user info and tokens
	return &schemas.LoginResponse{
		User: schemas.UserProfile{
			Id:             user.Id,
			Name:           user.Name,
			FirstLastName:  user.FirstLastName,
			SecondLastName: user.SecondLastName,
			Email:          user.Email,
			Rol:            user.Rol,
			ImageUrl:       user.ImageUrl,
		},
		Tokens: tokenResponse,
	}, nil
}

// Registers a new user and returns extended user info with tokens
func (l *Login) Register(
	name string,
	firstLastName string,
	secondLastName *string,
	email string,
	password string,
	imageUrl string,
) (*schemas.LoginResponse, *errors.Error) {
	// Check if user already exists
	existingUser, _ := l.Adapter.User.GetPostgresqlUserByEmail(email)
	if existingUser != nil {
		return nil, &errors.ConflictError.UserAlreadyExists
	}

	// Create new user with default role CLIENT
	user, err := l.Adapter.User.CreatePostgresqlUser(
		name,
		firstLastName,
		secondLastName,
		password, // In production, this should be hashed with bcrypt
		email,
		string(schemas.UserRolClient), // Default to CLIENT role
		imageUrl,
		"SYSTEM", // Updated by system for self-registration
	)
	if err != nil {
		return nil, err
	}

	// Generate user roles array
	userRoles := []string{string(user.Rol)}

	// Generate token for the new user
	tokenResponse := l.Auth.GenerateToken(
		user.Id,
		user.Email,
		user.Password,
		userRoles,
		time.Hour*2, // 2 hours expiration
	)

	// Return login response with user info and tokens
	return &schemas.LoginResponse{
		User: schemas.UserProfile{
			Id:             user.Id,
			Name:           user.Name,
			FirstLastName:  user.FirstLastName,
			SecondLastName: user.SecondLastName,
			Email:          user.Email,
			Rol:            user.Rol,
			ImageUrl:       user.ImageUrl,
		},
		Tokens: tokenResponse,
	}, nil
}
