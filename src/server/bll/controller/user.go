package controller

import (
	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
	bllAdapter "onichankimochi.com/astro_cat_backend/src/server/bll/adapter"
	errors "onichankimochi.com/astro_cat_backend/src/server/errors"
	schemas "onichankimochi.com/astro_cat_backend/src/server/schemas"
	"onichankimochi.com/astro_cat_backend/src/server/utils"
)

type User struct {
	logger      logging.Logger
	Adapter     *bllAdapter.AdapterCollection
	EnvSettings *schemas.EnvSettings
}

func NewUserController(
	logger logging.Logger,
	adapter *bllAdapter.AdapterCollection,
	envSettings *schemas.EnvSettings,
) *User {
	return &User{
		logger:      logger,
		Adapter:     adapter,
		EnvSettings: envSettings,
	}
}

func (u *User) GetUser(userId uuid.UUID) (*schemas.User, *errors.Error) {
	return u.Adapter.User.GetPostgresqlUser(userId)
}

func (u *User) FetchUsers() (*schemas.Users, *errors.Error) {
	users, err := u.Adapter.User.FetchPostgresqlUsers()
	if err != nil {
		return nil, err
	}

	return &schemas.Users{Users: users}, nil
}

func (u *User) CreateUser(
	createUserRequest schemas.CreateUserRequest,
	updatedBy string,
) (*schemas.User, *errors.Error) {
	var secondLastName *string
	if createUserRequest.SecondLastName != "" {
		secondLastName = &createUserRequest.SecondLastName
	}

	return u.Adapter.User.CreatePostgresqlUser(
		createUserRequest.Name,
		createUserRequest.FirstLastName,
		secondLastName,
		createUserRequest.Password,
		createUserRequest.Email,
		createUserRequest.Rol,
		createUserRequest.ImageUrl,
		updatedBy,
		createUserRequest.Memberships,
		createUserRequest.Onboarding,
	)
}

func (u *User) UpdateUser(
	userId uuid.UUID,
	updateUserRequest schemas.UpdateUserRequest,
	updatedBy string,
) (*schemas.User, *errors.Error) {
	return u.Adapter.User.UpdatePostgresqlUser(
		userId,
		updateUserRequest.Name,
		updateUserRequest.FirstLastName,
		updateUserRequest.SecondLastName,
		updateUserRequest.Password,
		updateUserRequest.Email,
		updateUserRequest.Rol,
		updateUserRequest.ImageUrl,
		updateUserRequest.Memberships,
		updateUserRequest.Onboarding,
		updatedBy,
	)
}

func (u *User) DeleteUser(userId uuid.UUID) *errors.Error {
	return u.Adapter.User.DeletePostgresqlUser(userId)
}

func (u *User) BulkCreateUsers(
	createUsersData []*schemas.CreateUserRequest,
	updatedBy string,
) ([]*schemas.User, *errors.Error) {
	return u.Adapter.User.BulkCreatePostgresqlUser(
		createUsersData,
		updatedBy,
	)
}

func (u *User) BulkDeleteUsers(
	bulkDeleteUsersData schemas.BulkDeleteUserRequest,
) *errors.Error {
	return u.Adapter.User.BulkDeletePostgresqlUser(
		bulkDeleteUsersData.Users,
	)
}

func (u *User) CheckUserExistsByEmail(email string) (*schemas.CheckUserExistsResponse, *errors.Error) {
	_, err := u.Adapter.User.GetPostgresqlUserByEmail(email)

	// Si no hay error, el usuario existe
	exists := err == nil

	return &schemas.CheckUserExistsResponse{
		Email:  email,
		Exists: exists,
	}, nil
}

func (u *User) ChangePassword(
	email string,
	request schemas.ChangePasswordInput,
) *errors.Error {
	// Buscar usuario por email
	user, err := u.Adapter.User.GetPostgresqlUserByEmail(email)
	if err != nil {
		return &errors.ObjectNotFoundError.UserNotFound
	}

	// Hashear la nueva contraseña
	hashedPassword, hashErr := utils.HashPassword(request.NewPassword)
	if hashErr != nil {
		return &errors.InternalServerError.Default
	}

	// Actualizar contraseña usando el ID del usuario encontrado
	updateErr := u.Adapter.User.UpdateUserPassword(user.Id, hashedPassword)
	if updateErr != nil {
		return &errors.BadRequestError.UserPasswordNotUpdated
	}

	return nil
}
