package onboarding_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

func TestFetchOnboardingsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple onboarding records exist in the database
		WHEN:  FetchOnboardings is called
		THEN:  All onboarding records should be returned
	*/
	// GIVEN
	onboardingController, _, db := controllerTest.NewOnboardingControllerTestWrapper(t)

	// Create users for the onboardings
	user1 := &model.User{
		Email:         utilsTest.GenerateRandomEmail(),
		Name:          "John",
		FirstLastName: "Doe",
		Rol:           model.UserRolClient,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	user2 := &model.User{
		Email:         utilsTest.GenerateRandomEmail(),
		Name:          "Jane",
		FirstLastName: "Smith",
		Rol:           model.UserRolClient,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err := db.Create([]*model.User{user1, user2}).Error
	assert.NoError(t, err)

	// Create onboarding records
	onboardings := []*model.Onboarding{
		{
			UserId:         user1.Id,
			DocumentType:   model.DocumentTypeDni,
			DocumentNumber: "12345678",
			PhoneNumber:    "987654321",
			PostalCode:     "15001",
			Address:        "Av. Principal 123",
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			UserId:         user2.Id,
			DocumentType:   model.DocumentTypeDni,
			DocumentNumber: "87654321",
			PhoneNumber:    "123456789",
			PostalCode:     "15002",
			Address:        "Av. Secundaria 456",
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
	}
	err = db.Create(onboardings).Error
	assert.NoError(t, err)

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
