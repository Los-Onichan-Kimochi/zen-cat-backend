package membership_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	adapterTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/adapter"
)

func TestMembershipAdapter_CreatePostgresqlMembership_Success(t *testing.T) {
	// GIVEN: Valid membership data with user, community, and plan
	membershipAdapter, _, db := adapterTest.NewMembershipAdapterTestWrapper(t)

	// Create dependencies
	user := factories.NewUserModel(db, factories.UserModelF{})
	community := factories.NewCommunityModel(db, factories.CommunityModelF{})
	plan := factories.NewPlanModel(db, factories.PlanModelF{})

	startDate := time.Now()
	endDate := startDate.AddDate(1, 0, 0) // One year later
	status := schemas.MembershipStatusActive
	reservationsUsed := 0

	// WHEN: CreatePostgresqlMembership is called
	result, err := membershipAdapter.CreatePostgresqlMembership(
		"Standard membership for community",
		startDate,
		endDate,
		status,
		&reservationsUsed,
		community.Id,
		user.Id,
		plan.Id,
		"test_user",
	)

	// THEN: A new membership is created and returned
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Standard membership for community", result.Description)
	assert.Equal(t, startDate.Format(time.RFC3339), result.StartDate.Format(time.RFC3339))
	assert.Equal(t, endDate.Format(time.RFC3339), result.EndDate.Format(time.RFC3339))
	assert.Equal(t, status, result.Status)
	assert.Equal(t, community.Id, result.CommunityId)
	assert.Equal(t, user.Id, result.UserId)
	assert.Equal(t, plan.Id, result.PlanId)
	assert.NotEqual(t, "", result.Id)
}

func TestMembershipAdapter_CreatePostgresqlMembership_EmptyUpdatedBy(t *testing.T) {
	// GIVEN: Valid membership data but empty updatedBy
	membershipAdapter, _, db := adapterTest.NewMembershipAdapterTestWrapper(t)

	// Create dependencies
	user := factories.NewUserModel(db, factories.UserModelF{})
	community := factories.NewCommunityModel(db, factories.CommunityModelF{})
	plan := factories.NewPlanModel(db, factories.PlanModelF{})

	startDate := time.Now()
	endDate := startDate.AddDate(1, 0, 0)
	status := schemas.MembershipStatusActive
	reservationsUsed := 0

	// WHEN: CreatePostgresqlMembership is called
	result, err := membershipAdapter.CreatePostgresqlMembership(
		"Standard membership for community",
		startDate,
		endDate,
		status,
		&reservationsUsed,
		community.Id,
		user.Id,
		plan.Id,
		"",
	)

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "Invalid updated by value", err.Message)
}

func TestMembershipAdapter_CreatePostgresqlMembership_ExpiredStatus(t *testing.T) {
	// GIVEN: Valid membership data with expired status
	membershipAdapter, _, db := adapterTest.NewMembershipAdapterTestWrapper(t)

	// Create dependencies
	user := factories.NewUserModel(db, factories.UserModelF{})
	community := factories.NewCommunityModel(db, factories.CommunityModelF{})
	plan := factories.NewPlanModel(db, factories.PlanModelF{})

	startDate := time.Now()
	endDate := startDate.AddDate(0, 6, 0) // Six months later
	status := schemas.MembershipStatusExpired
	reservationsUsed := 0

	// WHEN: CreatePostgresqlMembership is called
	result, err := membershipAdapter.CreatePostgresqlMembership(
		"Expired membership for review",
		startDate,
		endDate,
		status,
		&reservationsUsed,
		community.Id,
		user.Id,
		plan.Id,
		"test_admin",
	)

	// THEN: A new expired membership is created and returned
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Expired membership for review", result.Description)
	assert.Equal(t, status, result.Status)
	assert.Equal(t, community.Id, result.CommunityId)
	assert.Equal(t, user.Id, result.UserId)
	assert.Equal(t, plan.Id, result.PlanId)
	assert.NotEqual(t, "", result.Id)
}

func TestGetMembershipSuccessfully(t *testing.T) {
	/*
		GIVEN: A membership exists in the database
		WHEN:  GetPostgresqlMembership is called
		THEN:  The membership is returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewMembershipAdapterTestWrapper(t)

	membership := factories.NewMembershipModel(db, factories.MembershipModelF{})

	// WHEN
	result, err := adapter.GetPostgresqlMembership(membership.Id)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, membership.Id, result.Id)
	assert.Equal(t, membership.Description, result.Description)
}
