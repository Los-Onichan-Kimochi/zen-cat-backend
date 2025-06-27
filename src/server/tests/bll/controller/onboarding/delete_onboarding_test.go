package onboarding_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

func TestDeleteOnboardingSuccessfully(t *testing.T) {
	/*
		GIVEN: An onboarding record exists in the database
		WHEN:  DeleteOnboarding is called with valid ID
		THEN:  The onboarding record should be deleted successfully
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

	// WHEN
	errResult := onboardingController.DeleteOnboarding(onboarding.Id)

	// THEN
	assert.Nil(t, errResult)

	// Verify the onboarding was deleted
	var deletedOnboarding model.Onboarding
	err = db.Where("id = ?", onboarding.Id).First(&deletedOnboarding).Error
	assert.Error(t, err) // Should return error because record doesn't exist
}

func TestDeleteOnboardingNotFound(t *testing.T) {
	/*
		GIVEN: No onboarding record exists with the given ID
		WHEN:  DeleteOnboarding is called with non-existent ID
		THEN:  An error should be returned
	*/
	// GIVEN
	onboardingController, _, _ := controllerTest.NewOnboardingControllerTestWrapper(t)
	nonExistentId := uuid.New()

	// WHEN
	errResult := onboardingController.DeleteOnboarding(nonExistentId)

	// THEN
	assert.NotNil(t, errResult)
	assert.Contains(t, errResult.Message, "not found")
}

func TestDeleteOnboardingByUserIdSuccessfully(t *testing.T) {
	/*
		GIVEN: An onboarding record exists for a user
		WHEN:  DeleteOnboardingByUserId is called with valid user ID
		THEN:  The onboarding record should be deleted successfully
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

	// WHEN
	errResult := onboardingController.DeleteOnboardingByUserId(user.Id)

	// THEN
	assert.Nil(t, errResult)

	// Verify the onboarding was deleted
	var deletedOnboarding model.Onboarding
	err = db.Where("user_id = ?", user.Id).First(&deletedOnboarding).Error
	assert.Error(t, err) // Should return error because record doesn't exist
}

func TestDeleteOnboardingByUserIdNotFound(t *testing.T) {
	/*
		GIVEN: No onboarding record exists for the given user ID
		WHEN:  DeleteOnboardingByUserId is called with non-existent user ID
		THEN:  An error should be returned
	*/
	// GIVEN
	onboardingController, _, _ := controllerTest.NewOnboardingControllerTestWrapper(t)
	nonExistentUserId := uuid.New()

	// WHEN
	errResult := onboardingController.DeleteOnboardingByUserId(nonExistentUserId)

	// THEN
	assert.NotNil(t, errResult)
	assert.Contains(t, errResult.Message, "not found")
}
