package onboarding_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	adapterTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/adapter"
)

func TestOnboardingAdapter_CreatePostgresqlOnboarding_Success(t *testing.T) {
	// GIVEN: Valid onboarding data with existing user
	onboardingAdapter, _, db := adapterTest.NewOnboardingAdapterTestWrapper(t)

	// Create a user first
	user := factories.NewUserModel(db, factories.UserModelF{})

	birthDate := "1990-05-15"
	gender := schemas.GenderMale
	district := "Lima"
	province := "Lima"
	region := "Lima"

	// WHEN: CreatePostgresqlOnboarding is called
	result, err := onboardingAdapter.CreatePostgresqlOnboarding(
		user.Id,
		schemas.DocumentTypeDNI,
		"12345678",
		"987654321",
		&birthDate,
		&gender,
		"15001",
		&district,
		&province,
		&region,
		"Av. Example 123",
		"test_user",
	)

	// THEN: A new onboarding is created and returned
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, schemas.DocumentTypeDNI, result.DocumentType)
	assert.Equal(t, "12345678", result.DocumentNumber)
	assert.Equal(t, "987654321", result.PhoneNumber)
	assert.NotNil(t, result.BirthDate)
	assert.Equal(t, schemas.GenderMale, *result.Gender)
	assert.Equal(t, "15001", result.PostalCode)
	assert.Equal(t, "Lima", *result.District)
	assert.Equal(t, "Lima", *result.Province)
	assert.Equal(t, "Lima", *result.Region)
	assert.Equal(t, "Av. Example 123", result.Address)
	assert.Equal(t, user.Id, result.UserId)
	assert.NotEqual(t, "", result.Id)
}

func TestOnboardingAdapter_CreatePostgresqlOnboarding_EmptyUpdatedBy(t *testing.T) {
	// GIVEN: Valid onboarding data but empty updatedBy
	onboardingAdapter, _, db := adapterTest.NewOnboardingAdapterTestWrapper(t)

	// Create a user first
	user := factories.NewUserModel(db, factories.UserModelF{})

	birthDate := "1990-05-15"
	gender := schemas.GenderFemale

	// WHEN: CreatePostgresqlOnboarding is called
	result, err := onboardingAdapter.CreatePostgresqlOnboarding(
		user.Id,
		schemas.DocumentTypePassport,
		"A1234567",
		"123456789",
		&birthDate,
		&gender,
		"15002",
		nil,
		nil,
		nil,
		"Av. Test 456",
		"",
	)

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "Invalid updated by value", err.Message)
}
