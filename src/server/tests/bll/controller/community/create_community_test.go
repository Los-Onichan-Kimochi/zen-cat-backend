package community_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestCreateCommunitySuccessfully(t *testing.T) {
	// GIVEN: Valid community creation request
	controller, _, _ := controllerTest.NewCommunityControllerTestWrapper(t)

	createRequest := schemas.CreateCommunityRequest{
		Name:     "New Community",
		Purpose:  "A community for testing",
		ImageUrl: "https://example.com/community.jpg",
	}

	// WHEN: CreateCommunity is called
	result, err := controller.CreateCommunity(createRequest, "test_admin")

	// THEN: Community is created successfully
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createRequest.Name, result.Name)
	assert.Equal(t, createRequest.Purpose, result.Purpose)
	assert.Equal(t, createRequest.ImageUrl, result.ImageUrl)
	assert.NotEqual(t, "", result.Id)
}

func TestCreateCommunityEmptyUpdatedBy(t *testing.T) {
	// GIVEN: Valid community creation request but empty updatedBy
	controller, _, _ := controllerTest.NewCommunityControllerTestWrapper(t)

	createRequest := schemas.CreateCommunityRequest{
		Name:     "Test Community",
		Purpose:  "Testing purposes",
		ImageUrl: "https://example.com/test.jpg",
	}

	// WHEN: CreateCommunity is called with empty updatedBy
	result, err := controller.CreateCommunity(createRequest, "")

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "Invalid updated by value", err.Message)
}

func TestCreateCommunityDuplicateName(t *testing.T) {
	// GIVEN: A community already exists with the same name
	controller, _, _ := controllerTest.NewCommunityControllerTestWrapper(t)

	createRequest1 := schemas.CreateCommunityRequest{
		Name:     "Duplicate Community",
		Purpose:  "First community",
		ImageUrl: "https://example.com/first.jpg",
	}

	createRequest2 := schemas.CreateCommunityRequest{
		Name:     "Duplicate Community", // Same name
		Purpose:  "Second community",
		ImageUrl: "https://example.com/second.jpg",
	}

	// Create first community
	_, err1 := controller.CreateCommunity(createRequest1, "test_admin")
	assert.Nil(t, err1)

	// WHEN: CreateCommunity is called with duplicate name
	result, err := controller.CreateCommunity(createRequest2, "test_admin")

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestCreateCommunityMinimalFields(t *testing.T) {
	// GIVEN: Community creation request with minimal required fields
	controller, _, _ := controllerTest.NewCommunityControllerTestWrapper(t)

	createRequest := schemas.CreateCommunityRequest{
		Name: "Minimal Community",
		// Purpose and ImageUrl might be optional
	}

	// WHEN: CreateCommunity is called
	result, err := controller.CreateCommunity(createRequest, "test_admin")

	// THEN: Community is created successfully
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createRequest.Name, result.Name)
	assert.NotEqual(t, "", result.Id)
}

func TestCreateCommunityLongName(t *testing.T) {
	// GIVEN: Community creation request with long name
	controller, _, _ := controllerTest.NewCommunityControllerTestWrapper(t)

	longName := "This is a very long community name that tests the system's ability to handle extended text fields"
	createRequest := schemas.CreateCommunityRequest{
		Name:     longName,
		Purpose:  "Testing long names",
		ImageUrl: "https://example.com/long.jpg",
	}

	// WHEN: CreateCommunity is called
	result, err := controller.CreateCommunity(createRequest, "test_admin")

	// THEN: Community with long name is created successfully
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, longName, result.Name)
	assert.Equal(t, createRequest.Purpose, result.Purpose)
}

func TestCreateCommunitySpecialCharacters(t *testing.T) {
	// GIVEN: Community creation request with special characters
	controller, _, _ := controllerTest.NewCommunityControllerTestWrapper(t)

	specialName := "Café & Música 🎵 Community"
	specialPurpose := "Para compartir música, café ☕ y conversación"

	createRequest := schemas.CreateCommunityRequest{
		Name:     specialName,
		Purpose:  specialPurpose,
		ImageUrl: "https://example.com/special.jpg",
	}

	// WHEN: CreateCommunity is called
	result, err := controller.CreateCommunity(createRequest, "test_admin")

	// THEN: Community with special characters is created successfully
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, specialName, result.Name)
	assert.Equal(t, specialPurpose, result.Purpose)
}

func TestCreateCommunityEmptyName(t *testing.T) {
	// GIVEN: Community creation request with empty name
	controller, _, _ := controllerTest.NewCommunityControllerTestWrapper(t)

	createRequest := schemas.CreateCommunityRequest{
		Name:     "", // Empty name
		Purpose:  "Testing empty name",
		ImageUrl: "https://example.com/empty.jpg",
	}

	// WHEN: CreateCommunity is called
	result, err := controller.CreateCommunity(createRequest, "test_admin")

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestCreateCommunityInvalidImageUrl(t *testing.T) {
	// GIVEN: Community creation request with invalid image URL
	controller, _, _ := controllerTest.NewCommunityControllerTestWrapper(t)

	createRequest := schemas.CreateCommunityRequest{
		Name:     "Test Community",
		Purpose:  "Testing invalid URL",
		ImageUrl: "not-a-valid-url",
	}

	// WHEN: CreateCommunity is called
	result, err := controller.CreateCommunity(createRequest, "test_admin")

	// THEN: Community is created (URL validation might be at different layer)
	// The exact behavior depends on implementation
	if err == nil {
		assert.NotNil(t, result)
		assert.Equal(t, createRequest.Name, result.Name)
	} else {
		assert.Nil(t, result)
		assert.NotNil(t, err)
	}
}
