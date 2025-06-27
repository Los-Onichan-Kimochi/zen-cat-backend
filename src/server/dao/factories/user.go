package factories

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/utils"
)

type UserModelF struct {
	Id             *uuid.UUID
	Name           *string
	FirstLastName  *string
	SecondLastName *string
	Password       *string
	Email          *string
	Rol            *model.UserRol
	ImageUrl       *string
}

// Create a new user on DB
func NewUserModel(db *gorm.DB, option ...UserModelF) *model.User {
	// Default password
	hashedPassword, err := utils.HashPassword("testpassword123")
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
	}

	user := &model.User{
		Id:             uuid.New(),
		Name:           "TestUser",
		FirstLastName:  "TestLastName",
		SecondLastName: nil,
		Password:       hashedPassword,
		Email:          "test@example.com",
		Rol:            model.UserRolClient,
		ImageUrl:       "https://example.com/avatar.jpg",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}

	if len(option) > 0 {
		parameters := option[0]
		if parameters.Id != nil {
			user.Id = *parameters.Id
		}
		if parameters.Name != nil {
			user.Name = *parameters.Name
		}
		if parameters.FirstLastName != nil {
			user.FirstLastName = *parameters.FirstLastName
		}
		if parameters.SecondLastName != nil {
			user.SecondLastName = parameters.SecondLastName
		}
		if parameters.Password != nil {
			hashedPassword, err := utils.HashPassword(*parameters.Password)
			if err != nil {
				log.Fatalf("Error hashing password: %v", err)
			}
			user.Password = hashedPassword
		}
		if parameters.Email != nil {
			user.Email = *parameters.Email
		}
		if parameters.Rol != nil {
			user.Rol = *parameters.Rol
		}
		if parameters.ImageUrl != nil {
			user.ImageUrl = *parameters.ImageUrl
		}
	}

	result := db.Create(user)
	if result.Error != nil {
		log.Fatalf("Error when trying to create user: %v", result.Error)
	}

	return user
}

// Create size number of new users on DB
func NewUserModelBatch(
	db *gorm.DB,
	size int,
	option ...UserModelF,
) []*model.User {
	users := []*model.User{}
	for i := 0; i < size; i++ {
		var user *model.User
		if len(option) > 0 {
			user = NewUserModel(db, option[0])
		} else {
			user = NewUserModel(db)
		}
		users = append(users, user)
	}
	return users
}
