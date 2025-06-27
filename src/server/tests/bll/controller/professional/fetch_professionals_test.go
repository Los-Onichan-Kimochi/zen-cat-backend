package professional_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

func TestFetchProfessionalsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple professional records exist in the database
		WHEN:  FetchProfessionals is called
		THEN:  All professional records should be returned
	*/
	// GIVEN
	professionalController, _, db := controllerTest.NewProfessionalControllerTestWrapper(t)

	// Create professional records
	professionals := []*model.Professional{
		{
			Name:          "Dr. John",
			FirstLastName: "Doe",
			Specialty:     "Cardiology",
			Email:         utilsTest.GenerateRandomEmail(),
			PhoneNumber:   "987654321",
			Type:          model.ProfessionalTypeMedic,
			ImageUrl:      "https://example.com/john.jpg",
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Name:          "Jane",
			FirstLastName: "Smith",
			Specialty:     "Fitness Training",
			Email:         utilsTest.GenerateRandomEmail(),
			PhoneNumber:   "123456789",
			Type:          model.ProfessionalTypeGymTrainer,
			ImageUrl:      "https://example.com/jane.jpg",
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
	}
	err := db.Create(professionals).Error
	assert.NoError(t, err)

	// WHEN
	result, errResult := professionalController.FetchProfessionals()

	// THEN
	assert.Nil(t, errResult)
	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, len(result.Professionals), 2)
}

func TestFetchProfessionalsEmpty(t *testing.T) {
	/*
		GIVEN: No professional records exist in the database
		WHEN:  FetchProfessionals is called
		THEN:  An empty list should be returned
	*/
	// GIVEN
	professionalController, _, _ := controllerTest.NewProfessionalControllerTestWrapper(t)

	// WHEN
	result, errResult := professionalController.FetchProfessionals()

	// THEN
	assert.Nil(t, errResult)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result.Professionals))
}
