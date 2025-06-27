package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	adapterTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/adapter"
)

func TestServiceAdapter_CreatePostgresqlService_Success(t *testing.T) {
	/*
		GIVEN: Valid service data
		WHEN:  CreatePostgresqlService is called
		THEN:  A new service is created and returned
	*/
	// GIVEN
	serviceAdapter, _, _ := adapterTest.NewServiceAdapterTestWrapper(t)

	// WHEN
	result, err := serviceAdapter.CreatePostgresqlService(
		"Therapy Session",
		"Individual therapy session",
		"https://example.com/therapy.jpg",
		false,
		"test_user",
	)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Therapy Session", result.Name)
	assert.Equal(t, "Individual therapy session", result.Description)
	assert.Equal(t, "https://example.com/therapy.jpg", result.ImageUrl)
	assert.Equal(t, false, result.IsVirtual)
	assert.NotEqual(t, "", result.Id)
}

func TestServiceAdapter_CreatePostgresqlService_VirtualService(t *testing.T) {
	/*
		GIVEN: Valid virtual service data
		WHEN:  CreatePostgresqlService is called
		THEN:  A new virtual service is created and returned
	*/
	// GIVEN
	serviceAdapter, _, _ := adapterTest.NewServiceAdapterTestWrapper(t)

	// WHEN
	result, err := serviceAdapter.CreatePostgresqlService(
		"Online Consultation",
		"Virtual consultation session",
		"https://example.com/virtual.jpg",
		true,
		"test_user",
	)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Online Consultation", result.Name)
	assert.Equal(t, "Virtual consultation session", result.Description)
	assert.Equal(t, "https://example.com/virtual.jpg", result.ImageUrl)
	assert.Equal(t, true, result.IsVirtual)
	assert.NotEqual(t, "", result.Id)
}

func TestServiceAdapter_CreatePostgresqlService_EmptyUpdatedBy(t *testing.T) {
	/*
		GIVEN: Valid service data but empty updatedBy
		WHEN:  CreatePostgresqlService is called
		THEN:  An error is returned
	*/
	// GIVEN
	serviceAdapter, _, _ := adapterTest.NewServiceAdapterTestWrapper(t)

	// WHEN
	result, err := serviceAdapter.CreatePostgresqlService(
		"Therapy Session",
		"Individual therapy session",
		"https://example.com/therapy.jpg",
		false,
		"",
	)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "Invalid updated by value", err.Message)
}

func TestServiceAdapter_CreatePostgresqlService_WithDatabaseUser(t *testing.T) {
	/*
		GIVEN: Valid service data with database user
		WHEN:  CreatePostgresqlService is called
		THEN:  A new service is created and returned
	*/
	// GIVEN
	serviceAdapter, _, db := adapterTest.NewServiceAdapterTestWrapper(t)

	// Using database factory for a real user
	testUser := factories.NewUserModel(db, factories.UserModelF{})

	// WHEN
	result, err := serviceAdapter.CreatePostgresqlService(
		"Group Therapy",
		"Group therapy session for multiple clients",
		"https://example.com/group.jpg",
		false,
		testUser.Name,
	)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Group Therapy", result.Name)
	assert.Equal(t, "Group therapy session for multiple clients", result.Description)
	assert.Equal(t, "https://example.com/group.jpg", result.ImageUrl)
	assert.Equal(t, false, result.IsVirtual)
	assert.NotEqual(t, "", result.Id)
}
