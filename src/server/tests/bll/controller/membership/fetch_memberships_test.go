package membership_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

func TestFetchMembershipsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple membership records exist in the database
		WHEN:  FetchMemberships is called
		THEN:  All membership records should be returned
	*/
	// GIVEN
	membershipController, _, db := controllerTest.NewMembershipControllerTestWrapper(t)

	// Create users
	user1 := &model.User{
		Email:         utilsTest.GenerateRandomEmail(),
		Name:          "John",
		FirstLastName: "Doe",
		Rol:           model.UserRolClient,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	user2 := &model.User{
		Email:         utilsTest.GenerateRandomEmail(),
		Name:          "Jane",
		FirstLastName: "Smith",
		Rol:           model.UserRolClient,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err := db.Create([]*model.User{user1, user2}).Error
	assert.NoError(t, err)

	// Create a community
	community := &model.Community{
		Name:    "Test Community",
		Purpose: "A test community",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err = db.Create(community).Error
	assert.NoError(t, err)

	// Create a plan
	reservationLimit := 10
	plan := &model.Plan{
		Fee:              100.00,
		Type:             model.PlanTypeMonthly,
		ReservationLimit: &reservationLimit,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err = db.Create(plan).Error
	assert.NoError(t, err)

	// Create membership records
	startDate := time.Now()
	endDate := time.Now().AddDate(0, 1, 0) // 1 month later
	memberships := []*model.Membership{
		{
			Description: "Monthly membership 1",
			StartDate:   startDate,
			EndDate:     endDate,
			Status:      model.MembershipStatusActive,
			CommunityId: community.Id,
			UserId:      user1.Id,
			PlanId:      plan.Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Description: "Monthly membership 2",
			StartDate:   startDate,
			EndDate:     endDate,
			Status:      model.MembershipStatusActive,
			CommunityId: community.Id,
			UserId:      user2.Id,
			PlanId:      plan.Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
	}
	err = db.Create(memberships).Error
	assert.NoError(t, err)

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
