package controller

import (
	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
	bllAdapter "onichankimochi.com/astro_cat_backend/src/server/bll/adapter"
	errors "onichankimochi.com/astro_cat_backend/src/server/errors"
	schemas "onichankimochi.com/astro_cat_backend/src/server/schemas"
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

		updatedBy,
	)
}

func (u *User) DeleteUser(userId uuid.UUID) *errors.Error {
	return u.Adapter.User.DeletePostgresqlUser(userId)
}

// Bulk deletes users.
func (u *User) BulkDeleteUsers(
	bulkDeleteUserData schemas.BulkDeleteUserRequest,
) *errors.Error {
	return u.Adapter.User.BulkDeletePostgresqlUsers(
		bulkDeleteUserData.Users,
	)
}

// todo : bulk-create-users
