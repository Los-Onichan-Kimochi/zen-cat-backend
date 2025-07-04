package controller

import (
	"context"
	"time"

	"google.golang.org/api/idtoken"

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
		nil,
		nil,
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

func (l *Login) GoogleLogin(
	ctx context.Context,
	idToken string,
) (*schemas.GoogleLoginResponse, *errors.Error) {
	// 1. Validar token de Google
	payload, err := idtoken.Validate(ctx, idToken, "")
	if err != nil {
		l.Logger.Error("Token de Google inválido:", err)
		return nil, &errors.AuthenticationError.InvalidAccessToken
	}

	// 2. Extraer datos del token
	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)

	// 3. Buscar usuario por email
	user, userErr := l.Adapter.User.GetPostgresqlUserByEmail(email)
	if userErr != nil {
		// 4. Crear usuario si no existe
		newUser, createErr := l.Adapter.User.CreatePostgresqlUser(
			name,
			"",  // firstLastName
			nil, // secondLastName
			"",  // password
			email,
			string(schemas.UserRolClient),
			picture,
			"SYSTEM",
			nil,
			nil,
		)
		if createErr != nil {
			l.Logger.Error("Error al crear usuario con Google:", createErr)
			return nil, &errors.InternalServerError.Default
		}

		// reconstruir user localmente (porque Create no lo retorna)
		user = newUser
	}

	userRoles := []string{string(user.Rol)}

	// 5. Generar tokens JWT
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

	// 6. Retornar estructura de respuesta
	return &schemas.GoogleLoginResponse{
		User:   *user,
		Tokens: tokenResponse,
	}, nil
}
