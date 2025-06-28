package community_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestGetCommunitySuccessfully(t *testing.T) {
	// GIVEN: An existing community
	controller, _, db := controllerTest.NewCommunityControllerTestWrapper(t)

	// Create a test community
	name := "Test Community"
	purpose := "Testing purposes"
	imageUrl := "https://example.com/community.jpg"

	testCommunity := factories.NewCommunityModel(db, factories.CommunityModelF{
		Name:     &name,
		Purpose:  &purpose,
		ImageUrl: &imageUrl,
	})

	// WHEN: GetCommunity is called
	result, err := controller.GetCommunity(testCommunity.Id)

	// THEN: Community is returned successfully
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testCommunity.Id, result.Id)
	assert.Equal(t, testCommunity.Name, result.Name)
	assert.Equal(t, testCommunity.Purpose, result.Purpose)
	assert.Equal(t, testCommunity.ImageUrl, result.ImageUrl)
}

func TestGetCommunityNotFound(t *testing.T) {
	// GIVEN: A non-existent community ID
	controller, _, _ := controllerTest.NewCommunityControllerTestWrapper(t)

	nonExistentId := uuid.New()

	// WHEN: GetCommunity is called with non-existent ID
	result, err := controller.GetCommunity(nonExistentId)

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Message, "not found")
}

func TestGetCommunityWithNilFields(t *testing.T) {
	// GIVEN: A community with some nil optional fields
	controller, _, db := controllerTest.NewCommunityControllerTestWrapper(t)

	// Create a community with minimal required fields
	name := "Minimal Community"
	testCommunity := factories.NewCommunityModel(db, factories.CommunityModelF{
		Name: &name,
		// Purpose and ImageUrl might be optional
	})

	// WHEN: GetCommunity is called
	result, err := controller.GetCommunity(testCommunity.Id)

	// THEN: Community is returned successfully
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testCommunity.Id, result.Id)
	assert.Equal(t, testCommunity.Name, result.Name)
	// Other fields should match what was created
	assert.Equal(t, testCommunity.Purpose, result.Purpose)
	assert.Equal(t, testCommunity.ImageUrl, result.ImageUrl)
}

func TestGetCommunityWithLongName(t *testing.T) {
	// GIVEN: A community with a long name
	controller, _, db := controllerTest.NewCommunityControllerTestWrapper(t)

	longName := "This is a very long community name that tests the system's ability to handle extended text fields and ensure proper storage and retrieval"
	purpose := "Testing long names"

	testCommunity := factories.NewCommunityModel(db, factories.CommunityModelF{
		Name:    &longName,
		Purpose: &purpose,
	})

	// WHEN: GetCommunity is called
	result, err := controller.GetCommunity(testCommunity.Id)

	// THEN: Community with long name is returned successfully
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testCommunity.Id, result.Id)
	assert.Equal(t, longName, result.Name)
	assert.Equal(t, purpose, result.Purpose)
}

func TestGetCommunityWithSpecialCharacters(t *testing.T) {
	// GIVEN: A community with special characters in name and purpose
	controller, _, db := controllerTest.NewCommunityControllerTestWrapper(t)

	specialName := "Café & Música 🎵 Community"
	specialPurpose := "Para compartir música, café ☕ y conversación en español"

	testCommunity := factories.NewCommunityModel(db, factories.CommunityModelF{
		Name:    &specialName,
		Purpose: &specialPurpose,
	})

	// WHEN: GetCommunity is called
	result, err := controller.GetCommunity(testCommunity.Id)

	// THEN: Community with special characters is returned successfully
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testCommunity.Id, result.Id)
	assert.Equal(t, specialName, result.Name)
	assert.Equal(t, specialPurpose, result.Purpose)
}
