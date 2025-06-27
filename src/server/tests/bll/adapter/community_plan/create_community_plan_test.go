package community_plan_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	adapterTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/adapter"
)

func TestCommunityPlanAdapter_CreatePostgresqlCommunityPlan_Success(t *testing.T) {
	// GIVEN: Valid community and plan exist
	communityPlanAdapter, _, db := adapterTest.NewCommunityPlanAdapterTestWrapper(t)

	// Create dependencies
	community := factories.NewCommunityModel(db, factories.CommunityModelF{})
	plan := factories.NewPlanModel(db, factories.PlanModelF{})

	// WHEN: CreatePostgresqlCommunityPlan is called
	result, err := communityPlanAdapter.CreatePostgresqlCommunityPlan(
		community.Id,
		plan.Id,
		"test_user",
	)

	// THEN: A new community-plan association is created
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, community.Id, result.CommunityId)
	assert.Equal(t, plan.Id, result.PlanId)
	assert.NotEqual(t, "", result.Id)
}

func TestCommunityPlanAdapter_CreatePostgresqlCommunityPlan_EmptyUpdatedBy(t *testing.T) {
	// GIVEN: Valid community and plan but empty updatedBy
	communityPlanAdapter, _, db := adapterTest.NewCommunityPlanAdapterTestWrapper(t)

	community := factories.NewCommunityModel(db, factories.CommunityModelF{})
	plan := factories.NewPlanModel(db, factories.PlanModelF{})

	// WHEN: CreatePostgresqlCommunityPlan is called with empty updatedBy
	result, err := communityPlanAdapter.CreatePostgresqlCommunityPlan(
		community.Id,
		plan.Id,
		"",
	)

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "Invalid updated by value", err.Message)
}

func TestCommunityPlanAdapter_GetPostgresqlCommunityPlan_Success(t *testing.T) {
	// GIVEN: A community-plan association exists
	communityPlanAdapter, _, db := adapterTest.NewCommunityPlanAdapterTestWrapper(t)

	community := factories.NewCommunityModel(db, factories.CommunityModelF{})
	plan := factories.NewPlanModel(db, factories.PlanModelF{})

	// Create the association
	created, createErr := communityPlanAdapter.CreatePostgresqlCommunityPlan(
		community.Id,
		plan.Id,
		"test_user",
	)
	assert.Nil(t, createErr)

	// WHEN: GetPostgresqlCommunityPlan is called
	result, err := communityPlanAdapter.GetPostgresqlCommunityPlan(
		community.Id,
		plan.Id,
	)

	// THEN: The association is returned
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, created.Id, result.Id)
	assert.Equal(t, community.Id, result.CommunityId)
	assert.Equal(t, plan.Id, result.PlanId)
}

func TestCommunityPlanAdapter_GetPostgresqlCommunityPlan_NotFound(t *testing.T) {
	// GIVEN: No community-plan association exists
	communityPlanAdapter, _, db := adapterTest.NewCommunityPlanAdapterTestWrapper(t)

	community := factories.NewCommunityModel(db, factories.CommunityModelF{})
	plan := factories.NewPlanModel(db, factories.PlanModelF{})

	// WHEN: GetPostgresqlCommunityPlan is called for non-existent association
	result, err := communityPlanAdapter.GetPostgresqlCommunityPlan(
		community.Id,
		plan.Id,
	)

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Message, "not found")
}

func TestCommunityPlanAdapter_DeletePostgresqlCommunityPlan_Success(t *testing.T) {
	// GIVEN: A community-plan association exists
	communityPlanAdapter, _, db := adapterTest.NewCommunityPlanAdapterTestWrapper(t)

	community := factories.NewCommunityModel(db, factories.CommunityModelF{})
	plan := factories.NewPlanModel(db, factories.PlanModelF{})

	// Create the association
	_, createErr := communityPlanAdapter.CreatePostgresqlCommunityPlan(
		community.Id,
		plan.Id,
		"test_user",
	)
	assert.Nil(t, createErr)

	// WHEN: DeletePostgresqlCommunityPlan is called
	err := communityPlanAdapter.DeletePostgresqlCommunityPlan(
		community.Id,
		plan.Id,
	)

	// THEN: The association is deleted successfully
	assert.Nil(t, err)

	// Verify it was deleted
	_, getErr := communityPlanAdapter.GetPostgresqlCommunityPlan(
		community.Id,
		plan.Id,
	)
	assert.NotNil(t, getErr)
}

func TestCommunityPlanAdapter_BulkCreatePostgresqlCommunityPlans_Success(t *testing.T) {
	// GIVEN: Multiple communities and plans
	communityPlanAdapter, _, db := adapterTest.NewCommunityPlanAdapterTestWrapper(t)

	community1 := factories.NewCommunityModel(db, factories.CommunityModelF{})
	community2 := factories.NewCommunityModel(db, factories.CommunityModelF{})
	plan1 := factories.NewPlanModel(db, factories.PlanModelF{})
	plan2 := factories.NewPlanModel(db, factories.PlanModelF{})

	requests := []*schemas.CreateCommunityPlanRequest{
		{
			CommunityId: community1.Id,
			PlanId:      plan1.Id,
		},
		{
			CommunityId: community2.Id,
			PlanId:      plan2.Id,
		},
	}

	// WHEN: BulkCreatePostgresqlCommunityPlans is called
	results, err := communityPlanAdapter.BulkCreatePostgresqlCommunityPlans(
		requests,
		"test_admin",
	)

	// THEN: Multiple associations are created
	assert.Nil(t, err)
	assert.NotNil(t, results)
	assert.Len(t, results, 2)

	assert.Equal(t, community1.Id, results[0].CommunityId)
	assert.Equal(t, plan1.Id, results[0].PlanId)
	assert.Equal(t, community2.Id, results[1].CommunityId)
	assert.Equal(t, plan2.Id, results[1].PlanId)
}

func TestCommunityPlanAdapter_BulkCreatePostgresqlCommunityPlans_EmptyUpdatedBy(t *testing.T) {
	// GIVEN: Valid data but empty updatedBy
	communityPlanAdapter, _, db := adapterTest.NewCommunityPlanAdapterTestWrapper(t)

	community := factories.NewCommunityModel(db, factories.CommunityModelF{})
	plan := factories.NewPlanModel(db, factories.PlanModelF{})

	requests := []*schemas.CreateCommunityPlanRequest{
		{
			CommunityId: community.Id,
			PlanId:      plan.Id,
		},
	}

	// WHEN: BulkCreatePostgresqlCommunityPlans is called with empty updatedBy
	results, err := communityPlanAdapter.BulkCreatePostgresqlCommunityPlans(
		requests,
		"",
	)

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, results)
	assert.Equal(t, "Invalid updated by value", err.Message)
}

func TestCommunityPlanAdapter_FetchPostgresqlCommunityPlans_Success(t *testing.T) {
	// GIVEN: Multiple community-plan associations exist
	communityPlanAdapter, _, db := adapterTest.NewCommunityPlanAdapterTestWrapper(t)

	community := factories.NewCommunityModel(db, factories.CommunityModelF{})
	plan1 := factories.NewPlanModel(db, factories.PlanModelF{})
	plan2 := factories.NewPlanModel(db, factories.PlanModelF{})

	// Create associations
	created1, err1 := communityPlanAdapter.CreatePostgresqlCommunityPlan(
		community.Id,
		plan1.Id,
		"test_user",
	)
	assert.Nil(t, err1)

	created2, err2 := communityPlanAdapter.CreatePostgresqlCommunityPlan(
		community.Id,
		plan2.Id,
		"test_user",
	)
	assert.Nil(t, err2)

	// WHEN: FetchPostgresqlCommunityPlans is called with community filter
	associations, err := communityPlanAdapter.FetchPostgresqlCommunityPlans(
		&community.Id,
		nil,
	)

	// THEN: The associations are returned
	assert.Nil(t, err)
	assert.NotNil(t, associations)
	assert.GreaterOrEqual(t, len(associations), 2)

	// Verify both associations are included
	foundAssoc1 := false
	foundAssoc2 := false
	for _, assoc := range associations {
		if assoc.Id == created1.Id {
			foundAssoc1 = true
		}
		if assoc.Id == created2.Id {
			foundAssoc2 = true
		}
	}
	assert.True(t, foundAssoc1)
	assert.True(t, foundAssoc2)
}

func TestCommunityPlanAdapter_FetchPostgresqlCommunityPlans_WithPlanFilter(t *testing.T) {
	// GIVEN: Multiple community-plan associations exist
	communityPlanAdapter, _, db := adapterTest.NewCommunityPlanAdapterTestWrapper(t)

	community1 := factories.NewCommunityModel(db, factories.CommunityModelF{})
	community2 := factories.NewCommunityModel(db, factories.CommunityModelF{})
	plan := factories.NewPlanModel(db, factories.PlanModelF{})

	// Create associations
	created1, err1 := communityPlanAdapter.CreatePostgresqlCommunityPlan(
		community1.Id,
		plan.Id,
		"test_user",
	)
	assert.Nil(t, err1)

	created2, err2 := communityPlanAdapter.CreatePostgresqlCommunityPlan(
		community2.Id,
		plan.Id,
		"test_user",
	)
	assert.Nil(t, err2)

	// WHEN: FetchPostgresqlCommunityPlans is called with plan filter
	associations, err := communityPlanAdapter.FetchPostgresqlCommunityPlans(
		nil,
		&plan.Id,
	)

	// THEN: The associations for that plan are returned
	assert.Nil(t, err)
	assert.NotNil(t, associations)
	assert.GreaterOrEqual(t, len(associations), 2)

	// Verify both associations are included
	foundAssoc1 := false
	foundAssoc2 := false
	for _, assoc := range associations {
		if assoc.Id == created1.Id {
			foundAssoc1 = true
		}
		if assoc.Id == created2.Id {
			foundAssoc2 = true
		}
	}
	assert.True(t, foundAssoc1)
	assert.True(t, foundAssoc2)
}

func TestCommunityPlanAdapter_FetchPostgresqlCommunityPlans_NoFilters(t *testing.T) {
	// GIVEN: Multiple community-plan associations exist
	communityPlanAdapter, _, db := adapterTest.NewCommunityPlanAdapterTestWrapper(t)

	community := factories.NewCommunityModel(db, factories.CommunityModelF{})
	plan := factories.NewPlanModel(db, factories.PlanModelF{})

	// Create association
	created, createErr := communityPlanAdapter.CreatePostgresqlCommunityPlan(
		community.Id,
		plan.Id,
		"test_user",
	)
	assert.Nil(t, createErr)

	// WHEN: FetchPostgresqlCommunityPlans is called without filters
	associations, err := communityPlanAdapter.FetchPostgresqlCommunityPlans(nil, nil)

	// THEN: All associations are returned
	assert.Nil(t, err)
	assert.NotNil(t, associations)
	assert.GreaterOrEqual(t, len(associations), 1)

	// Verify our association is included
	found := false
	for _, assoc := range associations {
		if assoc.Id == created.Id {
			found = true
			break
		}
	}
	assert.True(t, found)
}
