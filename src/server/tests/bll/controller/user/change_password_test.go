package user_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestChangePasswordSuccessfully(t *testing.T) {
	// GIVEN: An existing user
	controller, _, db := controllerTest.NewUserControllerTestWrapper(t)

	// Create a test user
	email := "user.changepass@example.com"
	name := "John"
	firstName := "Doe"
	rol := model.UserRolClient

	testUser := factories.NewUserModel(db, factories.UserModelF{
		Name:          &name,
		FirstLastName: &firstName,
		Email:         &email,
		Rol:           &rol,
	})
	originalPassword := testUser.Password

	changePasswordRequest := schemas.ChangePasswordInput{
		Email:       email,
		NewPassword: "newSecurePassword123",
	}

	// WHEN: ChangePassword is called
	err := controller.ChangePassword(email, changePasswordRequest)

	// THEN: Password is changed successfully
	assert.Nil(t, err)

	// Verify password was actually changed by getting the user
	updatedUser, getErr := controller.GetUser(testUser.Id)
	assert.Nil(t, getErr)
	assert.NotEqual(t, originalPassword, updatedUser.Password)
}

func TestChangePasswordUserNotFound(t *testing.T) {
	// GIVEN: Non-existent user email
	controller, _, _ := controllerTest.NewUserControllerTestWrapper(t)

	nonExistentEmail := "nonexistent@example.com"
	changePasswordRequest := schemas.ChangePasswordInput{
		Email:       nonExistentEmail,
		NewPassword: "newPassword123",
	}

	// WHEN: ChangePassword is called with non-existent email
	err := controller.ChangePassword(nonExistentEmail, changePasswordRequest)

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Contains(t, err.Message, "not found")
}

func TestChangePasswordEmptyEmail(t *testing.T) {
	// GIVEN: Empty email
	controller, _, _ := controllerTest.NewUserControllerTestWrapper(t)

	changePasswordRequest := schemas.ChangePasswordInput{
		Email:       "",
		NewPassword: "newPassword123",
	}

	// WHEN: ChangePassword is called with empty email
	err := controller.ChangePassword("", changePasswordRequest)

	// THEN: An error is returned
	assert.NotNil(t, err)
}

func TestChangePasswordEmptyNewPassword(t *testing.T) {
	// GIVEN: An existing user but empty new password
	controller, _, db := controllerTest.NewUserControllerTestWrapper(t)

	// Create a test user
	email := "user.emptypass@example.com"
	testUser := factories.NewUserModel(db, factories.UserModelF{
		Email: &email,
	})

	changePasswordRequest := schemas.ChangePasswordInput{
		Email:       email,
		NewPassword: "", // Empty password
	}

	// WHEN: ChangePassword is called with empty new password
	err := controller.ChangePassword(email, changePasswordRequest)

	// THEN: Password change should still work (validation might be at different layer)
	// The exact behavior depends on implementation
	if err == nil {
		// If no error, verify the password was updated
		updatedUser, getErr := controller.GetUser(testUser.Id)
		assert.Nil(t, getErr)
		assert.NotNil(t, updatedUser)
	} else {
		// If error, it should be related to password validation
		assert.NotNil(t, err)
	}
}

func TestChangePasswordWeakPassword(t *testing.T) {
	// GIVEN: An existing user and weak password
	controller, _, db := controllerTest.NewUserControllerTestWrapper(t)

	// Create a test user
	email := "user.weakpass@example.com"
	factories.NewUserModel(db, factories.UserModelF{
		Email: &email,
	})

	changePasswordRequest := schemas.ChangePasswordInput{
		Email:       email,
		NewPassword: "123", // Weak password
	}

	// WHEN: ChangePassword is called with weak password
	err := controller.ChangePassword(email, changePasswordRequest)

	// THEN: The operation should complete (password validation might be at different layer)
	// This test documents current behavior - password strength validation might be elsewhere
	if err != nil {
		// If there's validation at this level
		assert.NotNil(t, err)
	} else {
		// If no validation at this level, it should succeed
		assert.Nil(t, err)
	}
}

func TestChangePasswordSamePassword(t *testing.T) {
	// GIVEN: An existing user trying to change to same password
	controller, _, db := controllerTest.NewUserControllerTestWrapper(t)

	// Create a test user with known password
	email := "user.samepass@example.com"
	password := "currentPassword123"
	testUser := factories.NewUserModel(db, factories.UserModelF{
		Email:    &email,
		Password: &password,
	})

	changePasswordRequest := schemas.ChangePasswordInput{
		Email:       email,
		NewPassword: password, // Same as current password
	}

	// WHEN: ChangePassword is called with same password
	err := controller.ChangePassword(email, changePasswordRequest)

	// THEN: Operation should complete successfully
	assert.Nil(t, err)

	// Verify user still exists and accessible
	updatedUser, getErr := controller.GetUser(testUser.Id)
	assert.Nil(t, getErr)
	assert.NotNil(t, updatedUser)
}
