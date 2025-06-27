package professional_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	adapterTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/adapter"
)

func TestProfessionalAdapter_CreatePostgresqlProfessional_Success(t *testing.T) {
	// GIVEN: Valid professional data
	professionalAdapter, _, _ := adapterTest.NewProfessionalAdapterTestWrapper(t)

	secondLastName := "Second Last Name"

	// WHEN: CreatePostgresqlProfessional is called
	result, err := professionalAdapter.CreatePostgresqlProfessional(
		"John",
		"Doe",
		&secondLastName,
		"Therapist",
		"john.doe@example.com",
		"987654321",
		"Licensed",
		"https://example.com/john.jpg",
		"test_user",
	)

	// THEN: A new professional is created and returned
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "John", result.Name)
	assert.Equal(t, "Doe", result.FirstLastName)
	assert.Equal(t, "Second Last Name", *result.SecondLastName)
	assert.Equal(t, "Therapist", result.Specialty)
	assert.Equal(t, "john.doe@example.com", result.Email)
	assert.Equal(t, "987654321", result.PhoneNumber)
	assert.Equal(t, "Licensed", result.Type)
	assert.Equal(t, "https://example.com/john.jpg", result.ImageUrl)
	assert.NotEqual(t, "", result.Id)
}

func TestProfessionalAdapter_CreatePostgresqlProfessional_EmptyUpdatedBy(t *testing.T) {
	// GIVEN: Valid professional data but empty updatedBy
	professionalAdapter, _, _ := adapterTest.NewProfessionalAdapterTestWrapper(t)

	secondLastName := "Second Last Name"

	// WHEN: CreatePostgresqlProfessional is called
	result, err := professionalAdapter.CreatePostgresqlProfessional(
		"John",
		"Doe",
		&secondLastName,
		"Therapist",
		"john.doe@example.com",
		"987654321",
		"Licensed",
		"https://example.com/john.jpg",
		"",
	)

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "Invalid updated by value", err.Message)
}

func TestProfessionalAdapter_CreatePostgresqlProfessional_WithoutSecondLastName(t *testing.T) {
	// GIVEN: Valid professional data but without second last name
	professionalAdapter, _, db := adapterTest.NewProfessionalAdapterTestWrapper(t)

	// Using database factory for a real user
	testUser := factories.NewUserModel(db, factories.UserModelF{})

	// WHEN: CreatePostgresqlProfessional is called
	result, err := professionalAdapter.CreatePostgresqlProfessional(
		"Jane",
		"Smith",
		nil,
		"Psychologist",
		"jane.smith@example.com",
		"123456789",
		"Certified",
		"https://example.com/jane.jpg",
		testUser.Name,
	)

	// THEN: A new professional is created and returned
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Jane", result.Name)
	assert.Equal(t, "Smith", result.FirstLastName)
	assert.Nil(t, result.SecondLastName)
	assert.Equal(t, "Psychologist", result.Specialty)
	assert.Equal(t, "jane.smith@example.com", result.Email)
	assert.Equal(t, "123456789", result.PhoneNumber)
	assert.Equal(t, "Certified", result.Type)
	assert.Equal(t, "https://example.com/jane.jpg", result.ImageUrl)
	assert.NotEqual(t, "", result.Id)
}
