package membership_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

func TestGetMembershipSuccessfully(t *testing.T) {
	/*
		GIVEN: A membership record exists in the database
		WHEN:  GetMembership is called with valid ID
		THEN:  The membership record should be returned
	*/
	// GIVEN
	membershipController, _, db := controllerTest.NewMembershipControllerTestWrapper(t)

	// Create a user
	user := &model.User{
		Email:         utilsTest.GenerateRandomEmail(),
		Name:          "John",
		FirstLastName: "Doe",
		Rol:           model.UserRolClient,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err := db.Create(user).Error
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

	// Create a membership record
	startDate := time.Now()
	endDate := time.Now().AddDate(0, 1, 0) // 1 month later
	membership := &model.Membership{
		Description: "Monthly membership",
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      model.MembershipStatusActive,
		CommunityId: community.Id,
		UserId:      user.Id,
		PlanId:      plan.Id,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err = db.Create(membership).Error
	assert.NoError(t, err)

	// WHEN
	result, errResult := membershipController.GetMembership(membership.Id)

	// THEN
	assert.Nil(t, errResult)
	assert.NotNil(t, result)
	assert.Equal(t, membership.Id, result.Id)
	assert.Equal(t, membership.Description, result.Description)
	assert.Equal(t, string(membership.Status), string(result.Status))
	assert.Equal(t, membership.CommunityId, result.CommunityId)
	assert.Equal(t, membership.UserId, result.UserId)
	assert.Equal(t, membership.PlanId, result.PlanId)
}

func TestGetMembershipNotFound(t *testing.T) {
	/*
		GIVEN: No membership record exists with the given ID
		WHEN:  GetMembership is called with non-existent ID
		THEN:  An error should be returned
	*/
	// GIVEN
	membershipController, _, _ := controllerTest.NewMembershipControllerTestWrapper(t)
	nonExistentId := uuid.New()

	// WHEN
	result, errResult := membershipController.GetMembership(nonExistentId)

	// THEN
	assert.NotNil(t, errResult)
	assert.Nil(t, result)
	assert.Contains(t, errResult.Message, "not found")
}

func TestGetMembershipWithNilId(t *testing.T) {
	/*
		GIVEN: A nil UUID
		WHEN:  GetMembership is called with nil UUID
		THEN:  An error should be returned
	*/
	// GIVEN
	membershipController, _, _ := controllerTest.NewMembershipControllerTestWrapper(t)
	nilId := uuid.Nil

	// WHEN
	result, errResult := membershipController.GetMembership(nilId)

	// THEN
	assert.NotNil(t, errResult)
	assert.Nil(t, result)
}
