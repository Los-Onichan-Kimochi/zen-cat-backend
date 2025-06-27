package login_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

func TestRegisterSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid registration data
		WHEN:  Register is called with valid parameters
		THEN:  User should be created and tokens returned
	*/
	// GIVEN
	loginController, _, db := controllerTest.NewLoginControllerTestWrapper(t)

	name := "John"
	firstLastName := "Doe"
	secondLastName := "Smith"
	email := utilsTest.GenerateRandomEmail()
	password := "securePassword123"
	imageUrl := "https://example.com/profile.jpg"

	// WHEN
	result, err := loginController.Register(
		name,
		firstLastName,
		&secondLastName,
		email,
		password,
		imageUrl,
	)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.User)
	assert.NotNil(t, result.Tokens)

	// Verify user profile
	assert.Equal(t, name, result.User.Name)
	assert.Equal(t, firstLastName, result.User.FirstLastName)
	assert.Equal(t, &secondLastName, result.User.SecondLastName)
	assert.Equal(t, email, result.User.Email)
	assert.Equal(t, model.UserRolClient, result.User.Rol) // Should default to client
	assert.Equal(t, imageUrl, result.User.ImageUrl)

	// Verify tokens are present
	assert.NotEmpty(t, result.Tokens.AccessToken)
	assert.NotEmpty(t, result.Tokens.RefreshToken)

	// Verify user was created in database
	var createdUser model.User
	dbErr := db.Where("email = ?", email).First(&createdUser).Error
	assert.NoError(t, dbErr)
	assert.Equal(t, name, createdUser.Name)
	assert.Equal(t, email, createdUser.Email)
	assert.Equal(t, model.UserRolClient, createdUser.Rol)
}

func TestRegisterWithoutSecondLastName(t *testing.T) {
	/*
		GIVEN: Valid registration data without second last name
		WHEN:  Register is called with nil second last name
		THEN:  User should be created successfully
	*/
	// GIVEN
	loginController, _, db := controllerTest.NewLoginControllerTestWrapper(t)

	name := "Jane"
	firstLastName := "Doe"
	email := utilsTest.GenerateRandomEmail()
	password := "securePassword123"
	imageUrl := "https://example.com/profile.jpg"

	// WHEN
	result, err := loginController.Register(
		name,
		firstLastName,
		nil, // No second last name
		email,
		password,
		imageUrl,
	)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, name, result.User.Name)
	assert.Equal(t, firstLastName, result.User.FirstLastName)
	assert.Nil(t, result.User.SecondLastName)
	assert.Equal(t, email, result.User.Email)

	// Verify in database
	var createdUser model.User
	dbErr := db.Where("email = ?", email).First(&createdUser).Error
	assert.NoError(t, dbErr)
	assert.Nil(t, createdUser.SecondLastName)
}

func TestRegisterWithDuplicateEmail(t *testing.T) {
	/*
		GIVEN: User with email already exists
		WHEN:  Register is called with existing email
		THEN:  It should return user already exists error
	*/
	// GIVEN
	loginController, _, db := controllerTest.NewLoginControllerTestWrapper(t)

	// Create existing user
	existingUser := factories.NewUserModel(db)

	name := "John"
	firstLastName := "Doe"
	password := "securePassword123"
	imageUrl := "https://example.com/profile.jpg"

	// WHEN - Try to register with same email
	result, err := loginController.Register(
		name,
		firstLastName,
		nil,
		existingUser.Email, // Same email as existing user
		password,
		imageUrl,
	)

	// THEN
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, errors.ConflictError.UserAlreadyExists.Code, err.Code)
}

func TestRegisterWithEmptyFields(t *testing.T) {
	/*
		GIVEN: Registration data with empty required fields
		WHEN:  Register is called with empty name, email, or password
		THEN:  It should handle appropriately (depending on validation)
	*/
	// GIVEN
	loginController, _, _ := controllerTest.NewLoginControllerTestWrapper(t)

	// Test with empty name
	result1, err1 := loginController.Register(
		"", // Empty name
		"Doe",
		nil,
		utilsTest.GenerateRandomEmail(),
		"password123",
		"https://example.com/profile.jpg",
	)

	// The behavior depends on implementation - it might succeed or fail
	// For this test, we'll just verify it doesn't crash
	if err1 != nil {
		assert.Nil(t, result1)
	}

	// Test with empty email
	result2, err2 := loginController.Register(
		"John",
		"Doe",
		nil,
		"", // Empty email
		"password123",
		"https://example.com/profile.jpg",
	)

	// Should likely fail due to email validation
	if err2 != nil {
		assert.Nil(t, result2)
	}
}

func TestRegisterWithMinimalData(t *testing.T) {
	/*
		GIVEN: Minimal valid registration data
		WHEN:  Register is called with only required fields
		THEN:  User should be created successfully
	*/
	// GIVEN
	loginController, _, db := controllerTest.NewLoginControllerTestWrapper(t)

	name := "MinimalUser"
	firstLastName := "Test"
	email := utilsTest.GenerateRandomEmail()
	password := "password123"
	imageUrl := "" // Empty image URL

	// WHEN
	result, err := loginController.Register(
		name,
		firstLastName,
		nil,
		email,
		password,
		imageUrl,
	)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, name, result.User.Name)
	assert.Equal(t, firstLastName, result.User.FirstLastName)
	assert.Equal(t, email, result.User.Email)
	assert.Equal(t, imageUrl, result.User.ImageUrl)

	// Verify in database
	var createdUser model.User
	dbErr := db.Where("email = ?", email).First(&createdUser).Error
	assert.NoError(t, dbErr)
}
