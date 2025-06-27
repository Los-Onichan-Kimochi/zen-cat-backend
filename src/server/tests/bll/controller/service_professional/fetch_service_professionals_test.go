package service_professional_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestFetchServiceProfessionalsWithServiceFilter(t *testing.T) {
	/*
		GIVEN: Multiple service-professional associations exist
		WHEN:  FetchServiceProfessionals is called with service ID filter
		THEN:  Only associations for that service should be returned
	*/
	// GIVEN
	serviceProfessionalController, _, db := controllerTest.NewServiceProfessionalControllerTestWrapper(t)

	// Create test services and professionals
	testService1 := factories.NewServiceModel(db)
	testService2 := factories.NewServiceModel(db)
	testProfessional1 := factories.NewProfessionalModel(db)
	testProfessional2 := factories.NewProfessionalModel(db)

	// Create service-professional associations
	associations := []*model.ServiceProfessional{
		{
			ServiceId:      testService1.Id,
			ProfessionalId: testProfessional1.Id,
			AuditFields:    model.AuditFields{UpdatedBy: "TEST_USER"},
		},
		{
			ServiceId:      testService1.Id,
			ProfessionalId: testProfessional2.Id,
			AuditFields:    model.AuditFields{UpdatedBy: "TEST_USER"},
		},
		{
			ServiceId:      testService2.Id,
			ProfessionalId: testProfessional1.Id,
			AuditFields:    model.AuditFields{UpdatedBy: "TEST_USER"},
		},
	}

	for _, association := range associations {
		err := db.Create(association).Error
		assert.NoError(t, err)
	}

	// WHEN
	result, err := serviceProfessionalController.FetchServiceProfessionals(
		testService1.Id.String(),
		"",
	)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.ServiceProfessionals, 2) // Only associations for testService1
	for _, sp := range result.ServiceProfessionals {
		assert.Equal(t, testService1.Id, sp.ServiceId)
	}
}

func TestFetchServiceProfessionalsWithProfessionalFilter(t *testing.T) {
	/*
		GIVEN: Multiple service-professional associations exist
		WHEN:  FetchServiceProfessionals is called with professional ID filter
		THEN:  Only associations for that professional should be returned
	*/
	// GIVEN
	serviceProfessionalController, _, db := controllerTest.NewServiceProfessionalControllerTestWrapper(t)

	// Create test services and professionals
	testService1 := factories.NewServiceModel(db)
	testService2 := factories.NewServiceModel(db)
	testProfessional1 := factories.NewProfessionalModel(db)
	testProfessional2 := factories.NewProfessionalModel(db)

	// Create service-professional associations
	associations := []*model.ServiceProfessional{
		{
			ServiceId:      testService1.Id,
			ProfessionalId: testProfessional1.Id,
			AuditFields:    model.AuditFields{UpdatedBy: "TEST_USER"},
		},
		{
			ServiceId:      testService2.Id,
			ProfessionalId: testProfessional1.Id,
			AuditFields:    model.AuditFields{UpdatedBy: "TEST_USER"},
		},
		{
			ServiceId:      testService1.Id,
			ProfessionalId: testProfessional2.Id,
			AuditFields:    model.AuditFields{UpdatedBy: "TEST_USER"},
		},
	}

	for _, association := range associations {
		err := db.Create(association).Error
		assert.NoError(t, err)
	}

	// WHEN
	result, err := serviceProfessionalController.FetchServiceProfessionals(
		"",
		testProfessional1.Id.String(),
	)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.ServiceProfessionals, 2) // Only associations for testProfessional1
	for _, sp := range result.ServiceProfessionals {
		assert.Equal(t, testProfessional1.Id, sp.ProfessionalId)
	}
}

func TestFetchServiceProfessionalsWithInvalidServiceId(t *testing.T) {
	/*
		GIVEN: Invalid service ID format
		WHEN:  FetchServiceProfessionals is called with invalid service ID
		THEN:  It should return invalid service ID error
	*/
	// GIVEN
	serviceProfessionalController, _, _ := controllerTest.NewServiceProfessionalControllerTestWrapper(t)

	invalidServiceId := "invalid-uuid"

	// WHEN
	result, err := serviceProfessionalController.FetchServiceProfessionals(
		invalidServiceId,
		"",
	)

	// THEN
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, errors.UnprocessableEntityError.InvalidServiceId.Code, err.Code)
}

func TestFetchServiceProfessionalsWithNonExistentService(t *testing.T) {
	/*
		GIVEN: Service does not exist
		WHEN:  FetchServiceProfessionals is called with non-existent service ID
		THEN:  It should return service not found error
	*/
	// GIVEN
	serviceProfessionalController, _, _ := controllerTest.NewServiceProfessionalControllerTestWrapper(t)

	nonExistentServiceId := uuid.New()

	// WHEN
	result, err := serviceProfessionalController.FetchServiceProfessionals(
		nonExistentServiceId.String(),
		"",
	)

	// THEN
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, errors.ObjectNotFoundError.ServiceNotFound.Code, err.Code)
}
