package plan_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	adapterTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/adapter"
)

func TestCreatePlanSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid plan data
		WHEN:  CreatePostgresqlPlan is called
		THEN:  A new plan is created and returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewPlanAdapterTestWrapper(t)

	fee := 29.99
	planType := model.PlanTypeMonthly
	reservationLimit := 5
	updatedBy := "test-user"

	// WHEN
	plan, err := adapter.CreatePostgresqlPlan(fee, planType, &reservationLimit, updatedBy)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, plan)
	assert.NotEmpty(t, plan.Id)
	assert.Equal(t, fee, plan.Fee)
	assert.Equal(t, planType, plan.Type)
	assert.Equal(t, &reservationLimit, plan.ReservationLimit)
}

func TestCreatePlanWithEmptyUpdatedBy(t *testing.T) {
	/*
		GIVEN: Valid plan data but empty updatedBy
		WHEN:  CreatePostgresqlPlan is called
		THEN:  An error is returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewPlanAdapterTestWrapper(t)

	fee := 29.99
	planType := model.PlanTypeMonthly
	reservationLimit := 5
	updatedBy := ""

	// WHEN
	plan, err := adapter.CreatePostgresqlPlan(fee, planType, &reservationLimit, updatedBy)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, plan)
	assert.Equal(t, errors.BadRequestError.InvalidUpdatedByValue, *err)
}

func TestCreatePlanWithZeroFee(t *testing.T) {
	/*
		GIVEN: Plan data with zero fee
		WHEN:  CreatePostgresqlPlan is called
		THEN:  A new plan is created with zero fee
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewPlanAdapterTestWrapper(t)

	fee := 0.0
	planType := model.PlanTypeMonthly
	reservationLimit := 2
	updatedBy := "test-user"

	// WHEN
	plan, err := adapter.CreatePostgresqlPlan(fee, planType, &reservationLimit, updatedBy)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, plan)
	assert.Equal(t, fee, plan.Fee)
	assert.Equal(t, planType, plan.Type)
	assert.Equal(t, &reservationLimit, plan.ReservationLimit)
}

func TestCreatePlanWithHighFee(t *testing.T) {
	/*
		GIVEN: Plan data with high fee
		WHEN:  CreatePostgresqlPlan is called
		THEN:  A new plan is created with the high fee
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewPlanAdapterTestWrapper(t)

	fee := 99.99
	planType := model.PlanTypeAnual
	reservationLimit := 20
	updatedBy := "test-user"

	// WHEN
	plan, err := adapter.CreatePostgresqlPlan(fee, planType, &reservationLimit, updatedBy)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, plan)
	assert.Equal(t, fee, plan.Fee)
	assert.Equal(t, planType, plan.Type)
	assert.Equal(t, &reservationLimit, plan.ReservationLimit)
}
