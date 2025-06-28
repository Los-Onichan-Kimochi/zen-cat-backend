package professional_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

func TestUpdateProfessionalSuccessfully(t *testing.T) {
	/*
		GIVEN: A professional record exists in the database
		WHEN:  UpdateProfessional is called with valid data
		THEN:  The professional record should be updated successfully
	*/
	// GIVEN
	professionalController, _, db := controllerTest.NewProfessionalControllerTestWrapper(t)

	// Create a professional record
	professional := &model.Professional{
		Name:          "Dr. John",
		FirstLastName: "Doe",
		Specialty:     "Cardiology",
		Email:         utilsTest.GenerateRandomEmail(),
		PhoneNumber:   "987654321",
		Type:          model.ProfessionalTypeMedic,
		ImageUrl:      "https://example.com/image.jpg",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err := db.Create(professional).Error
	assert.NoError(t, err)

	// Prepare update data
	newName := "Dr. Jane"
	newSpecialty := "Neurology"
	updateRequest := schemas.UpdateProfessionalRequest{
		Name:      &newName,
		Specialty: &newSpecialty,
	}
	updatedBy := "ADMIN"

	// WHEN
	result, errResult := professionalController.UpdateProfessional(
		professional.Id,
		updateRequest,
		updatedBy,
	)

	// THEN
	assert.Nil(t, errResult)
	assert.NotNil(t, result)
	assert.Equal(t, professional.Id, result.Id)
	assert.Equal(t, newName, result.Name)
	assert.Equal(t, newSpecialty, result.Specialty)
	assert.Equal(t, professional.FirstLastName, result.FirstLastName) // Should remain unchanged
}

func TestUpdateProfessionalNotFound(t *testing.T) {
	/*
		GIVEN: No professional record exists with the given ID
		WHEN:  UpdateProfessional is called with non-existent ID
		THEN:  An error should be returned
	*/
	// GIVEN
	professionalController, _, _ := controllerTest.NewProfessionalControllerTestWrapper(t)
	nonExistentId := uuid.New()
	updateRequest := schemas.UpdateProfessionalRequest{}
	updatedBy := "ADMIN"

	// WHEN
	result, errResult := professionalController.UpdateProfessional(
		nonExistentId,
		updateRequest,
		updatedBy,
	)

	// THEN
	assert.NotNil(t, errResult)
	assert.Nil(t, result)
	assert.Contains(t, errResult.Message, "not found")
}

func TestUpdateProfessionalWithPartialData(t *testing.T) {
	/*
		GIVEN: A professional record exists in the database
		WHEN:  UpdateProfessional is called with only some fields
		THEN:  Only the specified fields should be updated
	*/
	// GIVEN
	professionalController, _, db := controllerTest.NewProfessionalControllerTestWrapper(t)

	// Create a professional record
	professional := &model.Professional{
		Name:          "Dr. John",
		FirstLastName: "Doe",
		Specialty:     "Cardiology",
		Email:         utilsTest.GenerateRandomEmail(),
		PhoneNumber:   "987654321",
		Type:          model.ProfessionalTypeMedic,
		ImageUrl:      "https://example.com/image.jpg",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err := db.Create(professional).Error
	assert.NoError(t, err)

	// Prepare update data - only update phone number
	newPhoneNumber := "111222333"
	updateRequest := schemas.UpdateProfessionalRequest{
		PhoneNumber: &newPhoneNumber,
	}
	updatedBy := "ADMIN"

	// WHEN
	result, errResult := professionalController.UpdateProfessional(
		professional.Id,
		updateRequest,
		updatedBy,
	)

	// THEN
	assert.Nil(t, errResult)
	assert.NotNil(t, result)
	assert.Equal(t, professional.Id, result.Id)
	assert.Equal(t, newPhoneNumber, result.PhoneNumber)
	assert.Equal(t, professional.Name, result.Name)           // Should remain unchanged
	assert.Equal(t, professional.Specialty, result.Specialty) // Should remain unchanged
}
