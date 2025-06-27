package plan_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	adapterTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/adapter"
)

func TestBulkCreatePlansSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid plan data for bulk creation
		WHEN:  BulkCreatePostgresqlPlans is called
		THEN:  All plans are created and returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewPlanAdapterTestWrapper(t)

	reservationLimit1 := 10
	reservationLimit2 := 20

	plansData := []*schemas.CreatePlanRequest{
		{
			Fee:              29.99,
			Type:             model.PlanTypeMonthly,
			ReservationLimit: &reservationLimit1,
		},
		{
			Fee:              59.99,
			Type:             model.PlanTypeAnual,
			ReservationLimit: &reservationLimit2,
		},
		{
			Fee:              99.99,
			Type:             model.PlanTypeMonthly,
			ReservationLimit: nil, // Unlimited
		},
	}
	updatedBy := "test-admin"

	// WHEN
	plans, err := adapter.BulkCreatePostgresqlPlans(plansData, updatedBy)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, plans)
	assert.Equal(t, 3, len(plans))

	// Verify first plan
	assert.NotEmpty(t, plans[0].Id)
	assert.Equal(t, 29.99, plans[0].Fee)
	assert.Equal(t, model.PlanTypeMonthly, plans[0].Type)
	assert.Equal(t, &reservationLimit1, plans[0].ReservationLimit)

	// Verify second plan
	assert.NotEmpty(t, plans[1].Id)
	assert.Equal(t, 59.99, plans[1].Fee)
	assert.Equal(t, model.PlanTypeAnual, plans[1].Type)
	assert.Equal(t, &reservationLimit2, plans[1].ReservationLimit)

	// Verify third plan
	assert.NotEmpty(t, plans[2].Id)
	assert.Equal(t, 99.99, plans[2].Fee)
	assert.Equal(t, model.PlanTypeMonthly, plans[2].Type)
	assert.Nil(t, plans[2].ReservationLimit)
}

func TestBulkCreatePlansWithEmptyUpdatedBy(t *testing.T) {
	/*
		GIVEN: Valid plan data but empty updatedBy
		WHEN:  BulkCreatePostgresqlPlans is called
		THEN:  An error is returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewPlanAdapterTestWrapper(t)

	reservationLimit := 10
	plansData := []*schemas.CreatePlanRequest{
		{
			Fee:              29.99,
			Type:             model.PlanTypeMonthly,
			ReservationLimit: &reservationLimit,
		},
	}
	updatedBy := ""

	// WHEN
	plans, err := adapter.BulkCreatePostgresqlPlans(plansData, updatedBy)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, plans)
	assert.Equal(t, errors.BadRequestError.InvalidUpdatedByValue, *err)
}

func TestBulkCreatePlansWithEmptyList(t *testing.T) {
	/*
		GIVEN: Empty plan data list
		WHEN:  BulkCreatePostgresqlPlans is called
		THEN:  An error is returned (as the adapter doesn't handle empty lists)
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewPlanAdapterTestWrapper(t)

	plansData := []*schemas.CreatePlanRequest{}
	updatedBy := "test-admin"

	// WHEN
	plans, err := adapter.BulkCreatePostgresqlPlans(plansData, updatedBy)

	// THEN
	assert.NotNil(t, err) // Adapter returns error for empty list
	assert.Nil(t, plans)
	assert.Equal(t, errors.BadRequestError.PlanNotCreated, *err)
}

func TestBulkCreatePlansWithDifferentTypes(t *testing.T) {
	/*
		GIVEN: Plan data with all different plan types
		WHEN:  BulkCreatePostgresqlPlans is called
		THEN:  Plans are created with correct types
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewPlanAdapterTestWrapper(t)

	reservationLimit := 15

	plansData := []*schemas.CreatePlanRequest{
		{
			Fee:              19.99,
			Type:             model.PlanTypeMonthly,
			ReservationLimit: &reservationLimit,
		},
		{
			Fee:              49.99,
			Type:             model.PlanTypeAnual,
			ReservationLimit: &reservationLimit,
		},
		{
			Fee:              149.99,
			Type:             model.PlanTypeMonthly,
			ReservationLimit: nil,
		},
	}
	updatedBy := "test-admin"

	// WHEN
	plans, err := adapter.BulkCreatePostgresqlPlans(plansData, updatedBy)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, plans)
	assert.Equal(t, 3, len(plans))

	// Verify plan types
	assert.Equal(t, model.PlanTypeMonthly, plans[0].Type)
	assert.Equal(t, model.PlanTypeAnual, plans[1].Type)
	assert.Equal(t, model.PlanTypeMonthly, plans[2].Type)
}

func TestBulkCreatePlansWithZeroFee(t *testing.T) {
	/*
		GIVEN: Plan data with zero fee (free plan)
		WHEN:  BulkCreatePostgresqlPlans is called
		THEN:  Plans are created with zero fee
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewPlanAdapterTestWrapper(t)

	reservationLimit := 5

	plansData := []*schemas.CreatePlanRequest{
		{
			Fee:              0.0, // Free plan
			Type:             model.PlanTypeMonthly,
			ReservationLimit: &reservationLimit,
		},
	}
	updatedBy := "test-admin"

	// WHEN
	plans, err := adapter.BulkCreatePostgresqlPlans(plansData, updatedBy)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, plans)
	assert.Equal(t, 1, len(plans))
	assert.Equal(t, 0.0, plans[0].Fee)
	assert.Equal(t, model.PlanTypeMonthly, plans[0].Type)
	assert.Equal(t, &reservationLimit, plans[0].ReservationLimit)
}
