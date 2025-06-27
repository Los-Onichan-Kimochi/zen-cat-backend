package auth_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestGenerateTokenSuccessfully(t *testing.T) {
	// GIVEN: A valid user and token parameters
	controller, _, db := controllerTest.NewAuthControllerTestWrapper(t)

	// Create a test user
	email := "user.token@example.com"
	password := "securePassword123"
	name := "John"
	firstName := "Doe"
	rol := model.UserRolClient

	testUser := factories.NewUserModel(db, factories.UserModelF{
		Name:          &name,
		FirstLastName: &firstName,
		Email:         &email,
		Password:      &password,
		Rol:           &rol,
	})

	userRoles := []string{"CLIENT"}
	expirationDelta := time.Hour * 2

	// WHEN: GenerateToken is called
	result, err := controller.GenerateToken(
		testUser.Id,
		testUser.Email,
		testUser.Password,
		userRoles,
		expirationDelta,
	)

	// THEN: Token is generated successfully
	assert.Nil(t, err)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	assert.Equal(t, expirationDelta, result.ExpiresIn)
	assert.Equal(t, result.AccessToken, result.RefreshToken) // Based on implementation
}

func TestGenerateTokenUserNotFound(t *testing.T) {
	// GIVEN: A non-existent user ID
	controller, _, _ := controllerTest.NewAuthControllerTestWrapper(t)

	nonExistentUserId := uuid.New()
	userEmail := "nonexistent@example.com"
	userPassword := "password123"
	userRoles := []string{"CLIENT"}
	expirationDelta := time.Hour * 1

	// WHEN: GenerateToken is called with non-existent user
	result, err := controller.GenerateToken(
		nonExistentUserId,
		userEmail,
		userPassword,
		userRoles,
		expirationDelta,
	)

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Empty(t, result.AccessToken)
	assert.Empty(t, result.RefreshToken)
}

func TestGenerateTokenAdminRole(t *testing.T) {
	// GIVEN: A user with admin role
	controller, _, db := controllerTest.NewAuthControllerTestWrapper(t)

	// Create an admin user
	email := "admin.token@example.com"
	password := "adminPassword123"
	name := "Admin"
	firstName := "User"
	rol := model.UserRolAdmin

	testUser := factories.NewUserModel(db, factories.UserModelF{
		Name:          &name,
		FirstLastName: &firstName,
		Email:         &email,
		Password:      &password,
		Rol:           &rol,
	})

	userRoles := []string{"ADMINISTRATOR"}
	expirationDelta := time.Hour * 4

	// WHEN: GenerateToken is called for admin user
	result, err := controller.GenerateToken(
		testUser.Id,
		testUser.Email,
		testUser.Password,
		userRoles,
		expirationDelta,
	)

	// THEN: Token is generated successfully for admin
	assert.Nil(t, err)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	assert.Equal(t, expirationDelta, result.ExpiresIn)
}

func TestGenerateTokenMultipleRoles(t *testing.T) {
	// GIVEN: A user with multiple roles
	controller, _, db := controllerTest.NewAuthControllerTestWrapper(t)

	// Create a test user
	testUser := factories.NewUserModel(db, factories.UserModelF{})

	userRoles := []string{"CLIENT", "MODERATOR", "GUEST"}
	expirationDelta := time.Hour * 1

	// WHEN: GenerateToken is called with multiple roles
	result, err := controller.GenerateToken(
		testUser.Id,
		testUser.Email,
		testUser.Password,
		userRoles,
		expirationDelta,
	)

	// THEN: Token is generated successfully with multiple roles
	assert.Nil(t, err)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	assert.Equal(t, expirationDelta, result.ExpiresIn)
}

func TestGenerateTokenShortExpiration(t *testing.T) {
	// GIVEN: A valid user and short expiration time
	controller, _, db := controllerTest.NewAuthControllerTestWrapper(t)

	testUser := factories.NewUserModel(db, factories.UserModelF{})
	userRoles := []string{"CLIENT"}
	expirationDelta := time.Minute * 5 // Short expiration

	// WHEN: GenerateToken is called with short expiration
	result, err := controller.GenerateToken(
		testUser.Id,
		testUser.Email,
		testUser.Password,
		userRoles,
		expirationDelta,
	)

	// THEN: Token is generated successfully with short expiration
	assert.Nil(t, err)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	assert.Equal(t, expirationDelta, result.ExpiresIn)
}

func TestGenerateTokenLongExpiration(t *testing.T) {
	// GIVEN: A valid user and long expiration time
	controller, _, db := controllerTest.NewAuthControllerTestWrapper(t)

	testUser := factories.NewUserModel(db, factories.UserModelF{})
	userRoles := []string{"CLIENT"}
	expirationDelta := time.Hour * 24 * 7 // One week

	// WHEN: GenerateToken is called with long expiration
	result, err := controller.GenerateToken(
		testUser.Id,
		testUser.Email,
		testUser.Password,
		userRoles,
		expirationDelta,
	)

	// THEN: Token is generated successfully with long expiration
	assert.Nil(t, err)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	assert.Equal(t, expirationDelta, result.ExpiresIn)
}

func TestGenerateTokenEmptyRoles(t *testing.T) {
	// GIVEN: A valid user but empty roles
	controller, _, db := controllerTest.NewAuthControllerTestWrapper(t)

	testUser := factories.NewUserModel(db, factories.UserModelF{})
	userRoles := []string{} // Empty roles
	expirationDelta := time.Hour * 1

	// WHEN: GenerateToken is called with empty roles
	result, err := controller.GenerateToken(
		testUser.Id,
		testUser.Email,
		testUser.Password,
		userRoles,
		expirationDelta,
	)

	// THEN: Token is generated successfully (roles might be optional)
	assert.Nil(t, err)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	assert.Equal(t, expirationDelta, result.ExpiresIn)
}

func TestGenerateTokenZeroExpiration(t *testing.T) {
	// GIVEN: A valid user but zero expiration time
	controller, _, db := controllerTest.NewAuthControllerTestWrapper(t)

	testUser := factories.NewUserModel(db, factories.UserModelF{})
	userRoles := []string{"CLIENT"}
	expirationDelta := time.Duration(0) // Zero expiration

	// WHEN: GenerateToken is called with zero expiration
	result, err := controller.GenerateToken(
		testUser.Id,
		testUser.Email,
		testUser.Password,
		userRoles,
		expirationDelta,
	)

	// THEN: Token generation behavior depends on implementation
	// It might succeed with immediate expiration or fail
	if err == nil {
		assert.NotEmpty(t, result.AccessToken)
		assert.Equal(t, expirationDelta, result.ExpiresIn)
	} else {
		assert.NotNil(t, err)
		assert.Empty(t, result.AccessToken)
	}
}
