package service_professional_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	adapterTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/adapter"
)

func TestBulkCreateServiceProfessionalsSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid service-professional data for bulk creation
		WHEN:  BulkCreatePostgresqlServiceProfessionals is called
		THEN:  All associations are created and returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewServiceProfessionalAdapterTestWrapper(t)

	service1 := factories.NewServiceModel(db, factories.ServiceModelF{})
	service2 := factories.NewServiceModel(db, factories.ServiceModelF{})
	professional1 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	professional2 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	serviceProfessionalsData := []*schemas.CreateServiceProfessionalRequest{
		{
			ServiceId:      service1.Id,
			ProfessionalId: professional1.Id,
		},
		{
			ServiceId:      service2.Id,
			ProfessionalId: professional2.Id,
		},
	}
	updatedBy := "test-admin"

	// WHEN
	serviceProfessionals, err := adapter.BulkCreatePostgresqlServiceProfessionals(serviceProfessionalsData, updatedBy)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, serviceProfessionals)
	assert.Equal(t, 2, len(serviceProfessionals))

	// Verify first association
	assert.NotEmpty(t, serviceProfessionals[0].Id)
	assert.Equal(t, service1.Id, serviceProfessionals[0].ServiceId)
	assert.Equal(t, professional1.Id, serviceProfessionals[0].ProfessionalId)

	// Verify second association
	assert.NotEmpty(t, serviceProfessionals[1].Id)
	assert.Equal(t, service2.Id, serviceProfessionals[1].ServiceId)
	assert.Equal(t, professional2.Id, serviceProfessionals[1].ProfessionalId)
}

func TestBulkCreateServiceProfessionalsWithEmptyUpdatedBy(t *testing.T) {
	/*
		GIVEN: Valid service-professional data but empty updatedBy
		WHEN:  BulkCreatePostgresqlServiceProfessionals is called
		THEN:  An error is returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewServiceProfessionalAdapterTestWrapper(t)

	service := factories.NewServiceModel(db, factories.ServiceModelF{})
	professional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	serviceProfessionalsData := []*schemas.CreateServiceProfessionalRequest{
		{
			ServiceId:      service.Id,
			ProfessionalId: professional.Id,
		},
	}
	updatedBy := ""

	// WHEN
	serviceProfessionals, err := adapter.BulkCreatePostgresqlServiceProfessionals(serviceProfessionalsData, updatedBy)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, serviceProfessionals)
	assert.Equal(t, errors.BadRequestError.InvalidUpdatedByValue, *err)
}

func TestFetchServiceProfessionalsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple service-professional associations exist
		WHEN:  FetchPostgresqlServiceProfessionals is called
		THEN:  All matching associations are returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewServiceProfessionalAdapterTestWrapper(t)

	service1 := factories.NewServiceModel(db, factories.ServiceModelF{})
	service2 := factories.NewServiceModel(db, factories.ServiceModelF{})
	professional1 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	professional2 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	// Create associations
	serviceProfessional1, err1 := adapter.CreatePostgresqlServiceProfessional(service1.Id, professional1.Id, "test-admin")
	assert.Nil(t, err1)
	serviceProfessional2, err2 := adapter.CreatePostgresqlServiceProfessional(service2.Id, professional2.Id, "test-admin")
	assert.Nil(t, err2)

	// WHEN - Fetch all associations
	serviceProfessionals, err := adapter.FetchPostgresqlServiceProfessionals(nil, nil)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, serviceProfessionals)
	assert.GreaterOrEqual(t, len(serviceProfessionals), 2)

	// Find our created associations
	foundServiceProfessional1 := false
	foundServiceProfessional2 := false
	for _, serviceProfessional := range serviceProfessionals {
		if serviceProfessional.Id == serviceProfessional1.Id {
			foundServiceProfessional1 = true
		}
		if serviceProfessional.Id == serviceProfessional2.Id {
			foundServiceProfessional2 = true
		}
	}
	assert.True(t, foundServiceProfessional1)
	assert.True(t, foundServiceProfessional2)
}

func TestFetchServiceProfessionalsWithServiceFilter(t *testing.T) {
	/*
		GIVEN: Associations exist for different services
		WHEN:  FetchPostgresqlServiceProfessionals is called with service filter
		THEN:  Only associations for that service are returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewServiceProfessionalAdapterTestWrapper(t)

	service1 := factories.NewServiceModel(db, factories.ServiceModelF{})
	service2 := factories.NewServiceModel(db, factories.ServiceModelF{})
	professional1 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	professional2 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	// Create associations
	serviceProfessional1, err1 := adapter.CreatePostgresqlServiceProfessional(service1.Id, professional1.Id, "test-admin")
	assert.Nil(t, err1)
	_, err2 := adapter.CreatePostgresqlServiceProfessional(service2.Id, professional2.Id, "test-admin")
	assert.Nil(t, err2)

	// WHEN - Filter by service1
	serviceProfessionals, err := adapter.FetchPostgresqlServiceProfessionals(&service1.Id, nil)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, serviceProfessionals)
	assert.GreaterOrEqual(t, len(serviceProfessionals), 1)

	// Verify all returned associations are for service1
	for _, serviceProfessional := range serviceProfessionals {
		assert.Equal(t, service1.Id, serviceProfessional.ServiceId)
	}

	// Verify our association is in the results
	foundServiceProfessional1 := false
	for _, serviceProfessional := range serviceProfessionals {
		if serviceProfessional.Id == serviceProfessional1.Id {
			foundServiceProfessional1 = true
		}
	}
	assert.True(t, foundServiceProfessional1)
}

func TestFetchServiceProfessionalsWithProfessionalFilter(t *testing.T) {
	/*
		GIVEN: Associations exist for different professionals
		WHEN:  FetchPostgresqlServiceProfessionals is called with professional filter
		THEN:  Only associations for that professional are returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewServiceProfessionalAdapterTestWrapper(t)

	service1 := factories.NewServiceModel(db, factories.ServiceModelF{})
	service2 := factories.NewServiceModel(db, factories.ServiceModelF{})
	professional1 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	professional2 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	// Create associations
	serviceProfessional1, err1 := adapter.CreatePostgresqlServiceProfessional(service1.Id, professional1.Id, "test-admin")
	assert.Nil(t, err1)
	_, err2 := adapter.CreatePostgresqlServiceProfessional(service2.Id, professional2.Id, "test-admin")
	assert.Nil(t, err2)

	// WHEN - Filter by professional1
	serviceProfessionals, err := adapter.FetchPostgresqlServiceProfessionals(nil, &professional1.Id)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, serviceProfessionals)
	assert.GreaterOrEqual(t, len(serviceProfessionals), 1)

	// Verify all returned associations are for professional1
	for _, serviceProfessional := range serviceProfessionals {
		assert.Equal(t, professional1.Id, serviceProfessional.ProfessionalId)
	}

	// Verify our association is in the results
	foundServiceProfessional1 := false
	for _, serviceProfessional := range serviceProfessionals {
		if serviceProfessional.Id == serviceProfessional1.Id {
			foundServiceProfessional1 = true
		}
	}
	assert.True(t, foundServiceProfessional1)
}

func TestBulkDeleteServiceProfessionalsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple service-professional associations exist
		WHEN:  BulkDeletePostgresqlServiceProfessionals is called
		THEN:  All specified associations are deleted
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewServiceProfessionalAdapterTestWrapper(t)

	service1 := factories.NewServiceModel(db, factories.ServiceModelF{})
	service2 := factories.NewServiceModel(db, factories.ServiceModelF{})
	service3 := factories.NewServiceModel(db, factories.ServiceModelF{})
	professional1 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	professional2 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})
	professional3 := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	// Create associations
	_, err1 := adapter.CreatePostgresqlServiceProfessional(service1.Id, professional1.Id, "test-admin")
	assert.Nil(t, err1)
	_, err2 := adapter.CreatePostgresqlServiceProfessional(service2.Id, professional2.Id, "test-admin")
	assert.Nil(t, err2)
	_, err3 := adapter.CreatePostgresqlServiceProfessional(service3.Id, professional3.Id, "test-admin")
	assert.Nil(t, err3)

	deleteRequests := []*schemas.DeleteServiceProfessionalRequest{
		{
			ServiceId:      service1.Id,
			ProfessionalId: professional1.Id,
		},
		{
			ServiceId:      service2.Id,
			ProfessionalId: professional2.Id,
		},
	}

	// WHEN
	err := adapter.BulkDeletePostgresqlServiceProfessionals(deleteRequests)

	// THEN
	assert.Nil(t, err)

	// Verify deleted associations cannot be found
	_, getErr1 := adapter.GetPostgresqlServiceProfessional(service1.Id, professional1.Id)
	assert.NotNil(t, getErr1)
	assert.Equal(t, errors.ObjectNotFoundError.ServiceProfessionalNotFound, *getErr1)

	_, getErr2 := adapter.GetPostgresqlServiceProfessional(service2.Id, professional2.Id)
	assert.NotNil(t, getErr2)
	assert.Equal(t, errors.ObjectNotFoundError.ServiceProfessionalNotFound, *getErr2)

	// Verify non-deleted association still exists
	serviceProfessional3Result, getErr3 := adapter.GetPostgresqlServiceProfessional(service3.Id, professional3.Id)
	assert.Nil(t, getErr3)
	assert.NotNil(t, serviceProfessional3Result)
	assert.Equal(t, service3.Id, serviceProfessional3Result.ServiceId)
	assert.Equal(t, professional3.Id, serviceProfessional3Result.ProfessionalId)
}

func TestBulkDeleteServiceProfessionalsWithEmptyList(t *testing.T) {
	/*
		GIVEN: An empty list of delete requests
		WHEN:  BulkDeletePostgresqlServiceProfessionals is called
		THEN:  No error occurs
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewServiceProfessionalAdapterTestWrapper(t)

	emptyRequests := []*schemas.DeleteServiceProfessionalRequest{}

	// WHEN
	err := adapter.BulkDeletePostgresqlServiceProfessionals(emptyRequests)

	// THEN
	assert.Nil(t, err)
}

func TestBulkDeleteServiceProfessionalsWithInvalidIds(t *testing.T) {
	/*
		GIVEN: Invalid service or professional IDs in delete requests
		WHEN:  BulkDeletePostgresqlServiceProfessionals is called
		THEN:  An error is returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewServiceProfessionalAdapterTestWrapper(t)

	invalidRequests := []*schemas.DeleteServiceProfessionalRequest{
		{
			ServiceId:      uuid.Nil, // Invalid ID
			ProfessionalId: uuid.New(),
		},
	}

	// WHEN
	err := adapter.BulkDeletePostgresqlServiceProfessionals(invalidRequests)

	// THEN
	assert.NotNil(t, err)
	assert.Equal(t, errors.UnprocessableEntityError.InvalidServiceProfessionalId, *err)
}
