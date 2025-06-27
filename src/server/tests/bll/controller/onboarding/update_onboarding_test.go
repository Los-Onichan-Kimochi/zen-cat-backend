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

func TestUpdateOnboardingSuccessfully(t *testing.T) {
	/*
		GIVEN: An onboarding record exists in the database
		WHEN:  UpdateOnboarding is called with valid data
		THEN:  The onboarding record should be updated successfully
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

	// Create an onboarding record
	onboarding := &model.Onboarding{
		UserId:         user.Id,
		DocumentType:   model.DocumentTypeDni,
		DocumentNumber: "12345678",
		PhoneNumber:    "987654321",
		PostalCode:     "15001",
		Address:        "Av. Principal 123",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err = db.Create(onboarding).Error
	assert.NoError(t, err)

	// Prepare update data
	newPhoneNumber := "999888777"
	newPostalCode := "15003"
	updateRequest := schemas.UpdateOnboardingRequest{
		PhoneNumber: &newPhoneNumber,
		PostalCode:  &newPostalCode,
	}
	updatedBy := "ADMIN"

	// WHEN
	result, errResult := onboardingController.UpdateOnboarding(
		onboarding.Id,
		updateRequest,
		updatedBy,
	)

	// THEN
	assert.Nil(t, errResult)
	assert.NotNil(t, result)
	assert.Equal(t, onboarding.Id, result.Id)
	assert.Equal(t, newPhoneNumber, result.PhoneNumber)
	assert.Equal(t, newPostalCode, result.PostalCode)
	assert.Equal(t, "12345678", result.DocumentNumber) // Should remain unchanged
}

func TestUpdateOnboardingNotFound(t *testing.T) {
	/*
		GIVEN: No onboarding record exists with the given ID
		WHEN:  UpdateOnboarding is called with non-existent ID
		THEN:  An error should be returned
	*/
	// GIVEN
	onboardingController, _, _ := controllerTest.NewOnboardingControllerTestWrapper(t)
	nonExistentId := uuid.New()
	updateRequest := schemas.UpdateOnboardingRequest{}
	updatedBy := "ADMIN"

	// WHEN
	result, errResult := onboardingController.UpdateOnboarding(
		nonExistentId,
		updateRequest,
		updatedBy,
	)

	// THEN
	assert.NotNil(t, errResult)
	assert.Nil(t, result)
	assert.Contains(t, errResult.Message, "not found")
}

func TestUpdateOnboardingByUserIdSuccessfully(t *testing.T) {
	/*
		GIVEN: An onboarding record exists for a user
		WHEN:  UpdateOnboardingByUserId is called with valid data
		THEN:  The onboarding record should be updated successfully
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

	// Create an onboarding record
	onboarding := &model.Onboarding{
		UserId:         user.Id,
		DocumentType:   model.DocumentTypeDni,
		DocumentNumber: "12345678",
		PhoneNumber:    "987654321",
		PostalCode:     "15001",
		Address:        "Av. Principal 123",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err = db.Create(onboarding).Error
	assert.NoError(t, err)

	// Prepare update data
	newAddress := "Av. Nueva 789"
	updateRequest := schemas.UpdateOnboardingRequest{
		Address: &newAddress,
	}
	updatedBy := "ADMIN"

	// WHEN
	result, errResult := onboardingController.UpdateOnboardingByUserId(
		user.Id,
		updateRequest,
		updatedBy,
	)

	// THEN
	assert.Nil(t, errResult)
	assert.NotNil(t, result)
	assert.Equal(t, user.Id, result.UserId)
	assert.Equal(t, newAddress, result.Address)
	assert.Equal(t, "12345678", result.DocumentNumber) // Should remain unchanged
}

func TestUpdateOnboardingByUserIdNotFound(t *testing.T) {
	/*
		GIVEN: No onboarding record exists for the given user ID
		WHEN:  UpdateOnboardingByUserId is called with non-existent user ID
		THEN:  An error should be returned
	*/
	// GIVEN
	onboardingController, _, _ := controllerTest.NewOnboardingControllerTestWrapper(t)
	nonExistentUserId := uuid.New()
	updateRequest := schemas.UpdateOnboardingRequest{}
	updatedBy := "ADMIN"

	// WHEN
	result, errResult := onboardingController.UpdateOnboardingByUserId(
		nonExistentUserId,
		updateRequest,
		updatedBy,
	)

	// THEN
	assert.NotNil(t, errResult)
	assert.Nil(t, result)
	assert.Contains(t, errResult.Message, "not found")
}
