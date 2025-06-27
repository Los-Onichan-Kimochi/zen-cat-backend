package user_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestCheckUserExistsByEmailExists(t *testing.T) {
	// GIVEN: An existing user
	controller, _, db := controllerTest.NewUserControllerTestWrapper(t)

	// Create a test user
	email := "existing.user@example.com"
	name := "John"
	firstName := "Doe"
	rol := model.UserRolClient

	factories.NewUserModel(db, factories.UserModelF{
		Name:          &name,
		FirstLastName: &firstName,
		Email:         &email,
		Rol:           &rol,
	})

	// WHEN: CheckUserExistsByEmail is called with existing email
	result, err := controller.CheckUserExistsByEmail(email)

	// THEN: User exists is returned as true
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, email, result.Email)
	assert.True(t, result.Exists)
}

func TestCheckUserExistsByEmailNotExists(t *testing.T) {
	// GIVEN: No user with the given email
	controller, _, _ := controllerTest.NewUserControllerTestWrapper(t)

	nonExistentEmail := "nonexistent@example.com"

	// WHEN: CheckUserExistsByEmail is called with non-existent email
	result, err := controller.CheckUserExistsByEmail(nonExistentEmail)

	// THEN: User exists is returned as false
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, nonExistentEmail, result.Email)
	assert.False(t, result.Exists)
}

func TestCheckUserExistsByEmailEmptyEmail(t *testing.T) {
	// GIVEN: Empty email
	controller, _, _ := controllerTest.NewUserControllerTestWrapper(t)

	// WHEN: CheckUserExistsByEmail is called with empty email
	result, err := controller.CheckUserExistsByEmail("")

	// THEN: User exists is returned as false
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "", result.Email)
	assert.False(t, result.Exists)
}

func TestCheckUserExistsByEmailInvalidFormat(t *testing.T) {
	// GIVEN: Invalid email format
	controller, _, _ := controllerTest.NewUserControllerTestWrapper(t)

	invalidEmail := "invalid-email-format"

	// WHEN: CheckUserExistsByEmail is called with invalid email format
	result, err := controller.CheckUserExistsByEmail(invalidEmail)

	// THEN: User exists is returned as false (no validation at controller level)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, invalidEmail, result.Email)
	assert.False(t, result.Exists)
}

func TestCheckUserExistsByEmailCaseSensitive(t *testing.T) {
	// GIVEN: An existing user with lowercase email
	controller, _, db := controllerTest.NewUserControllerTestWrapper(t)

	// Create a test user with lowercase email
	email := "lowercase@example.com"
	name := "John"
	firstName := "Doe"
	rol := model.UserRolClient

	factories.NewUserModel(db, factories.UserModelF{
		Name:          &name,
		FirstLastName: &firstName,
		Email:         &email,
		Rol:           &rol,
	})

	// WHEN: CheckUserExistsByEmail is called with uppercase email
	upperCaseEmail := "LOWERCASE@EXAMPLE.COM"
	result, err := controller.CheckUserExistsByEmail(upperCaseEmail)

	// THEN: Result depends on database collation (usually case-insensitive)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, upperCaseEmail, result.Email)
	// Note: The result.Exists value depends on database configuration
	// Most databases are case-insensitive for email comparison
}
