package onboarding_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

func TestGetOnboardingSuccessfully(t *testing.T) {
	/*
		GIVEN: An onboarding record exists in the database
		WHEN:  GetOnboarding is called with valid ID
		THEN:  The onboarding record should be returned
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
	result, errResult := onboardingController.GetOnboarding(onboarding.Id)

	// THEN
	assert.Nil(t, errResult)
	assert.NotNil(t, result)
	assert.Equal(t, onboarding.Id, result.Id)
	assert.Equal(t, onboarding.UserId, result.UserId)
	assert.Equal(t, string(onboarding.DocumentType), string(result.DocumentType))
	assert.Equal(t, onboarding.DocumentNumber, result.DocumentNumber)
	assert.Equal(t, onboarding.PhoneNumber, result.PhoneNumber)
}

func TestGetOnboardingNotFound(t *testing.T) {
	/*
		GIVEN: No onboarding record exists with the given ID
		WHEN:  GetOnboarding is called with non-existent ID
		THEN:  An error should be returned
	*/
	// GIVEN
	onboardingController, _, _ := controllerTest.NewOnboardingControllerTestWrapper(t)
	nonExistentId := uuid.New()

	// WHEN
	result, errResult := onboardingController.GetOnboarding(nonExistentId)

	// THEN
	assert.NotNil(t, errResult)
	assert.Nil(t, result)
	assert.Contains(t, errResult.Message, "not found")
}

func TestGetOnboardingWithNilId(t *testing.T) {
	/*
		GIVEN: A nil UUID
		WHEN:  GetOnboarding is called with nil UUID
		THEN:  An error should be returned
	*/
	// GIVEN
	onboardingController, _, _ := controllerTest.NewOnboardingControllerTestWrapper(t)
	nilId := uuid.Nil

	// WHEN
	result, errResult := onboardingController.GetOnboarding(nilId)

	// THEN
	assert.NotNil(t, errResult)
	assert.Nil(t, result)
}
