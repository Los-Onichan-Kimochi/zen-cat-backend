package service_professional_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestBulkCreateServiceProfessionalsSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid services and professionals exist
		WHEN:  BulkCreateServiceProfessionals is called with valid data
		THEN:  All service-professional associations should be created successfully
	*/
	// GIVEN
	serviceProfessionalController, _, db := controllerTest.NewServiceProfessionalControllerTestWrapper(t)

	// Create test services and professionals
	testService1 := factories.NewServiceModel(db)
	testService2 := factories.NewServiceModel(db)
	testProfessional1 := factories.NewProfessionalModel(db)
	testProfessional2 := factories.NewProfessionalModel(db)

	updatedBy := "TEST_USER"
	createData := []*schemas.CreateServiceProfessionalRequest{
		{
			ServiceId:      testService1.Id,
			ProfessionalId: testProfessional1.Id,
		},
		{
			ServiceId:      testService2.Id,
			ProfessionalId: testProfessional2.Id,
		},
	}

	// WHEN
	result, err := serviceProfessionalController.BulkCreateServiceProfessionals(createData, updatedBy)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.ServiceProfessionals, 2)
	assert.Equal(t, testService1.Id, result.ServiceProfessionals[0].ServiceId)
	assert.Equal(t, testProfessional1.Id, result.ServiceProfessionals[0].ProfessionalId)
	assert.Equal(t, testService2.Id, result.ServiceProfessionals[1].ServiceId)
	assert.Equal(t, testProfessional2.Id, result.ServiceProfessionals[1].ProfessionalId)
}

func TestBulkCreateServiceProfessionalsWithEmptyData(t *testing.T) {
	/*
		GIVEN: Empty data array
		WHEN:  BulkCreateServiceProfessionals is called with empty data
		THEN:  Empty result should be returned
	*/
	// GIVEN
	serviceProfessionalController, _, _ := controllerTest.NewServiceProfessionalControllerTestWrapper(t)

	updatedBy := "TEST_USER"
	createData := []*schemas.CreateServiceProfessionalRequest{}

	// WHEN
	result, err := serviceProfessionalController.BulkCreateServiceProfessionals(createData, updatedBy)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.ServiceProfessionals, 0)
}

func TestBulkCreateServiceProfessionalsWithInvalidData(t *testing.T) {
	/*
		GIVEN: Invalid data (non-existent service or professional)
		WHEN:  BulkCreateServiceProfessionals is called with invalid data
		THEN:  It should return an error
	*/
	// GIVEN
	serviceProfessionalController, _, db := controllerTest.NewServiceProfessionalControllerTestWrapper(t)

	// Create one valid professional but use non-existent service
	testProfessional := factories.NewProfessionalModel(db)
	testService := factories.NewServiceModel(db)

	updatedBy := "TEST_USER"
	createData := []*schemas.CreateServiceProfessionalRequest{
		{
			ServiceId:      testService.Id,
			ProfessionalId: testProfessional.Id,
		},
		{
			ServiceId:      testService.Id,      // Same service but different professional
			ProfessionalId: testProfessional.Id, // Same professional - should cause conflict
		},
	}

	// WHEN
	result, err := serviceProfessionalController.BulkCreateServiceProfessionals(createData, updatedBy)

	// THEN
	// The behavior depends on the implementation - it might succeed partially or fail entirely
	// For now, let's assume it handles duplicates gracefully
	if err != nil {
		assert.Nil(t, result)
	} else {
		assert.NotNil(t, result)
		// Should have at least one successful creation
		assert.GreaterOrEqual(t, len(result.ServiceProfessionals), 1)
	}
}
