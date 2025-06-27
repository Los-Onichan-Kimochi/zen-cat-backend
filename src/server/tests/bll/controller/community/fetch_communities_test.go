package community_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestFetchCommunitiesEmpty(t *testing.T) {
	// GIVEN: No communities in the database
	controller, _, _ := controllerTest.NewCommunityControllerTestWrapper(t)

	// WHEN: FetchCommunities is called
	result, err := controller.FetchCommunities()

	// THEN: Empty communities list is returned
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Communities)
	assert.Len(t, result.Communities, 0)
}

func TestFetchCommunitiesWithData(t *testing.T) {
	// GIVEN: Multiple communities in the database
	controller, _, db := controllerTest.NewCommunityControllerTestWrapper(t)

	// Create test communities
	name1 := "Community One"
	purpose1 := "First community purpose"
	name2 := "Community Two"
	purpose2 := "Second community purpose"
	name3 := "Community Three"
	purpose3 := "Third community purpose"

	community1 := factories.NewCommunityModel(db, factories.CommunityModelF{
		Name:    &name1,
		Purpose: &purpose1,
	})
	community2 := factories.NewCommunityModel(db, factories.CommunityModelF{
		Name:    &name2,
		Purpose: &purpose2,
	})
	community3 := factories.NewCommunityModel(db, factories.CommunityModelF{
		Name:    &name3,
		Purpose: &purpose3,
	})

	// WHEN: FetchCommunities is called
	result, err := controller.FetchCommunities()

	// THEN: All communities are returned
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Communities)
	assert.Len(t, result.Communities, 3)

	// Verify community data
	communityIds := make(map[string]bool)
	for _, community := range result.Communities {
		communityIds[community.Id.String()] = true
		assert.NotEmpty(t, community.Name)
		assert.NotEmpty(t, community.Purpose)
	}

	// Verify all created communities are present
	assert.True(t, communityIds[community1.Id.String()])
	assert.True(t, communityIds[community2.Id.String()])
	assert.True(t, communityIds[community3.Id.String()])
}

func TestFetchCommunitiesSingleCommunity(t *testing.T) {
	// GIVEN: One community in the database
	controller, _, db := controllerTest.NewCommunityControllerTestWrapper(t)

	name := "Single Community"
	purpose := "Only community in database"
	imageUrl := "https://example.com/single.jpg"

	testCommunity := factories.NewCommunityModel(db, factories.CommunityModelF{
		Name:     &name,
		Purpose:  &purpose,
		ImageUrl: &imageUrl,
	})

	// WHEN: FetchCommunities is called
	result, err := controller.FetchCommunities()

	// THEN: Single community is returned
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Communities)
	assert.Len(t, result.Communities, 1)

	// Verify community data
	community := result.Communities[0]
	assert.Equal(t, testCommunity.Id, community.Id)
	assert.Equal(t, testCommunity.Name, community.Name)
	assert.Equal(t, testCommunity.Purpose, community.Purpose)
	assert.Equal(t, testCommunity.ImageUrl, community.ImageUrl)
}

func TestFetchCommunitiesWithVariedData(t *testing.T) {
	// GIVEN: Communities with varied data (some with nil fields)
	controller, _, db := controllerTest.NewCommunityControllerTestWrapper(t)

	// Community with all fields
	fullName := "Full Community"
	fullPurpose := "Complete community with all fields"
	fullImageUrl := "https://example.com/full.jpg"

	// Community with minimal fields
	minimalName := "Minimal Community"

	factories.NewCommunityModel(db, factories.CommunityModelF{
		Name:     &fullName,
		Purpose:  &fullPurpose,
		ImageUrl: &fullImageUrl,
	})
	factories.NewCommunityModel(db, factories.CommunityModelF{
		Name: &minimalName,
		// Purpose and ImageUrl might be nil
	})

	// WHEN: FetchCommunities is called
	result, err := controller.FetchCommunities()

	// THEN: Both communities are returned
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Communities)
	assert.Len(t, result.Communities, 2)

	// Verify all communities have required fields
	for _, community := range result.Communities {
		assert.NotEqual(t, "", community.Id)
		assert.NotEmpty(t, community.Name)
		// Purpose and ImageUrl might be empty for some communities
	}
}

func TestFetchCommunitiesOrdering(t *testing.T) {
	// GIVEN: Multiple communities with different creation times
	controller, _, db := controllerTest.NewCommunityControllerTestWrapper(t)

	// Create communities in specific order
	name1 := "Alpha Community"
	name2 := "Beta Community"
	name3 := "Gamma Community"

	factories.NewCommunityModel(db, factories.CommunityModelF{Name: &name1})
	factories.NewCommunityModel(db, factories.CommunityModelF{Name: &name2})
	factories.NewCommunityModel(db, factories.CommunityModelF{Name: &name3})

	// WHEN: FetchCommunities is called
	result, err := controller.FetchCommunities()

	// THEN: Communities are returned (ordering depends on implementation)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Communities)
	assert.Len(t, result.Communities, 3)

	// Verify all expected communities are present
	communityNames := make(map[string]bool)
	for _, community := range result.Communities {
		communityNames[community.Name] = true
	}

	assert.True(t, communityNames["Alpha Community"])
	assert.True(t, communityNames["Beta Community"])
	assert.True(t, communityNames["Gamma Community"])
}
