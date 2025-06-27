package login_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
	"onichankimochi.com/astro_cat_backend/src/server/utils"
)

func TestLoginSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid user exists with correct credentials
		WHEN:  Login is called with valid email and password
		THEN:  User should be authenticated and tokens returned
	*/
	// GIVEN
	loginController, _, db := controllerTest.NewLoginControllerTestWrapper(t)

	// Create test user with known password
	plainPassword := "testPassword123"
	hashedPassword, hashErr := utils.HashPassword(plainPassword)
	assert.NoError(t, hashErr)

	secondLastName := "Smith"
	testUser := &model.User{
		Name:           "John",
		FirstLastName:  "Doe",
		SecondLastName: &secondLastName,
		Email:          "john.doe@example.com",
		Password:       hashedPassword,
		Rol:            model.UserRolClient,
		ImageUrl:       "https://example.com/profile.jpg",
		AuditFields: model.AuditFields{
			UpdatedBy: "SYSTEM",
		},
	}
	err := db.Create(testUser).Error
	assert.NoError(t, err)

	// WHEN
	result, loginErr := loginController.Login(testUser.Email, plainPassword)

	// THEN
	assert.Nil(t, loginErr)
	assert.NotNil(t, result)
	assert.NotNil(t, result.User)
	assert.NotNil(t, result.Tokens)

	// Verify user profile
	assert.Equal(t, testUser.Id, result.User.Id)
	assert.Equal(t, testUser.Name, result.User.Name)
	assert.Equal(t, testUser.FirstLastName, result.User.FirstLastName)
	assert.Equal(t, testUser.SecondLastName, result.User.SecondLastName)
	assert.Equal(t, testUser.Email, result.User.Email)
	assert.Equal(t, string(testUser.Rol), string(result.User.Rol))
	assert.Equal(t, testUser.ImageUrl, result.User.ImageUrl)

	// Verify tokens are present
	assert.NotEmpty(t, result.Tokens.AccessToken)
	assert.NotEmpty(t, result.Tokens.RefreshToken)
}

func TestLoginWithInvalidPassword(t *testing.T) {
	/*
		GIVEN: Valid user exists
		WHEN:  Login is called with incorrect password
		THEN:  It should return unauthorized user error
	*/
	// GIVEN
	loginController, _, db := controllerTest.NewLoginControllerTestWrapper(t)

	// Create test user
	testUser := factories.NewUserModel(db)
	wrongPassword := "wrongPassword"

	// WHEN
	result, err := loginController.Login(testUser.Email, wrongPassword)

	// THEN
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, errors.AuthenticationError.UnauthorizedUser.Code, err.Code)
}

func TestLoginWithNonExistentUser(t *testing.T) {
	/*
		GIVEN: User does not exist
		WHEN:  Login is called with non-existent email
		THEN:  It should return unauthorized user error
	*/
	// GIVEN
	loginController, _, _ := controllerTest.NewLoginControllerTestWrapper(t)

	nonExistentEmail := "nonexistent@example.com"
	password := "anyPassword"

	// WHEN
	result, err := loginController.Login(nonExistentEmail, password)

	// THEN
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, errors.AuthenticationError.UnauthorizedUser.Code, err.Code)
}

func TestLoginWithEmptyCredentials(t *testing.T) {
	/*
		GIVEN: Empty email and password
		WHEN:  Login is called with empty credentials
		THEN:  It should return unauthorized user error
	*/
	// GIVEN
	loginController, _, _ := controllerTest.NewLoginControllerTestWrapper(t)

	emptyEmail := ""
	emptyPassword := ""

	// WHEN
	result, err := loginController.Login(emptyEmail, emptyPassword)

	// THEN
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, errors.AuthenticationError.UnauthorizedUser.Code, err.Code)
}

func TestLoginWithDifferentUserRoles(t *testing.T) {
	/*
		GIVEN: Users with different roles exist
		WHEN:  Users are created with different roles
		THEN:  Each should have the correct role assigned
	*/
	// GIVEN
	_, _, db := controllerTest.NewLoginControllerTestWrapper(t)

	// Create role variables for pointer assignment
	adminRole := model.UserRolAdmin
	clientRole := model.UserRolClient

	// Create admin user using factory
	adminUser := factories.NewUserModel(db, factories.UserModelF{
		Rol: &adminRole,
	})

	// Create client user using factory
	clientUser := factories.NewUserModel(db, factories.UserModelF{
		Rol: &clientRole,
	})

	// THEN - Verify roles are correctly assigned
	assert.Equal(t, model.UserRolAdmin, adminUser.Rol)
	assert.Equal(t, model.UserRolClient, clientUser.Rol)

	// Verify in database
	var dbAdminUser model.User
	err := db.Where("id = ?", adminUser.Id).First(&dbAdminUser).Error
	assert.NoError(t, err)
	assert.Equal(t, model.UserRolAdmin, dbAdminUser.Rol)

	var dbClientUser model.User
	err = db.Where("id = ?", clientUser.Id).First(&dbClientUser).Error
	assert.NoError(t, err)
	assert.Equal(t, model.UserRolClient, dbClientUser.Rol)
}
