package onboarding_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestDeleteOnboardingSuccessfully(t *testing.T) {
	/*
		GIVEN: An onboarding record exists in the database
		WHEN:  DeleteOnboarding is called with valid ID
		THEN:  The onboarding record should be deleted successfully
	*/
	// GIVEN
	onboardingController, _, db := controllerTest.NewOnboardingControllerTestWrapper(t)

	// Create user and onboarding using factories
	user := factories.NewUserModel(db)
	onboarding := factories.NewOnboardingModel(db, factories.OnboardingModelF{
		UserId: &user.Id,
	})

	// WHEN
	errResult := onboardingController.DeleteOnboarding(onboarding.Id)

	// THEN
	assert.Nil(t, errResult)

	// Verify the onboarding was deleted
	var deletedOnboarding model.Onboarding
	err := db.Where("id = ?", onboarding.Id).First(&deletedOnboarding).Error
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

	// Create user and onboarding using factories
	user := factories.NewUserModel(db)
	_ = factories.NewOnboardingModel(db, factories.OnboardingModelF{
		UserId: &user.Id,
	})

	// WHEN
	errResult := onboardingController.DeleteOnboardingByUserId(user.Id)

	// THEN
	assert.Nil(t, errResult)

	// Verify the onboarding was deleted
	var deletedOnboarding model.Onboarding
	err := db.Where("user_id = ?", user.Id).First(&deletedOnboarding).Error
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
