package onboarding_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

func TestCreateOnboardingForUserSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid onboarding data and existing user
		WHEN:  CreateOnboardingForUser is called
		THEN:  The onboarding record should be created successfully
	*/
	// GIVEN
	onboardingController, _, db := controllerTest.NewOnboardingControllerTestWrapper(t)

	// Create a user for the onboarding
	user := &model.User{
		Email:         utilsTest.GenerateRandomEmail(),
		Name:          "John",
		FirstLastName: "Doe",
		Rol:           model.UserRolClient,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err := db.Create(user).Error
	assert.NoError(t, err)

	// Prepare onboarding data
	district := "Lima"
	province := "Lima"
	region := "Lima"
	address := "Av. Principal 123"
	updatedBy := "ADMIN"

	// WHEN
	createRequest := schemas.CreateOnboardingRequest{
		DocumentType:   schemas.DocumentTypeDNI,
		DocumentNumber: "12345678",
		PhoneNumber:    "987654321",
		PostalCode:     "15001",
		Address:        address,
		District:       &district,
		Province:       &province,
		Region:         &region,
	}
	result, errResult := onboardingController.CreateOnboardingForUser(
		user.Id,
		createRequest,
		updatedBy,
	)

	// THEN
	assert.Nil(t, errResult)
	assert.NotNil(t, result)
	assert.Equal(t, user.Id, result.UserId)
	assert.Equal(t, string(schemas.DocumentTypeDNI), string(result.DocumentType))
	assert.Equal(t, "12345678", result.DocumentNumber)
	assert.Equal(t, "987654321", result.PhoneNumber)
	assert.Equal(t, "15001", result.PostalCode)

	// Verify the onboarding was created in the database
	var onboarding model.Onboarding
	err = db.Where("user_id = ?", user.Id).First(&onboarding).Error
	assert.NoError(t, err)
	assert.Equal(t, user.Id, onboarding.UserId)
	assert.Equal(t, model.DocumentTypeDni, onboarding.DocumentType)
	assert.Equal(t, "12345678", onboarding.DocumentNumber)
}

func TestCreateOnboardingForUserWithNonExistentUser(t *testing.T) {
	/*
		GIVEN: Non-existent user ID
		WHEN:  CreateOnboardingForUser is called with invalid user ID
		THEN:  An error should be returned
	*/
	// GIVEN
	onboardingController, _, _ := controllerTest.NewOnboardingControllerTestWrapper(t)
	nonExistentUserId := uuid.New()
	updatedBy := "ADMIN"

	// WHEN
	createRequest := schemas.CreateOnboardingRequest{
		DocumentType:   schemas.DocumentTypeDNI,
		DocumentNumber: "12345678",
		PhoneNumber:    "987654321",
		PostalCode:     "15001",
		Address:        "Av. Principal 123",
	}
	result, errResult := onboardingController.CreateOnboardingForUser(
		nonExistentUserId,
		createRequest,
		updatedBy,
	)

	// THEN
	assert.NotNil(t, errResult)
	assert.Nil(t, result)
}

func TestCreateOnboardingForUserWithInvalidDocumentType(t *testing.T) {
	/*
		GIVEN: Invalid document type
		WHEN:  CreateOnboardingForUser is called with invalid document type
		THEN:  An error should be returned
	*/
	// GIVEN
	onboardingController, _, db := controllerTest.NewOnboardingControllerTestWrapper(t)

	// Create a user for the onboarding
	user := &model.User{
		Email:         utilsTest.GenerateRandomEmail(),
		Name:          "John",
		FirstLastName: "Doe",
		Rol:           model.UserRolClient,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err := db.Create(user).Error
	assert.NoError(t, err)

	updatedBy := "ADMIN"

	// WHEN - Using invalid document type (empty string)
	createRequest := schemas.CreateOnboardingRequest{
		DocumentType:   schemas.DocumentType(""), // Invalid document type
		DocumentNumber: "12345678",
		PhoneNumber:    "987654321",
		PostalCode:     "15001",
		Address:        "Av. Principal 123",
	}
	result, errResult := onboardingController.CreateOnboardingForUser(
		user.Id,
		createRequest,
		updatedBy,
	)

	// THEN
	assert.NotNil(t, errResult)
	assert.Nil(t, result)
}

func TestCreateOnboardingForUserWithMinimalData(t *testing.T) {
	/*
		GIVEN: Minimal required onboarding data
		WHEN:  CreateOnboardingForUser is called with only required fields
		THEN:  The onboarding record should be created successfully
	*/
	// GIVEN
	onboardingController, _, db := controllerTest.NewOnboardingControllerTestWrapper(t)

	// Create a user for the onboarding
	user := &model.User{
		Email:         utilsTest.GenerateRandomEmail(),
		Name:          "John",
		FirstLastName: "Doe",
		Rol:           model.UserRolClient,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err := db.Create(user).Error
	assert.NoError(t, err)

	updatedBy := "ADMIN"

	// WHEN - Using only required fields
	createRequest := schemas.CreateOnboardingRequest{
		DocumentType:   schemas.DocumentTypeDNI,
		DocumentNumber: "12345678",
		PhoneNumber:    "987654321",
		PostalCode:     "15001",
		Address:        "Av. Principal 123",
	}
	result, errResult := onboardingController.CreateOnboardingForUser(
		user.Id,
		createRequest,
		updatedBy,
	)

	// THEN
	assert.Nil(t, errResult)
	assert.NotNil(t, result)
	assert.Equal(t, user.Id, result.UserId)
	assert.Equal(t, string(schemas.DocumentTypeDNI), string(result.DocumentType))
	assert.Equal(t, "12345678", result.DocumentNumber)
	assert.Equal(t, "987654321", result.PhoneNumber)
	assert.Equal(t, "15001", result.PostalCode)
}
