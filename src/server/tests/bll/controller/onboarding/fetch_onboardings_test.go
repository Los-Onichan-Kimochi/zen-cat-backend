package onboarding_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestFetchOnboardingsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple onboarding records exist in the database
		WHEN:  FetchOnboardings is called
		THEN:  All onboarding records should be returned
	*/
	// GIVEN
	onboardingController, _, db := controllerTest.NewOnboardingControllerTestWrapper(t)

	// Create users and onboardings using factories
	user1 := factories.NewUserModel(db)
	user2 := factories.NewUserModel(db)

	_ = factories.NewOnboardingModel(db, factories.OnboardingModelF{
		UserId: &user1.Id,
	})
	_ = factories.NewOnboardingModel(db, factories.OnboardingModelF{
		UserId: &user2.Id,
	})

	// WHEN
	result, errResult := onboardingController.FetchOnboardings()

	// THEN
	assert.Nil(t, errResult)
	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, len(result.Onboardings), 2)
}

func TestFetchOnboardingsEmpty(t *testing.T) {
	/*
		GIVEN: No onboarding records exist in the database
		WHEN:  FetchOnboardings is called
		THEN:  An empty list should be returned
	*/
	// GIVEN
	onboardingController, _, _ := controllerTest.NewOnboardingControllerTestWrapper(t)

	// WHEN
	result, errResult := onboardingController.FetchOnboardings()

	// THEN
	assert.Nil(t, errResult)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result.Onboardings))
}
