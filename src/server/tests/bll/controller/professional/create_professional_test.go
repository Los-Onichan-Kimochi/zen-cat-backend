package professional_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

func TestCreateProfessionalSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid professional data
		WHEN:  CreateProfessional is called
		THEN:  A new professional record should be created
	*/
	// GIVEN
	professionalController, _, db := controllerTest.NewProfessionalControllerTestWrapper(t)
	updatedBy := "ADMIN"

	createRequest := schemas.CreateProfessionalRequest{
		Name:           "Dr. John",
		FirstLastName:  "Doe",
		SecondLastName: "Smith",
		Specialty:      "Cardiology",
		Email:          utilsTest.GenerateRandomEmail(),
		PhoneNumber:    "987654321",
		Type:           string(model.ProfessionalTypeMedic),
		ImageUrl:       "https://example.com/image.jpg",
	}

	// WHEN
	result, errResult := professionalController.CreateProfessional(
		createRequest,
		updatedBy,
	)

	// THEN
	assert.Nil(t, errResult)
	assert.NotNil(t, result)
	assert.Equal(t, "Dr. John", result.Name)
	assert.Equal(t, "Doe", result.FirstLastName)
	assert.Equal(t, "Cardiology", result.Specialty)
	assert.Equal(t, string(model.ProfessionalTypeMedic), result.Type)

	// Verify the professional was created in the database
	var professional model.Professional
	err := db.Where("id = ?", result.Id).First(&professional).Error
	assert.NoError(t, err)
	assert.Equal(t, "Dr. John", professional.Name)
}

func TestCreateProfessionalWithEmptyName(t *testing.T) {
	/*
		GIVEN: Professional data with empty name
		WHEN:  CreateProfessional is called
		THEN:  An error should be returned
	*/
	// GIVEN
	professionalController, _, _ := controllerTest.NewProfessionalControllerTestWrapper(t)
	updatedBy := "ADMIN"

	createRequest := schemas.CreateProfessionalRequest{
		Name:          "", // Empty name
		FirstLastName: "Doe",
		Specialty:     "Cardiology",
		Email:         utilsTest.GenerateRandomEmail(),
		PhoneNumber:   "987654321",
		Type:          string(model.ProfessionalTypeMedic),
		ImageUrl:      "https://example.com/image.jpg",
	}

	// WHEN
	result, errResult := professionalController.CreateProfessional(
		createRequest,
		updatedBy,
	)

	// THEN
	assert.NotNil(t, errResult)
	assert.Nil(t, result)
}

func TestCreateProfessionalWithDuplicateEmail(t *testing.T) {
	/*
		GIVEN: A professional with the same email already exists
		WHEN:  CreateProfessional is called with duplicate email
		THEN:  An error should be returned
	*/
	// GIVEN
	professionalController, _, db := controllerTest.NewProfessionalControllerTestWrapper(t)
	updatedBy := "ADMIN"
	email := utilsTest.GenerateRandomEmail()

	// Create first professional
	professional1 := &model.Professional{
		Name:          "Dr. John",
		FirstLastName: "Doe",
		Specialty:     "Cardiology",
		Email:         email,
		PhoneNumber:   "987654321",
		Type:          model.ProfessionalTypeMedic,
		ImageUrl:      "https://example.com/image.jpg",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err := db.Create(professional1).Error
	assert.NoError(t, err)

	// WHEN - Try to create another professional with the same email
	createRequest := schemas.CreateProfessionalRequest{
		Name:          "Dr. Jane",
		FirstLastName: "Smith",
		Specialty:     "Neurology",
		Email:         email, // Same email
		PhoneNumber:   "123456789",
		Type:          string(model.ProfessionalTypeMedic),
		ImageUrl:      "https://example.com/jane.jpg",
	}
	result, errResult := professionalController.CreateProfessional(
		createRequest,
		updatedBy,
	)

	// THEN
	assert.NotNil(t, errResult)
	assert.Nil(t, result)
	assert.Contains(t, errResult.Message, "duplicate")
}
