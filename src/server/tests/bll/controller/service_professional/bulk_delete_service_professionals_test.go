package service_professional_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestBulkDeleteServiceProfessionalsSuccessfully(t *testing.T) {
	/*
		GIVEN: Service-professional associations exist
		WHEN:  BulkDeleteServiceProfessionals is called with valid data
		THEN:  All specified associations should be deleted successfully
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
			ProfessionalId: testProfessional2.Id,
			AuditFields:    model.AuditFields{UpdatedBy: "TEST_USER"},
		},
	}

	for _, association := range associations {
		err := db.Create(association).Error
		assert.NoError(t, err)
	}

	// Prepare delete data
	deleteData := schemas.BulkDeleteServiceProfessionalRequest{
		ServiceProfessionals: []*schemas.DeleteServiceProfessionalRequest{
			{
				ServiceId:      testService1.Id,
				ProfessionalId: testProfessional1.Id,
			},
			{
				ServiceId:      testService2.Id,
				ProfessionalId: testProfessional2.Id,
			},
		},
	}

	// WHEN
	err := serviceProfessionalController.BulkDeleteServiceProfessionals(deleteData)

	// THEN
	assert.Nil(t, err)

	// Verify deletion in database
	var remainingAssociations []model.ServiceProfessional
	dbErr := db.Where("service_id IN (?, ?) AND professional_id IN (?, ?)",
		testService1.Id, testService2.Id, testProfessional1.Id, testProfessional2.Id).Find(&remainingAssociations).Error
	assert.NoError(t, dbErr)
	assert.Len(t, remainingAssociations, 0) // Should be empty
}

func TestBulkDeleteServiceProfessionalsWithEmptyData(t *testing.T) {
	/*
		GIVEN: Empty delete data
		WHEN:  BulkDeleteServiceProfessionals is called with empty data
		THEN:  Operation should succeed without error
	*/
	// GIVEN
	serviceProfessionalController, _, _ := controllerTest.NewServiceProfessionalControllerTestWrapper(t)

	deleteData := schemas.BulkDeleteServiceProfessionalRequest{
		ServiceProfessionals: []*schemas.DeleteServiceProfessionalRequest{},
	}

	// WHEN
	err := serviceProfessionalController.BulkDeleteServiceProfessionals(deleteData)

	// THEN
	assert.Nil(t, err)
}

func TestBulkDeleteServiceProfessionalsWithNonExistentAssociations(t *testing.T) {
	/*
		GIVEN: Non-existent service-professional associations
		WHEN:  BulkDeleteServiceProfessionals is called with non-existent data
		THEN:  Operation should handle gracefully (implementation dependent)
	*/
	// GIVEN
	serviceProfessionalController, _, db := controllerTest.NewServiceProfessionalControllerTestWrapper(t)

	// Create test services and professionals but don't create associations
	testService := factories.NewServiceModel(db)
	testProfessional := factories.NewProfessionalModel(db)

	deleteData := schemas.BulkDeleteServiceProfessionalRequest{
		ServiceProfessionals: []*schemas.DeleteServiceProfessionalRequest{
			{
				ServiceId:      testService.Id,
				ProfessionalId: testProfessional.Id,
			},
		},
	}

	// WHEN
	err := serviceProfessionalController.BulkDeleteServiceProfessionals(deleteData)

	// THEN
	// The behavior depends on implementation - it might succeed or fail
	// For bulk operations, it's common to succeed even if some items don't exist
	// We'll assume it succeeds for non-existent items
	assert.Nil(t, err)
}
