package community_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	adapterTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/adapter"
)

func TestCreateCommunitySuccessfully(t *testing.T) {
	/*
		GIVEN: Valid community data
		WHEN:  CreatePostgresqlCommunity is called
		THEN:  A new community is created and returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewCommunityAdapterTestWrapper(t)

	name := "Test Community"
	purpose := "Testing purposes"
	imageUrl := "https://example.com/image.jpg"
	updatedBy := "test-user"

	// WHEN
	community, err := adapter.CreatePostgresqlCommunity(name, purpose, imageUrl, updatedBy)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, community)
	assert.NotEmpty(t, community.Id)
	assert.Equal(t, name, community.Name)
	assert.Equal(t, purpose, community.Purpose)
	assert.Equal(t, imageUrl, community.ImageUrl)
	assert.Equal(t, int(0), community.NumberSubscriptions)
}

func TestCreateCommunityWithDuplicateName(t *testing.T) {
	/*
		GIVEN: A community with the same name already exists
		WHEN:  CreatePostgresqlCommunity is called with duplicate name
		THEN:  An error is returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewCommunityAdapterTestWrapper(t)

	duplicateName := "Duplicate Community"
	existingCommunity := factories.NewCommunityModel(db, factories.CommunityModelF{
		Name: &duplicateName,
	})

	name := existingCommunity.Name
	purpose := "Testing purposes"
	imageUrl := "https://example.com/image.jpg"
	updatedBy := "test-user"

	// WHEN
	community, err := adapter.CreatePostgresqlCommunity(name, purpose, imageUrl, updatedBy)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, community)
	assert.Equal(t, errors.BadRequestError.DuplicateCommunityName, *err)
}

func TestCreateCommunityWithEmptyUpdatedBy(t *testing.T) {
	/*
		GIVEN: Valid community data but empty updatedBy
		WHEN:  CreatePostgresqlCommunity is called
		THEN:  An error is returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewCommunityAdapterTestWrapper(t)

	name := "Test Community"
	purpose := "Testing purposes"
	imageUrl := "https://example.com/image.jpg"
	updatedBy := ""

	// WHEN
	community, err := adapter.CreatePostgresqlCommunity(name, purpose, imageUrl, updatedBy)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, community)
	assert.Equal(t, errors.BadRequestError.InvalidUpdatedByValue, *err)
}

func TestCreateCommunityWithEmptyImageUrl(t *testing.T) {
	/*
		GIVEN: Valid community data with empty image URL
		WHEN:  CreatePostgresqlCommunity is called
		THEN:  A new community is created with empty image URL
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewCommunityAdapterTestWrapper(t)

	name := "Test Community"
	purpose := "Testing purposes"
	imageUrl := ""
	updatedBy := "test-user"

	// WHEN
	community, err := adapter.CreatePostgresqlCommunity(name, purpose, imageUrl, updatedBy)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, community)
	assert.Equal(t, name, community.Name)
	assert.Equal(t, purpose, community.Purpose)
	assert.Equal(t, imageUrl, community.ImageUrl)
}
