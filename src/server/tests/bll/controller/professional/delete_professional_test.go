package professional_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

func TestDeleteProfessionalSuccessfully(t *testing.T) {
	/*
		GIVEN: A professional record exists in the database
		WHEN:  DeleteProfessional is called with valid ID
		THEN:  The professional record should be deleted successfully
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

	// WHEN
	errResult := professionalController.DeleteProfessional(professional.Id)

	// THEN
	assert.Nil(t, errResult)

	// Verify the professional was deleted
	var deletedProfessional model.Professional
	err = db.Where("id = ?", professional.Id).First(&deletedProfessional).Error
	assert.Error(t, err) // Should return error because record doesn't exist
}

func TestDeleteProfessionalNotFound(t *testing.T) {
	/*
		GIVEN: No professional record exists with the given ID
		WHEN:  DeleteProfessional is called with non-existent ID
		THEN:  An error should be returned
	*/
	// GIVEN
	professionalController, _, _ := controllerTest.NewProfessionalControllerTestWrapper(t)
	nonExistentId := uuid.New()

	// WHEN
	errResult := professionalController.DeleteProfessional(nonExistentId)

	// THEN
	assert.NotNil(t, errResult)
	assert.Contains(t, errResult.Message, "not soft deleted")
}

func TestDeleteProfessionalWithNilId(t *testing.T) {
	/*
		GIVEN: A nil UUID
		WHEN:  DeleteProfessional is called with nil UUID
		THEN:  An error should be returned
	*/
	// GIVEN
	professionalController, _, _ := controllerTest.NewProfessionalControllerTestWrapper(t)
	nilId := uuid.Nil

	// WHEN
	errResult := professionalController.DeleteProfessional(nilId)

	// THEN
	assert.NotNil(t, errResult)
}
