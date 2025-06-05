package controller

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/utils"
)

type User struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

func NewUserController(logger logging.Logger, postgresqlDB *gorm.DB) *User {
	return &User{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

func (u *User) GetUser(userId uuid.UUID) (*model.User, error) {
	user := &model.User{}
	result := u.PostgresqlDB.
		Preload("Memberships").
		Preload("Memberships.Community").
		Preload("Memberships.Plan").
		First(&user, "id = ?", userId)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (u *User) GetUserByEmail(email string) (*model.User, error) {
	user := &model.User{}
	result := u.PostgresqlDB.
		Preload("Memberships").
		Preload("Memberships.Community").
		Preload("Memberships.Plan").
		First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (u *User) FetchUsers() ([]*model.User, error) {
	users := []*model.User{}
	result := u.PostgresqlDB.
		Preload("Memberships").
		Preload("Memberships.Community").
		Preload("Memberships.Plan").
		Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (u *User) CreateUser(user *model.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return u.PostgresqlDB.Create(user).Error
}

func (u *User) UpdateUser(
	id uuid.UUID,
	name *string,
	firstLastName *string,
	secondLastName *string,
	password *string,
	email *string,
	rol *string,
	imageUrl *string,
	updatedBy string,
) (*model.User, error) {
	updateFields := map[string]any{
		"updated_by": updatedBy,
	}

	if name != nil {
		updateFields["name"] = *name
	}
	if firstLastName != nil {
		updateFields["first_last_name"] = *firstLastName
	}
	if secondLastName != nil {
		updateFields["second_last_name"] = *secondLastName
	}
	if password != nil {
		updateFields["password"] = *password
	}
	if email != nil {
		updateFields["email"] = *email
	}
	if rol != nil {
		updateFields["rol"] = *rol
	}
	if imageUrl != nil {
		updateFields["image_url"] = *imageUrl
	}

	var user model.User
	if len(updateFields) == 1 {
		if err := u.PostgresqlDB.First(&user, "id = ?", id).Error; err != nil {
			return nil, err
		}
		return &user, nil
	}

	result := u.PostgresqlDB.Model(&user).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(updateFields)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (u *User) DeleteUser(userId uuid.UUID) error {
	result := u.PostgresqlDB.Model(&model.User{}).
		Where("id = ?", userId).
		Update("deleted_at", gorm.Expr("NOW()"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
