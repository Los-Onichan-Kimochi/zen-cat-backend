package membership_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestFetchMembershipsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple membership records exist in the database
		WHEN:  FetchMemberships is called
		THEN:  All membership records should be returned
	*/
	// GIVEN
	membershipController, _, db := controllerTest.NewMembershipControllerTestWrapper(t)

	// Create entities using factories
	user1 := factories.NewUserModel(db)
	user2 := factories.NewUserModel(db)
	community := factories.NewCommunityModel(db)
	plan := factories.NewPlanModel(db)

	// Create membership records using factories
	_ = factories.NewMembershipModel(db, factories.MembershipModelF{
		CommunityId: &community.Id,
		UserId:      &user1.Id,
		PlanId:      &plan.Id,
	})
	_ = factories.NewMembershipModel(db, factories.MembershipModelF{
		CommunityId: &community.Id,
		UserId:      &user2.Id,
		PlanId:      &plan.Id,
	})

	// WHEN
	result, errResult := membershipController.FetchMemberships()

	// THEN
	assert.Nil(t, errResult)
	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, len(result.Memberships), 2)
}

func TestFetchMembershipsEmpty(t *testing.T) {
	/*
		GIVEN: No membership records exist in the database
		WHEN:  FetchMemberships is called
		THEN:  An empty list should be returned
	*/
	// GIVEN
	membershipController, _, _ := controllerTest.NewMembershipControllerTestWrapper(t)

	// WHEN
	result, errResult := membershipController.FetchMemberships()

	// THEN
	assert.Nil(t, errResult)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result.Memberships))
}
