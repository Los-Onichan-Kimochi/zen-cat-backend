package controller

import (
	"time"

	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/bll/adapter"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	"onichankimochi.com/astro_cat_backend/src/server/utils"
)

type Login struct {
	Logger      logging.Logger
	Adapter     *adapter.AdapterCollection
	EnvSettings *schemas.EnvSettings
	Auth        *Auth
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

func (l *Login) Login(
	email string,
	password string,
) (*schemas.LoginResponse, *errors.Error) {
	user, err := l.Adapter.User.GetPostgresqlUserByEmail(email)
	if err != nil {
		return nil, &errors.AuthenticationError.UnauthorizedUser
	}

	if err := utils.CheckPasswordHash(password, user.Password); err != nil {
		return nil, &errors.AuthenticationError.UnauthorizedUser
	}

	userRoles := []string{string(user.Rol)}

	tokenResponse, tokenErr := l.Auth.GenerateToken(
		user.Id,
		user.Email,
		user.Password,
		userRoles,
		time.Hour*2,
	)
	if tokenErr != nil {
		return nil, tokenErr
	}

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

func (l *Login) Register(
	name string,
	firstLastName string,
	secondLastName *string,
	email string,
	password string,
	imageUrl string,
) (*schemas.LoginResponse, *errors.Error) {
	existingUser, _ := l.Adapter.User.GetPostgresqlUserByEmail(email)
	if existingUser != nil {
		return nil, &errors.ConflictError.UserAlreadyExists
	}

	user, err := l.Adapter.User.CreatePostgresqlUser(
		name,
		firstLastName,
		secondLastName,
		password,
		email,
		string(schemas.UserRolClient),
		imageUrl,
		"SYSTEM",
	)
	if err != nil {
		return nil, err
	}

	userRoles := []string{string(user.Rol)}

	tokenResponse, tokenErr := l.Auth.GenerateToken(
		user.Id,
		user.Email,
		user.Password,
		userRoles,
		time.Hour*2,
	)
	if tokenErr != nil {
		return nil, tokenErr
	}

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
