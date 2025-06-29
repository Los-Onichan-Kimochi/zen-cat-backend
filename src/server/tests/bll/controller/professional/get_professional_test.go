package professional_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

func TestGetProfessionalSuccessfully(t *testing.T) {
	/*
		GIVEN: A professional record exists in the database
		WHEN:  GetProfessional is called with valid ID
		THEN:  The professional record should be returned
	*/
	// GIVEN
	professionalController, _, db := controllerTest.NewProfessionalControllerTestWrapper(t)

	// Create a professional record
	secondLastName := "Smith"
	professional := &model.Professional{
		Name:           "Dr. John",
		FirstLastName:  "Doe",
		SecondLastName: &secondLastName,
		Specialty:      "Cardiology",
		Email:          utilsTest.GenerateRandomEmail(),
		PhoneNumber:    "987654321",
		Type:           model.ProfessionalTypeMedic,
		ImageUrl:       "https://example.com/image.jpg",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err := db.Create(professional).Error
	assert.NoError(t, err)

	// WHEN
	result, errResult := professionalController.GetProfessional(professional.Id)

	// THEN
	assert.Nil(t, errResult)
	assert.NotNil(t, result)
	assert.Equal(t, professional.Id, result.Id)
	assert.Equal(t, professional.Name, result.Name)
	assert.Equal(t, professional.FirstLastName, result.FirstLastName)
	assert.Equal(t, professional.Specialty, result.Specialty)
	assert.Equal(t, professional.Email, result.Email)
}

func TestGetProfessionalNotFound(t *testing.T) {
	/*
		GIVEN: No professional record exists with the given ID
		WHEN:  GetProfessional is called with non-existent ID
		THEN:  An error should be returned
	*/
	// GIVEN
	professionalController, _, _ := controllerTest.NewProfessionalControllerTestWrapper(t)
	nonExistentId := uuid.New()

	// WHEN
	result, errResult := professionalController.GetProfessional(nonExistentId)

	// THEN
	assert.NotNil(t, errResult)
	assert.Nil(t, result)
	assert.Contains(t, errResult.Message, "not found")
}

func TestGetProfessionalWithNilId(t *testing.T) {
	/*
		GIVEN: A nil UUID
		WHEN:  GetProfessional is called with nil UUID
		THEN:  An error should be returned
	*/
	// GIVEN
	professionalController, _, _ := controllerTest.NewProfessionalControllerTestWrapper(t)
	nilId := uuid.Nil

	// WHEN
	result, errResult := professionalController.GetProfessional(nilId)

	// THEN
	assert.NotNil(t, errResult)
	assert.Nil(t, result)
}
