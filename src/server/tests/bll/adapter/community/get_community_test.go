package community_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	adapterTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/adapter"
)

func TestGetCommunitySuccessfully(t *testing.T) {
	/*
		GIVEN: A community exists in the database
		WHEN:  GetPostgresqlCommunity is called with the community ID
		THEN:  The community is returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewCommunityAdapterTestWrapper(t)

	existingCommunity := factories.NewCommunityModel(db, factories.CommunityModelF{})

	// WHEN
	community, err := adapter.GetPostgresqlCommunity(existingCommunity.Id)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, community)
	assert.Equal(t, existingCommunity.Id, community.Id)
	assert.Equal(t, existingCommunity.Name, community.Name)
	assert.Equal(t, existingCommunity.Purpose, community.Purpose)
	assert.Equal(t, existingCommunity.ImageUrl, community.ImageUrl)
	assert.Equal(t, int(existingCommunity.NumberSubscriptions), community.NumberSubscriptions)
}

func TestGetCommunityNotFound(t *testing.T) {
	/*
		GIVEN: No community exists with the given ID
		WHEN:  GetPostgresqlCommunity is called with non-existent ID
		THEN:  A not found error is returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewCommunityAdapterTestWrapper(t)

	nonExistentId := uuid.New()

	// WHEN
	community, err := adapter.GetPostgresqlCommunity(nonExistentId)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, community)
	assert.Equal(t, errors.ObjectNotFoundError.CommunityNotFound, *err)
}

func TestFetchCommunitiesSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple communities exist in the database
		WHEN:  FetchPostgresqlCommunities is called
		THEN:  All communities are returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewCommunityAdapterTestWrapper(t)

	community1Name := "Community 1"
	community2Name := "Community 2"
	community1 := factories.NewCommunityModel(db, factories.CommunityModelF{
		Name: &community1Name,
	})
	community2 := factories.NewCommunityModel(db, factories.CommunityModelF{
		Name: &community2Name,
	})

	// WHEN
	communities, err := adapter.FetchPostgresqlCommunities()

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, communities)
	assert.GreaterOrEqual(t, len(communities), 2)

	// Find our created communities
	foundCommunity1 := false
	foundCommunity2 := false
	for _, community := range communities {
		if community.Id == community1.Id {
			foundCommunity1 = true
			assert.Equal(t, community1.Name, community.Name)
			assert.Equal(t, community1.Purpose, community.Purpose)
		}
		if community.Id == community2.Id {
			foundCommunity2 = true
			assert.Equal(t, community2.Name, community.Name)
			assert.Equal(t, community2.Purpose, community.Purpose)
		}
	}
	assert.True(t, foundCommunity1)
	assert.True(t, foundCommunity2)
}

func TestFetchCommunitiesEmpty(t *testing.T) {
	/*
		GIVEN: No communities exist in the database
		WHEN:  FetchPostgresqlCommunities is called
		THEN:  An empty list is returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewCommunityAdapterTestWrapper(t)

	// WHEN
	communities, err := adapter.FetchPostgresqlCommunities()

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, communities)
	assert.Equal(t, 0, len(communities))
}
