package forgot_password_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestGenerateResetPinSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid user exists
		WHEN:  GenerateResetPin is called with valid email
		THEN:  Reset pin should be generated and returned
	*/
	// GIVEN
	forgotPasswordController, _, db := controllerTest.NewForgotPasswordControllerTestWrapper(t)

	// Create test user
	testUser := factories.NewUserModel(db)

	// WHEN
	result, err := forgotPasswordController.GenerateResetPin(testUser.Email)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Pin)
	assert.Equal(t, "Código enviado al correo", result.Message)
	assert.Len(t, result.Pin, 6) // PIN should be 6 digits

	// Verify PIN is numeric
	for _, char := range result.Pin {
		assert.True(t, char >= '0' && char <= '9', "PIN should contain only digits")
	}
}

func TestGenerateResetPinWithNonExistentUser(t *testing.T) {
	/*
		GIVEN: User does not exist
		WHEN:  GenerateResetPin is called with non-existent email
		THEN:  It should return user not found error
	*/
	// GIVEN
	forgotPasswordController, _, _ := controllerTest.NewForgotPasswordControllerTestWrapper(t)

	nonExistentEmail := "nonexistent@example.com"

	// WHEN
	result, err := forgotPasswordController.GenerateResetPin(nonExistentEmail)

	// THEN
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, errors.ObjectNotFoundError.UserNotFound.Code, err.Code)
}

func TestGenerateResetPinWithEmptyEmail(t *testing.T) {
	/*
		GIVEN: Empty email
		WHEN:  GenerateResetPin is called with empty email
		THEN:  It should return user not found error
	*/
	// GIVEN
	forgotPasswordController, _, _ := controllerTest.NewForgotPasswordControllerTestWrapper(t)

	emptyEmail := ""

	// WHEN
	result, err := forgotPasswordController.GenerateResetPin(emptyEmail)

	// THEN
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, errors.ObjectNotFoundError.UserNotFound.Code, err.Code)
}

func TestGenerateResetPinMultipleTimes(t *testing.T) {
	/*
		GIVEN: Valid user exists
		WHEN:  GenerateResetPin is called multiple times for the same user
		THEN:  Different PINs should be generated each time
	*/
	// GIVEN
	forgotPasswordController, _, db := controllerTest.NewForgotPasswordControllerTestWrapper(t)

	// Create test user
	testUser := factories.NewUserModel(db)

	// WHEN - Generate first PIN
	result1, err1 := forgotPasswordController.GenerateResetPin(testUser.Email)
	assert.Nil(t, err1)
	assert.NotNil(t, result1)

	// WHEN - Generate second PIN
	result2, err2 := forgotPasswordController.GenerateResetPin(testUser.Email)
	assert.Nil(t, err2)
	assert.NotNil(t, result2)

	// THEN
	assert.NotEqual(t, result1.Pin, result2.Pin, "Different PINs should be generated")
	assert.Len(t, result1.Pin, 6)
	assert.Len(t, result2.Pin, 6)
}
