package local_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	adapterTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/adapter"
)

func TestLocalAdapter_CreatePostgresqlLocal_Success(t *testing.T) {
	// Given
	localAdapter, _, _ := adapterTest.NewLocalAdapterTestWrapper(t)

	// When
	result, err := localAdapter.CreatePostgresqlLocal(
		"Test Local",
		"Main Street",
		"123",
		"Downtown",
		"Lima",
		"Lima",
		"Near the park",
		50,
		"https://example.com/image.jpg",
		"test_user",
	)

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Local", result.LocalName)
	assert.Equal(t, "Main Street", result.StreetName)
	assert.Equal(t, "123", result.BuildingNumber)
	assert.Equal(t, "Downtown", result.District)
	assert.Equal(t, "Lima", result.Province)
	assert.Equal(t, "Lima", result.Region)
	assert.Equal(t, "Near the park", result.Reference)
	assert.Equal(t, 50, result.Capacity)
	assert.Equal(t, "https://example.com/image.jpg", result.ImageUrl)
	assert.NotEqual(t, "", result.Id)
}

func TestLocalAdapter_CreatePostgresqlLocal_EmptyUpdatedBy(t *testing.T) {
	// Given
	localAdapter, _, _ := adapterTest.NewLocalAdapterTestWrapper(t)

	// When
	result, err := localAdapter.CreatePostgresqlLocal(
		"Test Local",
		"Main Street",
		"123",
		"Downtown",
		"Lima",
		"Lima",
		"Near the park",
		50,
		"https://example.com/image.jpg",
		"",
	)

	// Then
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "Invalid updated by value", err.Message)
}

func TestLocalAdapter_CreatePostgresqlLocal_WithDatabaseUser(t *testing.T) {
	// Given
	localAdapter, _, db := adapterTest.NewLocalAdapterTestWrapper(t)

	// Using database factory for a real user
	testUser := factories.NewUserModel(db, factories.UserModelF{})

	// When
	result, err := localAdapter.CreatePostgresqlLocal(
		"Factory Local",
		"Factory Street",
		"456",
		"Factory District",
		"Factory Province",
		"Factory Region",
		"Factory Reference",
		100,
		"https://factory.com/image.jpg",
		testUser.Name,
	)

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Factory Local", result.LocalName)
	assert.Equal(t, "Factory Street", result.StreetName)
	assert.Equal(t, "456", result.BuildingNumber)
	assert.Equal(t, "Factory District", result.District)
	assert.Equal(t, "Factory Province", result.Province)
	assert.Equal(t, "Factory Region", result.Region)
	assert.Equal(t, "Factory Reference", result.Reference)
	assert.Equal(t, 100, result.Capacity)
	assert.Equal(t, "https://factory.com/image.jpg", result.ImageUrl)
	assert.NotEqual(t, "", result.Id)
}
